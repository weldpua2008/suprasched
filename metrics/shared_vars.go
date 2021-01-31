package metrics

import (
	// "sync"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventBusMessageProcessed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "suprasched",
		Subsystem: "eventbus",
		Name:      "latency_ns",
		Help:      "The latency distribution of messages processed by Eventbus",
	},
		[]string{"topic", "type"},
	)
	FetchMetadataLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "suprasched",
		Subsystem: "functions_tracing",
		Name:      "latency_ns",
		Help:      "The latency distribution of functions processed",
	},
		[]string{"function", "type"},
	)

	ApiCallsStatistics = promauto.NewCounterVec(
		prometheus.CounterOpts{
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

	ReqClustersTerminated = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "suprasched",
		Subsystem: "req_clusters",
		Name:      "terminated_total",
		Help:      "Number of API calls for Cluster termination to 3rd party API partitioned by Type.",
	},
		[]string{
			// For example Amazon
			"provider",
			// Which profile is used?
			"profile",
			// Of what type is the request?
			"type",
		},
	)
)
