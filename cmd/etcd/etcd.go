package cmd

import (
	"flag"
	"fmt"
	"github.com/weldpua2008/suprasched/etcd"
	"suprasched/cmd"

	"os"
)

func main() {

	//etcd object get subcommand
	getCmd := flag.NewFlagSet("get", flag.ExitOnError) //exit on error
	putCmd := flag.NewFlagSet("put", flag.ExitOnError) //exit on error
	//input for etcd get command
	//getAll := getCmd.Bool("all", false, "Get all keys and values")
	getKey := getCmd.String("key", "", "key values")
	getEtcdIP := getCmd.String("ip", "", "etcd IP")
	putClusterID := putCmd.String("clusterid", "", "ClusterId value")
	putRetry := putCmd.String("retry", "", "retry count value as string")
	getRevision := getCmd.String("revision", "", "Revision value as string")

	//check if the user past a subcommand
	if len(os.Args) < 2 {
		fmt.Println("expected 'get' or 'put' subcommand")
		os.Exit(1)
	}

	//look at the 2nd argument's value
	switch os.Args[1] {
	case "get": //if it's the 'get' command
		//handle get here
		cmd.HandleGet(getCmd, getKey, getEtcdIP, getRevision)

	case "put":
		HandlePut(putCmd, getKey, getEtcdIP, putClusterID, putRetry)

	default: // if we don't understand the input

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
		//get the key values by revisionaa

	}

}

func HandlePut(putCmd *flag.FlagSet, key *string, ip *string, clusterId *string, retry *string) {

	putCmd.Parse(os.Args[2:])

	if *key == "" || *ip == "" || *clusterId == "" || *retry == "" {
		fmt.Println("Key, IP, clsuterId, and retry are required in oreder to set values")
		putCmd.PrintDefaults()
		os.Exit(1)
	}
	if *key != "" || *ip != "" || *clusterId != "" || *retry != "" {
		//etcd.ptuKV(*key, *ip, *clusterId, *retry)
	}

}
