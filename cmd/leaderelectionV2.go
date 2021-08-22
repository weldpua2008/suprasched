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
package cmd

import (
	"fmt"
	//"github.com/weldpua2008/suprasched/cmd"
	//etcd "github.com/weldpua2008/suprasched/etcd" aa
	"context"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"os"
	"time"
)

//var newE string

//func init()  {
//	leaderelectionV2Cmd.PersistentFlags().StringVar(&newE, "new", "", "new election")

//}

// leaderelectionV2Cmd represents the leaderelectionV2 command
var leaderelectionV2Cmd = &cobra.Command{
	Use:   "leaderelectionV2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("leaderelectionV2 called")
	},
}

var newLeader = &cobra.Command{
	Use: "newElection",
	Short: "start a new election",
	Long: "with key and ip for endpoint newElection connect to a client(etcd) and start leader election",
	Run: func(cmd *cobra.Command, args []string) {

		HandleNewV2(endpoint, prefixkey)

	},
}

var prefixkey, endpoint string



func init() {
	rootCmd.AddCommand(leaderelectionV2Cmd)
	leaderelectionV2Cmd.AddCommand(newLeader)

	newLeader.Flags().StringVar(&endpoint, "endpoint", "", "IP for client")
	newLeader.Flags().StringVar(&prefixkey, "pfx", "", "a prefix key of the componentpfx")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// leaderelectionV2Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// leaderelectionV2Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func HandleNewV2(endpoint string, prefixkey string){
	if endpoint == "" || prefixkey == "" {
		fmt.Println("endpoint and prefixKey are required to get values")
		os.Exit(1)
	}
	if endpoint != "" && prefixkey != ""{
		newElectionV2(endpoint, prefixkey)

	}


}

func newElectionV2(endpoint string, prefixkey string)  {

	//create a etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{endpoint}})
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

	e := concurrency.NewElection(s, prefixkey) // "/leader-election/"
	ctx := context.Background()

	// Elect a leader (or wait that the leader resign)
	if err := e.Campaign(ctx, "e"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("leader election for ", prefixkey)

	fmt.Println("do something for in ", prefixkey)
	time.Sleep(5 * time.Second)

	if err := e.Resign(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("resign ", prefixkey)


}
