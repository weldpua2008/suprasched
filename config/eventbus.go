package config

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

const (
	TOPIC_CLUSTER_CREATED       = "cluster.created"
	TOPIC_CLUSTER_STARTING      = "cluster.starting"
	TOPIC_CLUSTER_BOOTSTRAPPING = "cluster.bootstraping"
	TOPIC_CLUSTER_RUNNING       = "cluster.running"
	TOPIC_CLUSTER_WAITING       = "cluster.waiting"

	MATCHER_CLUSTER_TERMINATING = "cluster.term.*"

	TOPIC_CLUSTER_TERMINATING            = "cluster.terminating"
	TOPIC_CLUSTER_TERMINATED             = "cluster.terminated"
	TOPIC_CLUSTER_TERMINATED_WITH_ERRORS = "cluster.terminated_with_errors"

	TOPIC_JOB_CANCELED = "job.canceled"

	TOPIC_JOB_CREATED  = "job.created"
	TOPIC_JOB_STARTING = "job.starting"
	TOPIC_JOB_PENDING  = "job.pending"
	TOPIC_JOB_RUNNING  = "job.running"

	TOPIC_JOB_FAILED                 = "job.failed"
	TOPIC_JOB_SUCCEEDED              = "job.succeeded"
	TOPIC_JOB_TERMINATING            = "job.terminating"
	TOPIC_JOB_TERMINATED             = "job.terminated"
	TOPIC_JOB_TERMINATED_WITH_ERRORS = "job.terminated_with_errors"
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
		TOPIC_CLUSTER_TERMINATED, TOPIC_CLUSTER_TERMINATED_WITH_ERRORS, TOPIC_JOB_PENDING,
		TOPIC_JOB_CREATED, TOPIC_JOB_CANCELED, TOPIC_JOB_RUNNING, TOPIC_JOB_FAILED, TOPIC_JOB_SUCCEEDED)

	Bus = b
	Monoton = *m
}
