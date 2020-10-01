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

	// Internal variables
	clusterStatuses = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "suprasched",
			Subsystem: "clusters",
			Name:      "statuses",
			Help:      "Number of cluster statuses, partitioned by Profile and Type.",
		},
		[]string{
			// Which profile is used?
			"profile",
			// Of what type is the cluster?
			"type",
			// What is the cluster status?
			"status",
		},
	)

	apiCallsStatistics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "suprasched",
			Subsystem: "api",
			Name:      "calls",
			Help:      "Number of API calls to 3rd party API partitioned by Type.",
		},
		[]string{
			// For example Amazon
			"provider",
			// Which profile is used?
			"profile",
			// Of what type is the request?
			"type",
			// What is the Operation?
			"operation",
		},
	)

	// clusterStatuses = prometheus.NewCounterVec(
	// 	prometheus.CounterOpts{
	// 		Name: "suprasched_clusters_statuses",
	// 		Help: "Number of hard-disk errors.",
	// 	},
	// 	[]string{"device"},
	// )

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
