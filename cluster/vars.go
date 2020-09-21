package cluster

import (
	"errors"
	"github.com/sirupsen/logrus"
	model "github.com/weldpua2008/suprasched/model"
)

var (
	log                            = logrus.WithFields(logrus.Fields{"package": "cluster"})
	ErrNoSuitableClustersFetcher   = errors.New("No suitable ClustersFetcher found")
	ErrNoSuitableClustersDescriber = errors.New("No suitable ClustersDescriber found")
	ErrEmptyClusterId              = errors.New("Cluster Id is empty")
	// Internal
	empty_clusters_chan = make(chan *model.Cluster, 1)
)
