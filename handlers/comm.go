package handlers

import (
	"github.com/mustafaturan/bus"
	config "github.com/weldpua2008/suprasched/config"
	metrics "github.com/weldpua2008/suprasched/metrics"
	model "github.com/weldpua2008/suprasched/model"
	"time"
)

// Init registers all handlers.
func Init() {
	// Start("tracing", Trace, ".*")
	Start("cluster_termination", ClusterTermination, config.MATCHER_CLUSTER_TERMINATING)
	Start("cluster_is_empty", EmptyCluster, config.MATCHER_CLUSTER_IS_EMPTY)
	Start("cluster_refresh_timeout", RefreshTimeoutCluster, config.MATCHER_CLUSTER_REFRESH_TIMEOUT)
	Start("job_force_timeout", CancelTimeoutJobs, config.MATCHER_JOB_FORCE_TIMEOUT)

}

func wrapMetrics(f func(e *bus.Event)) func(e *bus.Event) {
	return func(e *bus.Event) {
		start := time.Now()
		defer metrics.EventBusMessageProcessed.WithLabelValues(e.Topic,
			"wraped").Observe(float64(time.Now().Sub(start).Nanoseconds()))
		f(e)
	}
}

// Deregister all handlers.
func Deregister() {
	// Stop("tracing")
	Stop("cluster_termination")
	Stop("cluster_is_empty")
	Stop("cluster_refresh_timeout")
	Stop("job_force_timeout")

}

// Start registers the handler
func Start(name string, f func(e *bus.Event), Matcher string) {
	b := config.Bus
	h := bus.Handler{Handle: wrapMetrics(f), Matcher: Matcher}
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

// Stop is the deregister handler
func Stop(name string) {
	defer log.Tracef("Deregistered %v handler...", name)
	b := config.Bus
	b.DeregisterHandler(name)
}
