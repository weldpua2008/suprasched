package handlers

import (
	"github.com/mustafaturan/bus"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
)

// Init registers all handlers.
func Init() {
	// Start("tracing", Trace, ".*")
	Start("cluster_termination", ClusterTermination, config.MATCHER_CLUSTER_TERMINATING)
	Start("cluster_is_empty", EmptyCluster, config.MATCHER_CLUSTER_IS_EMPTY)

}

// Deregister all handlers.
func Deregister() {
	// Stop("tracing")
	Stop("cluster_termination")
	Stop("cluster_is_empty")

}

// Start registers the handler
func Start(name string, f func(e *bus.Event), Matcher string) {
	b := config.Bus
	h := bus.Handler{Handle: f, Matcher: Matcher}
	b.RegisterHandler(name, &h)
	log.Tracef("Registered %v handler...", name)
}

//
func startTestingHandler(f func(e *bus.Event)) {
	config.JobsRegistry = model.NewRegistry()
	config.ClusterRegistry = model.NewClusterRegistry()

	config.InitEvenBus()
	name := "testing"
	Matcher := name
	Start(name, f, Matcher)
}

func stopTestingHandler() {
	config.EvenBusTearDown()
	name := "testing"
	Stop(name)
}

// Stop deregisters the handler
func Stop(name string) {
	defer log.Tracef("Deregistered %v handler...", name)
	b := config.Bus
	b.DeregisterHandler(name)
}
