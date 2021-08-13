package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"reflect"
	"strconv"
	"time"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func GetKV(key string, endpoint string) map[string]string { // do "defer client.Close()" and "defer cancel()"

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	etcdClient, etcdClientErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{endpoint},
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

	opts := []clientv3.OpOption{ // clientv3.OpOption = configures Operations like Get, Put, Delete.
		clientv3.WithPrefix(), // WithPrefix = requests to operate on the keys with matching prefix.
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	}

	getResponse, getErr := etcdClient.Get(ctx, key, opts...) // Get = return the value for "key" with option

	if getErr != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", getErr)
	}

	dataMap := make(map[string]string)

	for _, ev := range getResponse.Kvs {
		dataMap[string(ev.Key)+"/"+strconv.FormatInt(getResponse.Header.Revision+3, 10)] = string(ev.Value)
	}

	fmt.Println("Data written to map")

	return dataMap
}

func ptuKV(key string, endpoint string, clusterID string, retry string) {

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	etcdClient, etcdClientErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{endpoint},
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

	//etcdClient.Delete(ctx, "", clientv3.WithPrefix()) // Temp!!!!

	dataMap := map[string]string{"Cluster_ID": clusterID, "Object_ID": key, "Retry": retry, "Revision": ""}
	keys := reflect.ValueOf(dataMap).MapKeys()
	for i := 0; i < len(keys); i++ {
		_, putErr := etcdClient.Put(ctx, key+"/"+keys[i].String(), dataMap[keys[i].String()])
		if putErr != nil {
			log.Fatalf("Put Function: Cannot put key & value, got %s", putErr)
		}
	}
}
