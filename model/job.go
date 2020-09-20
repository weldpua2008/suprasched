package model

import (
	"context"
	"fmt"
	utils "github.com/weldpua2008/suprasched/utils"
	"strings"
	"sync"
	"time"
)

// Job public structure
type Job struct {
	Id                      string    // Identificator for Job
	RunUID                  string    // Running indentification
	ExtraRunUID             string    // Extra indentification
	Priority                int64     // Priority for a Job
	CreateAt                time.Time // When Job was created
	StartAt                 time.Time // When command started
	LastActivityAt          time.Time // When job metadata last changed
	Status                  string    // Currentl status
	MaxAttempts             int       // Absoulute max num of attempts.
	MaxFails                int       // Absolute max number of failures.
	TTR                     uint64    // Time-to-run in Millisecond
	ClusterId               string    // Identificator for ClusterId
	ClusterType             string    // Identificator for Cluster Type
	ClusterStoreKey         string    // Cluster UUID for Registry
	PreviousClusterStoreKey string    // Previous Cluster UUID
	ClusterConfig           map[string]interface{}
	mu                      sync.RWMutex
	exitError               error
	ExitCode                int // Exit code
	ctx                     context.Context
	inTransition            bool // wether Job in some transaction
}

// StoreKey returns Job unique store key
func StoreKey(Id string, RunUID string, ExtraRunUID string) string {
	return fmt.Sprintf("%s:%s:%s", Id, RunUID, ExtraRunUID)
}

// StoreKey returns StoreKey
func (j *Job) StoreKey() string {
	return StoreKey(j.Id, j.RunUID, j.ExtraRunUID)
}

// GetStatus get job status.
func (j *Job) GetStatus() string {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.Status
}

// updatelastActivity for the Job
func (j *Job) updatelastActivity() {
	j.LastActivityAt = time.Now()
}

func (j *Job) IsInTransition() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if j.inTransition {
		return true
	}
	return false
}
func (j *Job) PutInTransition() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if !j.inTransition {
		j.inTransition = true
		return true
	}
	return false
}
func (j *Job) FinishTransition() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.inTransition {
		j.inTransition = false
		return true
	}
	return false
}

// UpdateStatus compare with cluster status string and updates.
// returns true if the cluster need update the status
func (j *Job) UpdateStatus(ext string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()

	if strings.ToLower(j.Status) != strings.ToLower(ext) {
		j.Status = ext
		return true
	}

	return false
}

func (j *Job) EventMetadata() map[string]string {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return map[string]string{
		"Id":          j.Id,
		"Status":      j.Status,
		"RunUID":      j.RunUID,
		"ExtraRunUID": j.ExtraRunUID,
		"ClusterId":   j.ClusterId,
		"ClusterType": j.ClusterType,
		"StoreKey":    j.StoreKey(),
	}
}

// updateStatus job status
func (j *Job) updateStatus(status string) error {
	log.Trace(fmt.Sprintf("Job %s status %s -> %s", j.Id, j.Status, status))
	j.Status = status
	return nil
}

func (j *Job) GetParams() map[string]interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()

	params := map[string]interface{}{
		"Id":          j.Id,
		"Status":      j.Status,
		"RunUID":      j.RunUID,
		"ExtraRunUID": j.ExtraRunUID,
		"ClusterId":   j.ClusterId,
	}
	return params
}

func (j *Job) ChangeClusterStoreKey(in string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.ClusterStoreKey != in {
		j.PreviousClusterStoreKey = j.ClusterStoreKey
		j.ClusterStoreKey = in
	}
}

func (j *Job) GetClusterStoreKey() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.ClusterStoreKey
}

func (j *Job) Cancel() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if !IsTerminalStatus(j.Status) {
		log.Tracef("Call Canceled for Job %s", j.Id)
		if errUpdate := j.updateStatus(JOB_STATUS_CANCELED); errUpdate != nil {
			log.Tracef("failed to change job %s status '%s' -> '%s'", j.Id, j.Status, JOB_STATUS_CANCELED)
		}
		j.updatelastActivity()

	} else {
		log.Trace(fmt.Sprintf("Job %s in terminal '%s' status ", j.Id, j.Status))
	}
	return nil
}

// NewJob return Job with defaults
func NewJob(id string) *Job {
	return &Job{
		Id:             id,
		CreateAt:       time.Now(),
		StartAt:        time.Now(),
		LastActivityAt: time.Now(),
		MaxFails:       1,
		MaxAttempts:    1,
		TTR:            0,
	}
}

func NewEmptyJob() *Job {
	return &Job{}
}
func NewJobFromMap(v map[string]interface{}) *Job {
	j := NewJob("")
	if found_val, ok := utils.GetFirstTimeFromMap(v, []string{"StartAt", "startAt", "StartDate", "startDate"}); ok {
		j.StartAt = found_val
	}

	if found_val, ok := utils.GetFirstStringFromMap(v, []string{"JobStatus", "jobStatus", "Job_Status", "Status", "status"}); ok {
		j.Status = found_val
	}
	if found_val, ok := utils.GetFirstStringFromMap(v, []string{"JobId", "jobId", "Job_ID", "Job_Id", "job_Id", "job_id"}); ok {
		j.Id = found_val
	}
	if found_val, ok := utils.GetFirstStringFromMap(v, []string{"ClusterId", "clusterId", "Cluster_Id", "Cluster_ID", "Cluster", "cluster", "clusterid"}); ok {
		j.ClusterId = found_val
	}

	if found_val, ok := utils.GetFirstStringFromMap(v, []string{"JobRunId", "jobRunId", "Job_RUN_ID", "Job_Run_Id", "job_run_id", "run_id",
		"run_uid", "RunId", "RunUID"}); ok {
		j.RunUID = found_val
	}

	if found_val, ok := utils.GetFirstStringFromMap(v, []string{"JobExtraRunId", "jobExtraRunId", "JOB_EXTRA_RUN_ID", "Job_Extra_Run_Id", "job_extra_run_id", "extra_run_id",
		"job_extra_run_uid", "extra_run_uid"}); ok {
		j.ExtraRunUID = found_val
	}

	return j
}
