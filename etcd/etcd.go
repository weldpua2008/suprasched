package etcd

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type MyKV struct {
	Key string
	Value  string
	Revision int64
}

type MyKVtwo struct {
	//Key string
	Value  string
	//Name string
	//Revision int64
}

type MyTemp struct {
	ClusterID MyKV
	ObjectID  MyKV // This is the key!!
	Retry MyKV
	Revision MyKV
}

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func KV_getExample() { // Get function for example, without accepting external values

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})

	if err != nil {
		log.Fatalf("Cannot start etcd client, got %s", err)
		return
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer client.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// NewKV wraps a KV instance so that all requests are prefixed with a given string.
	kvClient := clientv3.NewKV(client)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	kvClient.Delete(ctx, "", clientv3.WithPrefix())

	kvPut, err := kvClient.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s", err)
	}

	rev := kvPut.Header.Revision
	fmt.Println("Revision:", rev)

	resp, err := kvClient.Get(ctx, "foo")

	if err != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Println("Key:", string(ev.Key), ",Value:", string(ev.Value), ",Revision:", rev)
	}

	rev := kvPut.Header.Revision
	fmt.Println("Revision:", rev)

	resp, err := kvClient.Get(ctx, *flagKey)

	if err != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Println("Key:", string(ev.Key), ",Value:", string(ev.Value), ",Revision:", rev)
	}

}

func KV_getExampleWithFlag() { // Get function for example, with external values

	flagKey := flag.String("key", "tempKeyValue", "key for Get KeyValue function in etcd")

	flag.Parse()

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})

	if err != nil {
		log.Fatalf("Cannot start etcd client, got %s", err)
		return
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer client.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// NewKV wraps a KV instance so that all requests are prefixed with a given string.
	kvClient := clientv3.NewKV(client)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	kvClient.Delete(ctx, "", clientv3.WithPrefix())

	myTime := time.Now().String()
	kvPut, err := kvClient.Put(context.TODO(), *flagKey, myTime)
	if err != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s", err)
	}

	rev := kvPut.Header.Revision
	fmt.Println("Revision:", rev)

	resp, err := kvClient.Get(ctx, *flagKey)

	if err != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Println("Key:", string(ev.Key), ",Value:", string(ev.Value), ",Revision:", rev)
	}

}

func KV_getExampleWithWriteToYaml(){ // Get function for example, write to yaml with external values

	flagKey := flag.String("key", "tempKeyValue", "key for Get KeyValue function in etcd")

	flag.Parse()

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{"127.0.0.1:2379"},
	})

	if err != nil {
		log.Fatalf("Cannot start etcd client, got %s", err)
		return
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer client.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// NewKV wraps a KV instance so that all requests are prefixed with a given string.
	kvClient := clientv3.NewKV(client)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	kvClient.Delete(ctx, "", clientv3.WithPrefix())

	myTime := time.Now().String()
	kvPut, err := kvClient.Put(context.TODO(), *flagKey, myTime)
	if err != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s",err)
	}

	rev := kvPut.Header.Revision
	fmt.Println("Revision:", rev)

	resp, err := kvClient.Get(ctx, *flagKey)

	if err != nil {
		log.Fatalf("Get Function: Cannot get value, got %s",err)
	}

	for _, ev := range resp.Kvs {
		fmt.Println("Key:", string(ev.Key), ",Value:" ,string(ev.Value), ",Revision:", rev)
		myKV := map[string]MyKV{"kv " + myTime: {string(ev.Key), string(ev.Value), rev}}
		data, err := yaml.Marshal(&myKV)
		if err != nil {
			log.Fatal(err)
		}
		err2 := ioutil.WriteFile("etcd/etcdConfigFile.yaml", data, 0777)

		if err2 != nil {

			log.Fatal(err2)
		}

		fmt.Println("data written")
	}
}

func KV_getExampleWithWriteToYaml_two(){ // Get function for example, write to yaml with external values

	flagKey := flag.String("key", "tempKeyValue", "key for Get KeyValue function in etcd")

	// Must be called after all flags are defined and before flags are accessed by the program.
	flag.Parse()

	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	client, cliErr := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if cliErr != nil {
		log.Fatalf("Cannot start etcd client, got %s", cliErr)
		return
	}

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer client.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)

	// NewKV wraps a KV instance so that all requests are prefixed with a given string.
	kvClient := clientv3.NewKV(client)

	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	
	
	
	
	kvClient.Delete(ctx, "", clientv3.WithPrefix()) // Temp!!!!

	temp := MyTemp{MyKVtwo{"15848"}, MyKVtwo{*flagKey}, MyKVtwo{"2"}, MyKVtwo{""}}
	s := reflect.ValueOf(&temp).Elem()
	typeOfT := s.Type()
	flag := 1
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		kvc, putErr := kvClient.Put(ctx, *flagKey + "/" + typeOfT.Field(i).Name, strings.Trim(fmt.Sprint(f.Interface()), "{}"))
		if putErr != nil {
			log.Fatalf("Put Function: Cannot put key & value, got %s", putErr)
		}
		if flag == 1{
			temp.Revision.Value = strconv.FormatInt(kvc.Header.Revision, 10)
			flag = 0
		}
	}
	
	
	
	
	
	// clientv3.OpOption = configures Operations like Get, Put, Delete.
	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	}

	// Get = return the value for "key"
	resp, getErr := kvClient.Get(ctx, *flagKey, opts...)
	if getErr != nil {
		log.Fatalf("Get Function: Cannot get value, got %s", getErr)
	}

	myKV := make(map[string]MyKV)

	for _, ev := range resp.Kvs {
		myKV[string(ev.Key)] = MyKV{string(ev.Value)}
	}

	data, marshalErr := yaml.Marshal(&myKV)
	if marshalErr != nil {
		log.Fatalf("yaml.Marshal Function: got %s", marshalErr)
	}

	// WriteFile = writes data to a file named by filename
	writeErr := ioutil.WriteFile("etcdConfigFile.yaml", data, 0777)
	if writeErr != nil {
		log.Fatalf("WriteFile Function: Cannot write to file, got %s", writeErr)
	}

	fmt.Println("data written")
}

