package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	//"reflect"
	//"strconv"
	"strings"
	"time"
)

type UniversalObject interface {
	//Data interface{}
	ObjID() string
}

//func (job Job) ObjID() string {
//	return job.Id
//}

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 10 * time.Second
	//waitTimeout    = 2 * time.Second
	tempObj UniversalObject
)

func PutKV_v2(obj UniversalObject, key string, endpoint string){

	obj_type := CheckStructType(obj) // Checking the type of struct we received
	json_data, _ := json.Marshal(obj) // Converts the object to a json file as one long string

	// ---------------------------------------------------------------------------------------------------------

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	etcdClient, etcdClientErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{endpoint},
	})

	if etcdClientErr != nil {
		log.Fatalf("Cannot start etcd client, got %s", etcdClientErr)
		return
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer etcdClient.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	// ---------------------------------------------------------------------------------------------------------

	etcdClient.Delete(ctx, "", clientv3.WithPrefix())
	_, putErr := etcdClient.Put(ctx, key + "/" + obj_type + "/" + obj.ObjID(), string(json_data)) // put the obj id
	if putErr != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s", putErr)
	}
}

func GetKV_v2(key string, endpoint string) interface{} { // Get function with JSON

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	etcdClient, etcdClientErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{endpoint},
	})

	if etcdClientErr != nil {
		log.Fatalf("Cannot start etcd client, got %s", etcdClientErr)
		return nil
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer etcdClient.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	// ---------------------------------------------------------------------------------------------------------

	opts := []clientv3.OpOption{ // clientv3.OpOption = configures Operations like Get, Put, Delete.
		clientv3.WithPrefix(), // WithPrefix = requests to operate on the keys with matching prefix.
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	}

	getResponse, getErr := etcdClient.Get(ctx, key, opts...) // Get = return the value for "key" with option

	if getErr != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", getErr)
	}

	// ---------------------------------------------------------------------------------------------------------

	var dat_two string
	var key_name []string
	for _, ev := range getResponse.Kvs {
		dat_two = string(ev.Value)
		key_name = strings.Split(string(ev.Key), "/")
	}

	// ---------------------------------------------------------------------------------------------------------

	CheckGetStructType(key_name[1])
	unmarshal_err := json.Unmarshal([]byte(dat_two), &tempObj)
	if unmarshal_err != nil {
		fmt.Printf("json.Unmarshal Function: Cannot do unmarshal, got %s\n", unmarshal_err)
	}

	return tempObj
}

func CheckGetStructType(obj_type string){
	switch obj_type {
	case "Job":
		tempObj = new(Job)
	//case "Cluster":
	//	myStruct = new(Cluster)
	//case "ClusterRegistry":
	//	myStruct = new(ClusterRegistry)
	//case "Registry":
	//	myStruct = new(Registry)
	default:
		fmt.Println("Error unidentified object from json string in get function")
		return
	}
}

func CheckStructType(obj interface{}) string{
	switch obj.(type) {
	case *Job, Job:
		return "Job"
	//case *Cluster, Cluster:
	//	fmt.Println("Object type: Cluster")
	//	return "Cluster"
	//case *ClusterRegistry, ClusterRegistry:
	//	fmt.Println("Object type: ClusterRegistry")
	//	return "ClusterRegistry"
	//case *Registry, Registry:
	//	fmt.Println("Object type: Registry")
	//	return "Registry"
	default:
		fmt.Println("Error unidentified struct")
		return ""
	}
}
