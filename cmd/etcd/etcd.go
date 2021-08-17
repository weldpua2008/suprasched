package main

import (
	"flag"
	"fmt"
	etcd "github.com/weldpua2008/suprasched/etcd"
	"os"
)

func main() {

	//etcd object get subcommand
	getCmd := flag.NewFlagSet("get", flag.ExitOnError) //exit on error
	//input for etcd get command
	//getAll := getCmd.Bool("all", false, "Get all keys and values")
	getKey := getCmd.String("key", "", "key values")
	getEtcdIP := getCmd.String("ip", "", "etcd IP")
	//getObjectID := getCmd.String("objectId","", "ObjectId value")
	//getRetry := getCmd.String("retry","", "retry count value as string")
	getRevisoion := getCmd.String("revision", "", "Revision value as string")

	//check if the user past a subcommand
	if len(os.Args) < 2 {
		fmt.Println("expected 'get' or 'set' subcommand")
		os.Exit(1)
	}

	//look at the 2nd argument's value
	switch os.Args[1] {
	case "get": //if its the 'get' command
		//handle get here
		HandleGet(getCmd, getKey, getEtcdIP, getRevisoion)
	default: // if we dont understand the input

	}
}

func HandleGet(getCmd *flag.FlagSet, key *string, ip *string, revision *string) {
	getCmd.Parse(os.Args[2:])
	if *key == "" || *ip == "" {
		fmt.Println("key and IP are required to get values")
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	if *key != "" && *ip != "" && *revision == "" {
		m := etcd.GetKV(*key, *ip)
		fmt.Println("map:", m)

	}
	if *key != "" && *ip != "" && *revision != "" {
		//get the key values by revision

	}

}
