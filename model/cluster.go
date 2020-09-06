package model

import (
	"context"
	// "fmt"
	"sync"
	"time"
)

// Cluster public structure
type Cluster struct {
	ClusterId      string // Identificator for Job
	ClusterConfig  map[string]interface{}
	CreateAt       time.Time // When Job was created
	StartAt        time.Time // When command started
	LastActivityAt time.Time // When job metadata last changed
	Status         string    // Currentl status
	MaxAttempts    int       // Absoulute max num of attempts.
	MaxFails       int       // Absolute max number of failures.
	TTR            uint64    // Time-to-run in Millisecond
	mu             sync.RWMutex
	ExitCode       int // Exit code
	ctx            context.Context
	all            map[string]*Job
}

// NewCluster returns a new Clustec.
func NewCluster() *Cluster {
	return &Cluster{
		all: make(map[string]*Job),
	}
}

// Add a job.
// Returns false on duplicate or invalid job id.
func (c *Cluster) Add(rec *Job) bool {
	if rec == nil || rec.StoreKey() == "" {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.all[rec.StoreKey()]; ok {
		return false
	}

	c.all[rec.StoreKey()] = rec
	return true
}

// Len returns length of registry.
func (c *Cluster) Len() int {
	c.mu.RLock()
	t := len(c.all)
	c.mu.RUnlock()
	return t
}

// Delete a job by job ID.
// Return false if record does not exist.
func (c *Cluster) Delete(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.all[id]
	if !ok {
		return false
	}
	delete(c.all, id)
	return true
}

// Record fetch job by Job ID.
// Follows comma ok idiom
func (c *Cluster) Record(jid string) (*Job, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if rec, ok := c.all[jid]; ok {
		return rec, true
	}

	return nil, false
}
