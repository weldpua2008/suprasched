// Copyright 2020 Valeriy Soloviov. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// license that can be found in the LICENSE file.

// This is the main package for the `supraworker` application.

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/weldpua2008/suprasched/cmd"
)

func main() {
	log := logrus.WithFields(logrus.Fields{"package": "main"})

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
