package handlers

import (
	"github.com/mustafaturan/bus"
	// config "github.com/weldpua2008/suprasched/config"
	// model "github.com/weldpua2008/suprasched/model"
)

// RefreshTimeoutCluster handler when the cluster has jobs.
func RefreshTimeoutCluster(e *bus.Event) {
	if cl, err := eventGetCLuster(e); err == nil {
        cl.UpdateJobsLastActivity()
		cl.RefreshTimeout()
	}
}
