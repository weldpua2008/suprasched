package model

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cluster public structure
type Cluster struct {
	ClusterId      string // Identificator for Cluster
	ClusterPool    string // Identificator for Cluster Pool
	ClusterProfile string // Identificator for Cluster Profile

	ClusterConfig  map[string]interface{}
	CreateAt       time.Time // When cluster was created
	StartAt        time.Time // When cluster started
    StopAt        time.Time // When cluster started

	LastActivityAt time.Time // When cluster metadata last changed
	Status         string    // Currentl status
	// MaxAttempts    int       // Absoulute max num of attempts.
	// MaxFails       int       // Absolute max number of failures.
	// TTR            uint64    // Time-to-run in Millisecond
	mu             sync.RWMutex
	ExitCode       int // Exit code
	ctx            context.Context
	all            map[string]*Job
}

// NewCluster returns a new Clustec.
func NewCluster(clusterId string) *Cluster {
	return &Cluster{
		ClusterId: clusterId,
		all:       make(map[string]*Job),
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

// Len returns length of Clusters on cluster.
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

// Record fetch job by Cluster ID.
// Follows comma ok idiom
func (c *Cluster) Record(jid string) (*Job, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if rec, ok := c.all[jid]; ok {
		return rec, true
	}

	return nil, false
}

// ClusterStoreKey returns Cluster unique store key
func ClusterStoreKey(ClusterId string, ClusterPool string, ClusterProfile string) string {
	return fmt.Sprintf("%s:%s:%s", ClusterId, ClusterPool, ClusterProfile)
}

// StoreKey returns StoreKey
func (c *Cluster) StoreKey() string {
	return StoreKey(c.ClusterId, c.ClusterPool, c.ClusterProfile)
}
