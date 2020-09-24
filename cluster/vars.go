package cluster

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	model "github.com/weldpua2008/suprasched/model"
)

var (
	ErrNoSuitableClustersFetcher   = errors.New("No suitable ClustersFetcher found")
	ErrNoSuitableClustersDescriber = errors.New("No suitable ClustersDescriber found")
	ErrEmptyClusterId              = errors.New("Cluster Id is empty")
	ErrClusterIdIsNotValid         = errors.New("InvalidRequestException: Cluster id is not valid")
	// Internal
	clustersDescribed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_clusters_describe_total",
		Help: "The total number of described clusters",
	})
	clustersFetched = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_clusters_fetch_total",
		Help: "The total number of fetched clusters",
	})

	log                 = logrus.WithFields(logrus.Fields{"package": "cluster"})
	empty_clusters_chan = make(chan *model.Cluster, 1)
)
