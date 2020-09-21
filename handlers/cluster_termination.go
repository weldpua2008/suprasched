package handlers

import (
	// "fmt"
	// "time"
	// "context"
	"github.com/mustafaturan/bus"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	// communicator "github.com/weldpua2008/suprasched/communicator"
)

// ClusterTermination handler
// Simple Implementation:
// 1. Send Cluster termination (external)
// 2. Cancel all Jobs (external)
//
// TODO:
//
// 1. Send Cluster termination (external)
// 2. Find new cluster / Create a new cluster
// 3. Lock Cluster & sub jobs
// 4. Reassign Jobs (ext)
// 5. Unlock Cluster & sub jobs
func ClusterTermination(e *bus.Event) {
	if err := eventCLusterRunComms(e); err != nil {
		log.Tracef("%v", err)
	}
	if cl, err := eventGetCLuster(e); err == nil {
		cl.PutInTransition()
		defer cl.FinishTransition()
		eData := eventDataStringMapString(e)
		jobs := cl.All()
		for _, j := range jobs {
			if j == nil || model.IsTerminalStatus(j.GetStatus()) {
				continue
			}

			if err := eventJobRunComms(j, eData); err != nil {
				log.Tracef("%v", err)
			} else {
				log.Tracef("Canceled Job %v ", j)
			}
			// config.JobsRegistry.Delete(j.StoreKey())
			// cl.Delete(j.StoreKey())
			log.Tracef("Removed Job %v ", j.StoreKey())
			// j = nil
		}
		config.ClusterRegistry.Delete(cl.StoreKey())
		log.Infof("Terminated all jobs on %v, Event in %s: %+v", cl.ClusterId, e.Topic, e)
		// cl = nil
	}
}
