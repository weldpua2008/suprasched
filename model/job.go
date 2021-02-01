package model

import (
	"fmt"
	utils "github.com/weldpua2008/suprasched/utils"
	// communicator "github.com/weldpua2008/suprasched/communicator"
	// config "github.com/weldpua2008/suprasched/config"

	"strings"
	"sync"
	"time"
)

// Job public structure
type Job struct {
	Id                      string    // Identification for Job
	RunUID                  string    // Running Identification
	ExtraRunUID             string    // Extra Identification
	Priority                int64     // Priority for a Job
	CreateAt                time.Time // When Job was created
	StartAt                 time.Time // When command started
	LastActivityAt          time.Time // When job metadata last changed
	PreviousStatus          string    // Previous Status
	Status                  string    // Current status
	MaxAttempts             int       // Absolute max num of attempts.
	MaxFails                int       // Absolute max number of failures.
	TTR                     uint64    // Time-to-run in Millisecond
	ClusterId               string    // Identification for ClusterId
	ClusterType             string    // Identification for Cluster Type
	ClusterStoreKey         string    // Cluster UUID for Registry
	PreviousClusterStoreKey string    // Previous Cluster UUID
	ClusterConfig           map[string]interface{}
	mu                      sync.RWMutex
	//exitError               error
	ExitCode int // Exit code
	//ctx                     context.Context
	inTransition    bool // whether Job in some transaction
	ExtraSendParams map[string]string
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
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.Status
}

// updatesLastActivity for the Job
func (j *Job) updatesLastActivity() {
	j.LastActivityAt = time.Now()
}

func (j *Job) IsInTransition() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.inTransition
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

	if !strings.EqualFold(j.Status, ext) {
		if err := j.updateStatus(ext); err != nil {
			log.Tracef("Failed to update %v %v => %v", j.Id, j.Status, ext)
		}
		return true
	}

	return false
}

func (j *Job) EventMetadata() map[string]string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	previousStatus := j.PreviousStatus
	if len(previousStatus) < 1 {
		previousStatus = j.Status
	}

	return map[string]string{
		"Id":                      j.Id,
		"JobId":                   j.Id,
		"Status":                  j.Status,
		"JobStatus":               j.Status,
		"PreviousStatus":          previousStatus,
		"JobPreviousStatus":       previousStatus,
		"RunUID":                  j.RunUID,
		"ExtraRunUID":             j.ExtraRunUID,
		"ClusterId":               j.ClusterId,
		"ClusterType":             j.ClusterType,
		"StoreKey":                j.StoreKey(),
		"PreviousClusterStoreKey": j.PreviousClusterStoreKey,
		"ClusterStoreKey":         j.ClusterStoreKey,
	}
}

// updateStatus job status
func (j *Job) updateStatus(status string) error {
	clusterId := j.ClusterStoreKey

	if len(j.ClusterStoreKey) > 0 {
		clusterId = j.ClusterId
	}
	log.Trace(fmt.Sprintf("Job %s status %s -> %s [%v %v]", j.Id, j.Status, status, clusterId, j.ClusterType))
	j.PreviousStatus = j.Status
	j.Status = status
	j.updatesLastActivity()
	return nil
}

func (j *Job) GetParams() map[string]interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	previousStatus := j.PreviousStatus
	if len(previousStatus) < 1 {
		previousStatus = j.Status
	}
	params := map[string]interface{}{
		"Id":                      j.Id,
		"JobId":                   j.Id,
		"PreviousStatus":          previousStatus,
		"JobPreviousStatus":       previousStatus,
		"Status":                  j.Status,
		"RunUID":                  j.RunUID,
		"ExtraRunUID":             j.ExtraRunUID,
		"ClusterId":               j.ClusterId,
		"ClusterType":             j.ClusterType,
		"PreviousClusterStoreKey": j.PreviousClusterStoreKey,
		"ClusterStoreKey":         j.ClusterStoreKey,
		"StoreKey":                j.StoreKey(),
	}
	if j.ExtraSendParams != nil {
		for k, v := range j.ExtraSendParams {
			if _, ok := params[k]; !ok {
				params[k] = v
			}
		}

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
		j.updatesLastActivity()

	} else {
		log.Trace(fmt.Sprintf("Job %s in terminal '%s' status ", j.Id, j.Status))
	}
	return nil
}

//
//
// func (j *Job) communicator(section string) error {
// 	return nil
// }

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
	if foundVal, ok := utils.GetFirstTimeFromMap(v, []string{"StartAt", "startAt", "StartDate", "startDate"}); ok {
		j.StartAt = foundVal
	}

	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"JobStatus", "jobStatus", "Job_Status", "Status", "status"}); ok {
		j.Status = foundVal
	}
	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"previousJobStatus", "PreviousJobStatus", "Previous_Job_Status", "PreviousStatus", "previousstatus"}); ok {
		j.PreviousStatus = foundVal
	}
	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"JobId", "jobId", "Job_ID", "Job_Id", "job_Id", "job_id", "Id", "id"}); ok {
		j.Id = foundVal
	}
	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"ClusterId", "Clusterid", "clusterId", "Cluster_Id", "Cluster_ID", "Cluster", "cluster", "clusterid"}); ok {
		j.ClusterId = foundVal
	}

	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"JobRunId", "jobRunId", "Job_RUN_ID", "Job_Run_Id", "job_run_id", "run_id",
		"run_uid", "RunId", "RunUID"}); ok {
		j.RunUID = foundVal
	}

	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"JobExtraRunId", "jobExtraRunId", "JOB_EXTRA_RUN_ID", "Job_Extra_Run_Id", "job_extra_run_id", "extra_run_id",
		"job_extra_run_uid", "extra_run_uid"}); ok {
		j.ExtraRunUID = foundVal
	}
	ClusterPool := ""
	ClusterProfile := ""
	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"clusterPool", "ClusterPool", "Pool", "pool"}); ok {
		ClusterPool = foundVal
	}
	if foundVal, ok := utils.GetFirstStringFromMap(v, []string{"clusterProfile", "ClusterProfile", "AWSProfile", "AWS_Profile", "AWS_PROFILE"}); ok {
		ClusterProfile = foundVal
	}

	j.ClusterStoreKey = ClusterStoreKey(j.ClusterId, ClusterPool, ClusterProfile)

	return j
}
