/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	//"github.com/weldpua2008/suprasched/cmd"
	etcd "github.com/weldpua2008/suprasched/etcd"
	"os"
)

func main() {
	//cmd.Execute()

	//etcd object get subcommand
	getCmd := flag.NewFlagSet("get", flag.ExitOnError) //exit on error
	//input for etcd get command
	//getAll := getCmd.Bool("all", false, "Get all keys and values")
	getKey := getCmd.String("key","", "key values")
	getEtcdIP := getCmd.String("ip","","etcd IP")
	//getObjectID := getCmd.String("objectId","", "ObjectId value")
	//getRetry := getCmd.String("retry","", "retry count value as string")
	getRevisoion := getCmd.String("revision","", "Revision value as string")

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

func HandleGet(getCmd *flag.FlagSet, key *string, ip *string, revision *string)  {
	getCmd.Parse(os.Args[2:])
	if *key =="" || *ip =="" {
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