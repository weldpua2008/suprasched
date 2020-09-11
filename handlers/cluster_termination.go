package handlers

import (
	// config "github.com/weldpua2008/suprasched/config"
	"github.com/mustafaturan/bus"
)

// ClusterTermination handler
// 1. Send Cluster termination (external)
// 2. Find new cluster / Create a new cluster
// 3. Lock Cluster & sub jobs
// 4. Reassign Jobs (ext)
// 5. Unlock Cluster & sub jobs
func ClusterTermination(e *bus.Event) {
	log.Tracef("Event for %s: %+v", e.Topic, e)
}
