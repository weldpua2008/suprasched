package handlers

import (
	"github.com/mustafaturan/bus"
	config "github.com/weldpua2008/suprasched/config"
	// model "github.com/weldpua2008/suprasched/model"
)

// EmptyCluster handler when the cluster has no jobs.
func EmptyCluster(e *bus.Event) {
	if cl, err := eventGetCLuster(e); err == nil {
		if cl.IsFree() {
			if config.ClusterRegistry.MarkFree(cl.StoreKey()) {
				t := cl.RefreshTimeout()
				log.Tracef("Cluster %v will be freed in %v [%v]", cl.ClusterId, t, cl.TimeOutAt)
				// log.Warningf("Cluster %v will be freed in %v [%v]", cl.ClusterId, t, cl.TimeOutAt)

			}

		} else {
			config.ClusterRegistry.UnMarkFree(cl.StoreKey())
		}
	}
}
