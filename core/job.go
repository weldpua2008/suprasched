package core

import (
	"sync"
	"time"
)

// JobPhase defines the condition of pod.
type JobPhase string

// These are valid conditions of job.
const (
	JobUnassigned JobPhase = "UNASSIGNED"

	// JobPending means the job has been accepted by the system, but has not been started yet.
	// This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	JobPending JobPhase = "PENDING"

	// JobRunning means the job has been bound to a cluster and has been started.
	// At least one container is still running or is in the process of being restarted.
	JobRunning JobPhase = "RUNNING"

	// JobSucceeded means that the job voluntarily terminated
	// with an exit code of 0.
	JobSucceeded JobPhase = "SUCCEEDED"

	// JobFailed means that the jobs has terminated in a failure
	// (exited with a non-zero exit code or was stopped by the system).
	JobFailed JobPhase = "FAILED"

	JobCanceled JobPhase = "CANCELED"
	JobTimeout  JobPhase = "TIMEOUT"
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
	Priority       int64
	CreateAt       time.Time // When Job was created
	StartAt        time.Time // When command started
	LastActivityAt time.Time // When job metadata last changed
	Status         JobStatus // Represents information about the status
	TTR            uint64    // Time-to-run in Millisecond
	ClusterId      string    // Identification for ClusterId
	mu             sync.RWMutex
	// Note that this is calculated from dead Jobs. But those jobs are subject to
	// garbage collection.  This value will get capped at 5 by GC.
	RestartCount int32
	//exitError               error
	ExitCode int // Exit code
}

func (j Job) GetObjId() string {
	return string(j.UID)
}

// NewJob returns a new job
func NewJob(name string, ns Namespace, uid string) Job {
	return Job{
		Status: JobStatus{
			Phase: JobUnassigned,
		},
		CreateAt: time.Now(),
		ObjectMeta: ObjectMeta{
			Name:      name,
			Namespace: ns,
			UID:       UID(uid),
		},
		TypeMeta: TypeMeta{
			Kind:       "job",
			APIVersion: LatestApi,
		},
	}
}
