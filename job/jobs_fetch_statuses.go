package job

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	// "strings"
	communicator "github.com/weldpua2008/suprasched/communicator"
	"time"

	model "github.com/weldpua2008/suprasched/model"
)

var (
	log = logrus.WithFields(logrus.Fields{"package": "jobs"})
	// Registry for the Jobs
	JobsRegistry = model.NewRegistry()
)

// ApiJobRequest is struct for new jobs
type ApiJobRequest struct {
	JobStatus string `json:"job_status"`
	Limit     int64  `json:"limit"`
}

// An ApiJobResponse represents a Job response.
// Example response
// {
//   "job_id": "dbd618f0-a878-e477-7234-2ef24cb85ef6",
//   "jobStatus": "RUNNING",
//   "has_error": false,
//   "error_msg": "",
//   "run_uid": "0f37a129-eb52-96a7-198b-44515220547e",
//   "job_name": "Untitled",
//   "cmd": "su  - hadoop -c 'hdfs ls ''",
//   "parameters": [],
//   "createDate": "1583414512",
//   "lastUpdated": "1583415483",
//   "stopDate": "1586092912",
//   "extra_run_id": "scheduled__2020-03-05T09:21:40.961391+00:00"
// }
type ApiJobResponse struct {
	JobId       string   `json:"job_id"`
	JobStatus   string   `json:"jobStatus"`
	JobName     string   `json:"job_name"`
	RunUID      string   `json:"run_uid"`
	ExtraRunUID string   `json:"extra_run_id"`
	CMD         string   `json:"cmd"`
	Parameters  []string `json:"parameters"`
	CreateDate  string   `json:"createDate"`
	LastUpdated string   `json:"lastUpdated"`
	StopDate    string   `json:"stopDate"`
}

// NewApiJobRequest prepare struct for Jobs for execution request
func NewApiJobRequest() *ApiJobRequest {
	return &ApiJobRequest{
		JobStatus: "PENDING",
		Limit:     5,
	}
}

// StartFetchJobs goroutine for getting jobs statuses with internal
// exists on kill
func StartFetchJobs(ctx context.Context, comm communicator.Communicator, jobs chan *model.Job, interval time.Duration) error {

	return nil
}

// GracefullShutdown cancel all running jobs
// returns error in case any job failed to cancel
func GracefullShutdown(jobs <-chan *model.Job) bool {
	// empty jobs channel
	if len(jobs) > 0 {
		log.Trace(fmt.Sprintf("jobs chan still has size %v, empty it", len(jobs)))
		for len(jobs) > 0 {
			<-jobs
		}
	}
	JobsRegistry.GracefullShutdown()
	if JobsRegistry.Len() > 0 {
		log.Trace(fmt.Sprintf("GracefullShutdown failed, '%v' jobs left ", JobsRegistry.Len()))
		return false
	}
	return true

}
