package model

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job public structure
type Job struct {
	Id             string    // Identificator for Job
	RunUID         string    // Running indentification
	ExtraRunUID    string    // Extra indentification
	Priority       int64     // Priority for a Job
	CreateAt       time.Time // When Job was created
	StartAt        time.Time // When command started
	LastActivityAt time.Time // When job metadata last changed
	Status         string    // Currentl status
	MaxAttempts    int       // Absoulute max num of attempts.
	MaxFails       int       // Absolute max number of failures.
	TTR            uint64    // Time-to-run in Millisecond
	ClusterId      string    // Identificator for ClusterId
	ClusterConfig  map[string]interface{}
	mu             sync.RWMutex
	exitError      error
	ExitCode       int // Exit code
	ctx            context.Context
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

// updateStatus job status
func (j *Job) updateStatus(status string) error {
	log.Trace(fmt.Sprintf("Job %s status %s -> %s", j.Id, j.Status, status))
	j.Status = status
	return nil
}

// IsTerminalStatus returns true if status is terminal:
// - Failed
// - Canceled
// - Successful
func IsTerminalStatus(status string) bool {
	switch status {
	case JOB_STATUS_ERROR, JOB_STATUS_CANCELED, JOB_STATUS_SUCCESS:
		log.Tracef("IsTerminalStatus %s true", status)
		return true
	}
	log.Tracef("IsTerminalStatus %s false", status)
	return false
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
		// stage := "jobs.cancel"
		// params := j.GetAPIParams(stage)
		// if err, result := DoApiCall(j.ctx, params, stage); err != nil {
		// 	log.Tracef("failed to update api, got: %s and %s", result, err)
		// }

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
		Status:         JOB_STATUS_PENDING,
		MaxFails:       1,
		MaxAttempts:    1,
		TTR:            0,
	}
}
