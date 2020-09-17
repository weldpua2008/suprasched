package model

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Cluster public structure
type Cluster struct {
	ClusterId      string // Identificator for Cluster
	ClusterPool    string // Identificator for Cluster Pool
	ClusterProfile string // Identificator for Cluster Profile
	ClusterRegion  string // Identificator for Cluster Region
	ClusterType    string // Identificator for Cluster Type

	ClusterConfig map[string]interface{}
	CreateAt      time.Time // When cluster was created
	StartAt       time.Time // When cluster started
	StopAt        time.Time // When cluster started

	LastActivityAt time.Time // When cluster metadata last changed
	Status         string    // Currentl status
	// MaxAttempts    int       // Absoulute max num of attempts.
	// MaxFails       int       // Absolute max number of failures.
	// TTR            uint64    // Time-to-run in Millisecond
	mu       sync.RWMutex
	ExitCode int // Exit code
	ctx      context.Context
	all      map[string]*Job
}

// NewCluster returns a new Clustec.
func NewCluster(clusterId string) *Cluster {
	return &Cluster{
		ClusterId:   clusterId,
		all:         make(map[string]*Job),
		ClusterType: CLUSTER_TYPE_EMR,
	}
}

func (c *Cluster) GetParams() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	params := make(map[string]interface{})
	params["ClusterId"] = c.ClusterId
	params["ClusterPool"] = c.ClusterPool
	params["ClusterProfile"] = c.ClusterProfile
	params["ClusterRegion"] = c.ClusterRegion
	params["ClusterType"] = c.ClusterType
	params["Status"] = c.Status
	return params
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

// EventMetadata.
func (c *Cluster) EventMetadata() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]string{
		"ClusterStatus":  c.Status,
		"Status":         c.Status,
		"ClusterPool":    c.ClusterPool,
		"ClusterId":      c.ClusterId,
		"ClusterProfile": c.ClusterProfile,
		"ClusterRegion":  c.ClusterRegion,
		"ClusterType":    c.ClusterType,
	}
}

// UseExternaleStatus compare with another cluster status.
// returns true if the cluster need update the status
func (c *Cluster) UseExternaleStatus(ext *Cluster) bool {
	if ext.StoreKey() != c.StoreKey() {
		c.mu.RLock()
		defer c.mu.RUnlock()
	}
	if GetClusterStatusWeight(ext.Status) > GetClusterStatusWeight(c.Status) {
		return true
	} else if c.Status != ext.Status {
		return true
	}

	return false
}

// UseExternaleStatusString compare with cluster status string.
// returns true if the cluster need update the status
func (c *Cluster) UseExternaleStatusString(ext string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if GetClusterStatusWeight(ext) > GetClusterStatusWeight(c.Status) {
		return true
	} else if c.Status != ext {
		return true
	}

	return false
}

// UpdateStatus compare with cluster status string and updates.
// returns true if the cluster need update the status
func (c *Cluster) UpdateStatus(ext string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if GetClusterStatusWeight(ext) > GetClusterStatusWeight(c.Status) {
		c.Status = ext
		return true
	} else if c.Status != ext {
		c.Status = ext
		return true
	}

	return false
}

// ClusterStoreKey returns Cluster unique store key
func ClusterStoreKey(ClusterId string, ClusterPool string, ClusterProfile string) string {
	return fmt.Sprintf("%s:%s:%s", ClusterId, ClusterPool, ClusterProfile)
}

// StoreKey returns StoreKey
func (c *Cluster) StoreKey() string {
	return StoreKey(c.ClusterId, c.ClusterPool, c.ClusterProfile)
}

// GetClusterStatusWeight for Cluster Status
func GetClusterStatusWeight(status string) int {

	switch strings.ToLower(status) {

	case CLUSTER_STATUS_STARTING:
		return 1
	case CLUSTER_STATUS_BOOTSTRAPPING:
		return 2
	case CLUSTER_STATUS_RUNNING:
		return 3
	case CLUSTER_STATUS_WAITING:
		return 3
	case CLUSTER_STATUS_TERMINATING:
		return 4
	case CLUSTER_STATUS_TERMINATED:
		return 5
	case CLUSTER_STATUS_TERMINATED_WITH_ERRORS:
		return 5
	}
	return 0
}
