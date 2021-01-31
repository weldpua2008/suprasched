package handlers

import (
	"github.com/mustafaturan/bus"
	// config "github.com/weldpua2008/suprasched/config"
	// model "github.com/weldpua2008/suprasched/model"
	//"time"
)

// CancelTimeoutJobs handler when the job is timeout.
func CancelTimeoutJobs(e *bus.Event) {
	if j, err := eventGetJob(e); err == nil {
		j.PutInTransition()
		defer j.FinishTransition()
		// log.Tracef("Job %v will terminated %v", j.StoreKey(), time.Duration(time.Now().Sub(j.CreateAt).Seconds())*time.Second)

	}
}
