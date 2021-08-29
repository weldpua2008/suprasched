package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/weldpua2008/suprasched/core"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"sync"
	"time"
)

type EtcdClient struct {
	cfg *clientv3.Config
	client *clientv3.Client
}

//type Job struct {
//	UUID string
//}

var (
	tempObject Object
)

type ObjectCache struct {
	mu sync.RWMutex
	clusters map[core.Namespace]map[core.UID]*core.Cluster // a map from namespace to a map of clusters.
	prefixesChannels map[string]clientv3.WatchChan
	etcdClient *EtcdClient
}

type Object interface {
	GetObjID() string
	Store() bool
	Retrieve() interface{}
}

// Put, Check Object Type & Store --------------------------------------------------------------------------------------------------------------------------

func Store (obj Object, etcdClient *EtcdClient, ctx context.Context, prefix string) bool {
	jsonData, jsonDataErr := json.Marshal(&obj) // Converts the object to a json
	if jsonDataErr != nil {
		log.Fatal("Store function: Error converting object data to json\n")
		return false
	}
	objType := CheckObjectType(obj)
	if objType == "" {
		log.Fatal("Store function: Error unidentified object\n")
		return false
	}
	return etcdClient.Put(ctx, fmt.Sprintf("/%s/%s/%s", prefix, objType , obj.GetObjID()), base64.StdEncoding.EncodeToString(jsonData))
}

func CheckObjectType(obj interface{}) string{
	switch obj.(type) {
	case *core.Job, core.Job:
		return "Job"
	case *core.Cluster, core.Cluster:
		return "Cluster"
	default:
		return ""
	}
}

func (etcdClient *EtcdClient) Put (ctx context.Context, key string, value string) bool {
	_, putErr := etcdClient.client.Put(ctx, key, value)
	if putErr != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s\n", putErr)
		return false
	}
	return true
}

// Get & Retrieve --------------------------------------------------------------------------------------------------------------------------

func (etcdClient *EtcdClient) Retrieve (ctx context.Context, prefix string) interface{} {
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(), // WithPrefix = requests to operate on the keys with matching prefix.
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend), // WithSort = specifies the ordering in 'Get' request.
	}

	getResponse, getErr := etcdClient.client.Get(ctx, prefix, opts...) // Get = return the value for "key" with option
	if getErr != nil {
		log.Fatalf("Retrieve Function: Cannot get value, got %s", getErr)
	}

	for _, ev := range getResponse.Kvs {
		decodeJsonData, errDecode := base64.StdEncoding.DecodeString(string(ev.Value))
		if errDecode != nil { // This line checks if decode of the event.Kv.Value went good or not.
			log.Fatalf("Retrieve Function: Cannot get decode json string, git %s\n", errDecode)
			return nil
		}
		keyName := strings.Split(string(ev.Key), "/")
		return etcdClient.GetJob(decodeJsonData, keyName[1])
	}
	return nil
}

func CreateObjType (objType string) interface{} {
	switch objType {
	case "Job":
		return &core.Job{}
	case "Cluster":
		return &core.Cluster{}
	default:
		log.Fatalf("CreateObjType Function: Error unidentified object from json\n")
		return nil
	}
}

func (etcdClient *EtcdClient)GetJob(data []byte, objType string) interface{} {
	tempObject := CreateObjType(objType)
	if tempObject == nil {
		log.Fatalf("GetJob Function: Error unidentified object from json\n")
		return nil
	}
	unmarshalErr := json.Unmarshal(data, &tempObject)
	if unmarshalErr != nil {
		fmt.Printf("GetJob Function: Cannot do unmarshal אם גשאש, got %s\n", unmarshalErr)
	}
	return tempObject
}

// Watch, Refresh & Add Prefix --------------------------------------------------------------------------------------------------------------------------

func (etcdClient *EtcdClient)Watch(ctx context.Context, prefix string) clientv3.WatchChan {
	return etcdClient.client.Watch(ctx, prefix, clientv3.WithPrefix())
}

func (obj *ObjectCache)AddPrefix(ctx context.Context, prefix string) {

	obj.mu.Lock()
	defer obj.mu.Unlock()
	if len(prefix) < 1 {
		return
	}
	if _, ok := obj.prefixesChannels[prefix]; !ok {
		obj.prefixesChannels[prefix] = obj.etcdClient.Watch(ctx, prefix)
	}
}

func (obj *ObjectCache)Refresh(ctx context.Context, stopChannel chan bool) {
	var result chan *core.Cluster // This is a channel of type core.Cluster with size 1 -> need channel size ?!?
	for _, watchChannel := range obj.prefixesChannels { // This for loop takes channel from the obj.prefixesChannels map and put it in a var called "watchChannel"
		go func() {
			for true {
				select {
				case watchResp := <-watchChannel: // watchChannel put the channel data inside watchResp var.
					for _, event := range watchResp.Events { // This for loop takes the all the events and put it inside "event" var.
						fmt.Printf("%s %q : %q\n", event.Type, event.Kv.Key, event.Kv.Value)
						if dec, errDecode := base64.StdEncoding.DecodeString(string(event.Kv.Value)); errDecode != nil { // This line checks if decode of the event.Kv.Value went good or not.
							var cl *core.Cluster
							if errJsonUnmarshal := json.Unmarshal(dec, &cl); errJsonUnmarshal == nil { // This line checks if the json.Unmarshal put decode data inside the core.Cluster object
								result <- cl // Put the core.Cluster object inside the result chan *core.Cluster -> the new data is run over the old data ?!?
							}
						}
					}
				case <-stopChannel: // This statement indicates that the element receives data from the channel(stopChannel). If the result of the received statement is not going to use is also a valid statement.
					fmt.Println("Done watching.")
					//need to put here cancel() from context
					return
				}
			}
		}()
	}

	time.Sleep(time.Second)

	for clusterToCache := range result{
		obj.mu.Lock()
		obj.clusters = make(map[core.Namespace]map[core.UID]*core.Cluster)
		obj.clusters[clusterToCache.Namespace][clusterToCache.UID] = clusterToCache
		obj.mu.Unlock()
		select {
		case <-ctx.Done(): // Done returns a channel that's closed when work done on behalf of this context should be canceled. Done may return nil if this context can never be canceled.
			return
		}
	}
}