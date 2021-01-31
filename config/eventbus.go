package config

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

var (
	// Bus is a ref to bus.Bus
	Bus *bus.Bus
	// Monoton is an instance of monoton.Monoton
	Monoton    monoton.Monoton
	topicNames = []string{
		TOPIC_CLUSTER_CREATED, TOPIC_CLUSTER_STARTING, TOPIC_CLUSTER_BOOTSTRAPPING,
		TOPIC_CLUSTER_RUNNING, TOPIC_CLUSTER_WAITING, TOPIC_CLUSTER_TERMINATING, TOPIC_CLUSTER_IS_EMPTY,
		TOPIC_CLUSTER_TERMINATED, TOPIC_CLUSTER_TERMINATED_WITH_ERRORS, TOPIC_JOB_PENDING,
		TOPIC_CLUSTER_REFRESH_TIMEOUT,
		TOPIC_JOB_CREATED, TOPIC_JOB_CANCELED, TOPIC_JOB_RUNNING, TOPIC_JOB_FAILED, TOPIC_JOB_SUCCEEDED,
		TOPIC_JOB_SUCCESS, TOPIC_JOB_FORCE_TIMEOUT,
		"testing",
	}
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
	b.RegisterTopics(topicNames...)

	Bus = b
	Monoton = *m
}

func EvenBusTearDown() {
	Bus.DeregisterTopics(topicNames...)
}
