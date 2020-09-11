package handlers

import (
	"github.com/mustafaturan/bus"
	config "github.com/weldpua2008/suprasched/config"
)

// Init registers all handlers.
func Init() {
	Start("tracing", Trace, ".*")
	Start("cluster_termination", ClusterTermination, config.MATCHER_CLUSTER_TERMINATING)

}

// Deregister all handlers.
func Deregister() {
    defer Stop("tracing")

	defer Stop("cluster_termination")

}

// Start registers the handler
func Start(name string, f func(e *bus.Event), Matcher string) {
	b := config.Bus
	h := bus.Handler{Handle: f, Matcher: Matcher}
	b.RegisterHandler(name, &h)
	log.Infof("Registered %v handler...", name)
}

// Stop deregisters the handler
func Stop(name string) {
	defer log.Infof("Deregistered %v handler...", name)

	b := config.Bus
	b.DeregisterHandler(name)
}
