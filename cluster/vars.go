package cluster

import (
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	log                          = logrus.WithFields(logrus.Fields{"package": "cluster"})
	ErrNoSuitableClustersFetcher = errors.New("No suitable ClustersFetcher found")
)
