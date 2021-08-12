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
	"github.com/spf13/viper"
	"os"
)

var OBJECTID string = "/ObjectID"
var CLUSTERID string = "/ClusterID"
var RETRY string = "/Retry"
var REVISION string = "/Revision"
var VALUE string = ".value"

// etcdCmd represents the etcd command
var etcdCmd = &cobra.Command{
	Use:   "etcd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("etcd called")
		//prefixKey, _ := cmd.Flags().GetString("key")
		clusterId, _ := cmd.Flags().GetString("clusterId")
		objectId, _ := cmd.Flags().GetString("objectId")
		revision, _ := cmd.Flags().GetString("revision")
		retry, _ := cmd.Flags().GetString("retry")

<<<<<<< HEAD



=======
>>>>>>> ad71868ab97ec80cde2074b95af826856e8d9962
		//fmt.Println("the getKey value ", getKey)
		//keyPath := "kv." + getKey
		//fmt.Println("the keyPath value ", keyPath)
		//getKey = viper.GetString(keyPath)
		//fmt.Println(viper.GetString(keyPath))
		//fmt.Println(viper.GetString("kv.key"))
		//fmt.Println(getKey)

		//getKey, _ := cmd.Flags().GetString("value")
		//fmt.Println(getKey)
		viper.SetConfigName("etcdConfigFile") //config file name
		viper.SetConfigType("yaml")
		viper.AddConfigPath("etcd/")
		viper.AutomaticEnv()
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Fatal error config file: defualt \n", err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("clusterId") {
			clusterIdPath := clusterId + CLUSTERID + VALUE
			clusterIdValue := viper.GetString(clusterIdPath)
			fmt.Println(clusterIdPath + ":" + clusterIdValue)
		}
		if cmd.Flags().Changed("objectId") {
			objectIdPath := objectId + OBJECTID + VALUE
			objectIdValue := viper.GetString(objectIdPath)
			fmt.Println(objectIdPath + ":" + objectIdValue)
		}
		if cmd.Flags().Changed("revision") {
			revisionPath := revision + REVISION + VALUE
			revisionValue := viper.GetString(revisionPath)
			fmt.Println(revisionPath + ":" + revisionValue)
		}
		if cmd.Flags().Changed("retry") {
			retryPath := retry + RETRY + VALUE
			retryValue := viper.GetString(retryPath)
			fmt.Println(retryPath + ":" + retryValue)
		}

		//getKey, _:= cmd.Flags().GetString("key")
		//fmt.Println(getKey)
		//obj1 := getKey + ".key"
		//result := viper.GetString(obj1)
		//fmt.Println("result is: ", result)
		//obj2 := viper.GetString("kv")
		//fmt.Println("the obj2: ", obj2)
		//obj3 := viper.GetString("kv.key")
		//fmt.Println("the obj3: ", obj3)

		//getKey()
		//getObject()
		//etcdkey := viper.GetString("kv.key")
		//fmt.Println(etcdkey)
		//etcdvalue := viper.GetString("kv.value")
		//fmt.Println(etcdvalue)

	},
}

func init() {
	rootCmd.AddCommand(etcdCmd)
	//etcdCmd.Flags().StringP("key", "k", "", "Get all values of key")
	etcdCmd.Flags().StringP("clusterId", "c", "", "get clusterId of key")
	etcdCmd.Flags().StringP("objectId", "o", "", "get objectId of key")
	etcdCmd.Flags().StringP("retry", "r", "", "get retry count of key")
	etcdCmd.Flags().StringP("revision", "n", "", "get revision of key")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// etcdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// etcdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type eventObject struct {
	KEY   string
	VALUE string
	//CLUSTER_ID string
	REVISION string
}

//func newObject(key string)  *eventObject {

//	obj := eventObject{KEY: key}
//	obj.VALUE = "null"
//	obj.REVISION = "null"

//	return &obj

//obj := eventObject{}
//obj.KEY := viper.GetString("kv.key")
//fmt.Println("Key: %v", obj.KEY)
//obj.VALUE := viper.GetString("kv.value")
//fmt.Println("Value: %v", obj.VALUE)
//obj.REVISION := viper.GetString("kv.revision")
//fmt.Println("Revision: %v", obj.REVISION)
//fmt.Println(obj)

//}
