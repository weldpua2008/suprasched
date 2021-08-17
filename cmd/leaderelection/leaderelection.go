package main

import (
	"flag"
	"fmt"
	//"github.com/weldpua2008/suprasched/cmd"
	//etcd "github.com/weldpua2008/suprasched/etcd"
	"context"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"time"
)

func main() {

	getCmd := flag.NewFlagSet("new", flag.ExitOnError) //exit on error

	getEndpoint := getCmd.String("endpoint", "", "endpoint URL usually localhost:22379")
	getPrefixKey := getCmd.String("pfx", "", "prefix key for leader") // "/leader-election/"
	getID := getCmd.String("id", "", "id or name for the object")

	if len(os.Args) < 2 {
		fmt.Println("expected 'new' or 'set' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "new": //if its the 'new' command
		//handle get here
		HandleNew(getCmd, getEndpoint, getPrefixKey, getID)
	default: // if we dont understand the input
	}

}

func HandleNew(getCmd *flag.FlagSet, endpoint *string, prefixkey *string, id *string) {
	getCmd.Parse(os.Args[2:])
	if *endpoint == "" || *prefixkey == "" {
		fmt.Println("endpoint and prefixKey are required to get values")
		getCmd.PrintDefaults()
		os.Exit(1)
	}
	if *endpoint != "" && *prefixkey != "" && *id == "" {
		*id = "random"
		newElection(endpoint, prefixkey, id)

	}
	if *endpoint != "" && *prefixkey != "" && *id != "" {
		newElection(endpoint, prefixkey, id)

	}

}

func newElection(endpoint *string, prefixkey *string, id *string) {

	//create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{*endpoint}})
	if err != nil {
		log.Fatal(err)

	}
	defer cli.Close()

	// create a sessions to elect a leader
	s, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	e := concurrency.NewElection(s, *prefixkey) // "/leader-election/"
	ctx := context.Background()

	// Elect a leader (or wait that the leader resign)
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("leader election for ", *id)

	fmt.Println("do something for in ", *id)
	time.Sleep(5 * time.Second)

	if err := e.Resign(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("resign ", *id)
}
