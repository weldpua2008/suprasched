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

	clusterIdsAreNotValid = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "suprasched_cluster_ids_are_not_valid_total",
		Help: "The total number of not valid clusters Ids",
	})

	reqClustersDescribed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_req_clusters_describe_total",
		Help: "The total number of query for describe clusters",
	})

	reqClustersFailDescribed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_req_clusters_fail_describe_total",
		Help: "The total number of query for fail to describe clusters",
	})

	clustersFetched = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_clusters_fetch_total",
		Help: "The total number of fetched clusters",
	})

	log                 = logrus.WithFields(logrus.Fields{"package": "cluster"})
	empty_clusters_chan = make(chan *model.Cluster, 1)
)
