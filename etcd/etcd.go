package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type TempObj struct {
	Name string
	Age int
	Designation string
	Salary int
}

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func EtcdConnection() {
	// The etcd client object is instantiated, configured with the dial time and the endpoint to the local etcd server
	KVClient, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatalf("Cannot start etcd client, got %s", err)
		return
	}
	fmt.Println("connect success")

	// The defer call is guaranteed to be used at the end of the function and ensures all etcd resources are released
	defer KVClient.Close()

	// context = Go feature that allows code across a goroutine to access shared information in a safe manner and cancel operations
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	// releases resources if slowOperation completes before timeout elapses
	defer cancel()

	// NewKV wraps a KV instance so that all requests are prefixed with a given string.
	kv := clientv3.NewKV(KVClient)
}

func GetValue(ctx context.Context, kv clientv3.KV, key string) {
	// returns a response object with a header that contains the revision and a "Kvs" (key values) field, which is a slice of key-value pairs.

	opts := []clientv3.OpOption{
		//clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
	}

	gr, err := kv.Get(ctx, key, opts...)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range gr.Kvs {
		fmt.Println("Value:", string(item.Key), ",Value:" ,string(item.Value), ",Revision:", string(gr.Header.Revision))
	}
}

func GetPastValueFromRevision(ctx context.Context, kv clientv3.KV, key string, rev int64) {
	// WithRev() returns the historical version of the value for the target key, but the header's revision will always contain the current revision!.
	gr, _ := kv.Get(ctx, key, clientv3.WithRev(rev))
	fmt.Println("Historical version of the value:", string(gr.Kvs[0].Value), ",Past revision:", rev, ",Revision: ", gr.Header.Revision)
}

func InsertsSingleValue(ctx context.Context, kv clientv3.KV, key string, val string ){
	// Inserts a key value pair "key", "value" via the Put() method. the keys and values must be strings!
	pr, err := kv.Put(ctx, key, val)
	if err != nil {
		log.Fatalf("Cannot put key/value, got %s", err)
		return
	}
	rev := pr.Header.Revision
	fmt.Println("Revision:", rev)
}