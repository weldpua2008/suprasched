package core

import (
	"sync"
	"time"
)

// JobStatusType defines the condition of pod.
type JobStatusType string

// UID is a type that holds unique ID values, including UUIDs.
type UID string


// These are valid conditions of job.
const (
	JobStatusPending    JobStatusType = "PENDING"
	JobStatusInProgress JobStatusType = "RUNNING"
	JobStatusSuccess    JobStatusType = "SUCCESS"
	JobStatusError      JobStatusType = "ERROR"
	JobStatusCanceled   JobStatusType = "CANCELED"
	JobStatusTimeout    JobStatusType = "TIMEOUT"
)

// Job is a description of a job
type Job struct {
	TypeMeta
	ObjectMeta
	// If specified, the job's scheduling constraints
	// +optional
	Affinity *Affinity
	// Priority for a Job. The higher the value, the higher the priority.
	// +optional
	Priority                int64
	CreateAt                time.Time     // When Job was created
	StartAt                 time.Time     // When command started
	LastActivityAt          time.Time     // When job metadata last changed
	Status                  JobStatusType // Current status
	TTR                     uint64        // Time-to-run in Millisecond
	ClusterId               string        // Identification for ClusterId
	mu                      sync.RWMutex
	// Note that this is calculated from dead Jobs. But those jobs are subject to
	// garbage collection.  This value will get capped at 5 by GC.
	RestartCount int32
	//exitError               error
	ExitCode int // Exit code
}
