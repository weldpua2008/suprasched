package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func KV_getExample(){ // Get function for example, without accepting external values

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

	kvPut, err := kvClient.Put(context.TODO(), "foo", "bar")
	if err != nil {
		log.Fatalf("Put Function: Cannot put key & value, got %s",err)
	}

	rev := kvPut.Header.Revision
	fmt.Println("Revision:", rev)

	resp, err := kvClient.Get(ctx, "foo")

	if err != nil {
		log.Fatalf("Get Function: Cannot get value, got %s",err)
	}

	for _, ev := range resp.Kvs {
		fmt.Println("Key:", string(ev.Key), ",Value:" ,string(ev.Value), ",Revision:", rev)
	}
}

func KV_getExampleWithFlag(){ // Get function for example, with external values

	flagKey := flag.String("key", "tempKeyValue", "key for Get KeyValue function in etcd")

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
	}
}