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
	"github.com/spf13/cobra"
	etcd "github.com/weldpua2008/suprasched/etcd"
	"os"
)
var get, put string

func init()  {
	etcdV2Cmd.PersistentFlags().StringVar(&get, "get", "", "get subcommand")
	etcdV2Cmd.PersistentFlags().StringVar(&put, "put", "", "put subcommand")


}


// etcdV2Cmd represents the etcdV2 command
var etcdV2Cmd = &cobra.Command{
	Use:   "etcdV2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("etcdV2 called")

		fmt.Printf("get %v\n", get)
		fmt.Printf("put %v\n", put)
	},
}

var makeGet = &cobra.Command{
	Use: "makeGet",
	Short: "make a get command",
	Long: "By puting key and ip you get the value of the key from the etcd server",
	Run: func(cmd *cobra.Command, args []string) {

		HandleGetV2(getKey, getEtcdIP)

	},
}

var makePut = &cobra.Command{
	Use: "makePut",
	Short: "make a put command",
	Long: "By puting key and ip you get the value of the key from the etcd server",
	Run: func(cmd *cobra.Command, args []string) {

		//HandlePut(key string, ip string, clusterId string, retry string)

	},
}
//add commands
func addCommands()  {
	etcdV2Cmd.AddCommand(makeGet)
	etcdV2Cmd.AddCommand(makePut)

}

//add flags
var getKey, getEtcdIP string
var putKey, putEtcdIp, putClusterId, putRetry string

func init() {
	rootCmd.AddCommand(etcdV2Cmd)
	addCommands()
	makeGet.Flags().StringVar(&getKey, "getKey", "", "key values")
	makeGet.Flags().StringVar(&getEtcdIP, "ip", "", "etcd IP")

	makePut.Flags().StringVar(&putKey, "putKey", "", "set key to configure")
	makePut.Flags().StringVar(&putEtcdIp, "ip", "", "set ip for etcd client")
	makePut.Flags().StringVar(&putClusterId, "clusterid", "", "set clusterID")
	makePut.Flags().StringVar(&putRetry, "retry", "", "set retry count")

}

func HandleGetV2(key string, ip string)  {
	if key == "" || ip == "" {
		fmt.Println("key and IP are required to get values")
		os.Exit(1)
	}
	if key != "" && ip != "" {
		m := etcd.GetKV(key, ip)
		fmt.Println("map:", m)

	}

	
}
