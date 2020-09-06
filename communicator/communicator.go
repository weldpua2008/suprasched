package communicator

import (
	"context"
	"fmt"
)

// A Communicator is the interface used to communicate with APIs
// that will eventually return metadata. Communicators
// allow you to get information from remote APi, databases, etc.
//
// Communicators must be safe for concurrency, meaning multiple calls to
// any method may be called at the same time.
type Communicator interface {
	// Configured
	// Configured() bool
	// Configure Communicator
	Configure(map[string]interface{}) error
	// Fetch metadata from remote storage
	Fetch(context.Context, map[string]interface{}) ([]map[string]interface{}, error)
}

func GetCommunicator(communicator_type string) (error, Communicator) {
	switch communicator_type {
	case "http", "HTTP":
		return nil, NewRestCommunicator()
	default:
		return fmt.Errorf("Can't find sutable communicator for %s.\n", communicator_type), nil
	}
}
