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


func PutKV_v2(obj interface{}){


	fmt.Printf("obj: %v\n",obj)
	fmt.Println(obj)


	obj_type := CheckStructType(obj)
	fmt.Println(obj_type)
	json_data, _ := json.Marshal(obj)

	// ---------------------------------------------------------------------------------------------------------

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	etcdClient, etcdClientErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{"127.0.0.1:2379"},
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

	etcdClient.Delete(ctx, "", clientv3.WithPrefix()) // Temp!!!!
	_, putErr := etcdClient.Put(ctx, "tempKeyValue/" + obj_type + "/2693", string(json_data))
	if putErr != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s", putErr)
	}


func GetKV_v2(key string) { // Get function with JSON

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
		fmt.Println("Key:", string(ev.Key), ",Value:" ,string(ev.Value))
		dat_two = string(ev.Value)
		key_name = strings.Split(string(ev.Key), "/")
	}

	// ---------------------------------------------------------------------------------------------------------

	CheckGetStructType(key_name[1])
	fmt.Println("myStruct.MyID")
	//fmt.Println(myStruct.MyID())
	err := json.Unmarshal([]byte(dat_two), &myStruct)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myStruct.MyID())
	fmt.Println(myStruct)
}

func CheckGetStructType(struct_type string){
	switch struct_type {
	case "Job":
		myStruct = new(Job)
	//case "MyKVtwo":
	//	myStruct =  UniversalDTO{MyKVtwo{}}
	default:
		fmt.Println("Error interface")
		return
	}
}

func CheckStructType(x interface{}) string{
	switch x.(type) {
	case *Job, Job:
		fmt.Println("Object type: Job")
		return "Job"
	//case MyKVtwo:
	//	fmt.Println("MyKVtwo type")
	default:
		fmt.Println("Error")
		return ""
	}
}