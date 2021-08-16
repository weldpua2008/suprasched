/*
Copyright 2021 The Suprasched Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package app
import (
	"errors"
	"fmt"
	"github.com/weldpua2008/supraworker/metrics"
	"github.com/weldpua2008/supraworker/model"
	"github.com/weldpua2008/supraworker/utils"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)


// Worker run goroutine for executing commands and reporting to apiserver and/or syncyng with etcd
// Note that a WaitGroup must be passed to functions by
// pointer.
// There are several scenarios for the Job execution:
//	1). Job execution finished with error/success [Regular flow]
//	2). Cancelled because of TTR [Timeout]
//	3). Cancelled by Job's Registry because of Cleanup process (TTR) [Cancel]
//	4). Cancelled when we fetch external API (cancellation information) [Cancel]

func Worker(id int, jobs <-chan *model.Job, wg *sync.WaitGroup) {
	workerId := fmt.Sprintf("worker-%d", id)
	logWorker := log.WithField("worker", workerId)
	// On return, notify the WaitGroup that we're done.
	defer func() {
		logWorker.Debugf("[FINISHED]")
		wg.Done()
	}()
	logWorker.Info("Starting joblet worker")
	for j := range jobs {
		j.AddToContext(utils.CtxWorkerIdKey, workerId)
		logJob := j.GetLogger()
		logJob.Tracef("New Job with TTR %v", time.Duration(j.TTR)*time.Millisecond)
		atomic.AddInt64(&NumActiveJobs, 1)
		errJobRun := j.Run()
		logJob.Tracef("Run() Finished")
		dur := time.Since(j.StartAt)
		switch {
		// Execution stopped by TTR
		case errors.Is(errJobRun, model.ErrJobTimeout):
			if errTimeout := j.TimeoutWithCancel(TimeoutJobsAfter5MinInTerminalState); errTimeout != nil {
				logJob.Tracef("[Timeout()] got: %v ", errTimeout)
			}
			metrics.JobsTimeout.Inc()
		case errors.Is(errJobRun, model.ErrJobCancelled):
			if errTimeout := j.Cancel(); errTimeout != nil {
				logJob.Tracef("[Cancel()] got: %v ", errTimeout)
			}
			metrics.JobsTimeout.Inc()
		case errJobRun == nil:
			if err := j.Finish(); err != nil {
				logJob.Debugf("finished in %v got %v", dur, err)
			} else {
				logJob.Debugf("finished in %v", dur)
			}

		default:
			if errFail := j.Failed(); errFail != nil {
				logJob.Tracef("[Failed()] got: %v ", errFail)
			}
			logJob.Infof("Failed with %s", errJobRun)
		}

		atomic.AddInt64(&NumActiveJobs, -1)
		atomic.AddInt64(&NumProcessedJobs, 1)
		runtime.Gosched()

	}
}
