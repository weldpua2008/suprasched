package config

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

// Bus is a ref to bus.Bus
var Bus *bus.Bus

// Monoton is an instance of monoton.Monoton
var Monoton monoton.Monoton

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
	b.RegisterTopics("cluster.created", "cluster.waiting", "cluster.running",
		"cluster.terminated", "job.pending",
		"job.created", "job.canceled", "job.running", "job.failed", "job.succeeded")

	Bus = b
	Monoton = *m
}
