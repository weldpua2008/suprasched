package config

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

const (
	TOPIC_CLUSTER_CREATED                = "cluster.created"
    TOPIC_CLUSTER_STARTING               = "cluster.starting"
	TOPIC_CLUSTER_BOOTSTRAPPING          = "cluster.bootstraping"
	TOPIC_CLUSTER_RUNNING                = "cluster.running"
	TOPIC_CLUSTER_WAITING                = "cluster.waiting"
	TOPIC_CLUSTER_TERMINATING            = "cluster.terminating"
	TOPIC_CLUSTER_TERMINATED             = "cluster.terminated"
	TOPIC_CLUSTER_TERMINATED_WITH_ERRORS = "cluster.terminated_with_errors"

)

var (
	// Bus is a ref to bus.Bus
	Bus *bus.Bus

	// Monoton is an instance of monoton.Monoton
	Monoton monoton.Monoton
)

// Init inits the app config
func InitEvenBus() {
	// configure id generator (it doesn't have to be monoton)
	node := uint64(1)
	initialTime := uint64(0)
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		log.Panic(err)
	}

	// init an id generator
	var idGenerator bus.Next = (*m).Next

	// create a new bus instance
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}

	// maybe register topics in here
	b.RegisterTopics(TOPIC_CLUSTER_CREATED, TOPIC_CLUSTER_STARTING, TOPIC_CLUSTER_BOOTSTRAPPING,
		TOPIC_CLUSTER_RUNNING, TOPIC_CLUSTER_WAITING, TOPIC_CLUSTER_TERMINATING,
        TOPIC_CLUSTER_TERMINATED, TOPIC_CLUSTER_TERMINATED_WITH_ERRORS, "job.pending",
		"job.created", "job.canceled", "job.running", "job.failed", "job.succeeded")

	Bus = b
	Monoton = *m
}
