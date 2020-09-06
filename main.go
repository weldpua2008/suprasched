// Copyright 2020 Valeriy Soloviov. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// license that can be found in the LICENSE file.

// This is the main package for the `supraworker` application.

package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/weldpua2008/suprasched/cluster"
	"github.com/weldpua2008/suprasched/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logrus.WithFields(logrus.Fields{"package": "main"})

	cl := cluster.NewDescribeEMR()
	params := make(map[string]interface{})
	params["ClusterID"] = "j-3JTIEH8MDWQ21"
	params["aws_profile"] = "bi-use1"
	params["ctx"] = ctx

	clusterId, _ := cl.DescribeCluster(params)
	log.Infof("%s", clusterId)

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
