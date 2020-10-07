package model

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	// config "github.com/weldpua2008/suprasched/config"
	// communicator "github.com/weldpua2008/suprasched/communicator"
)

// Cluster public structure
type Cluster struct {
	ClusterId      string // Identificator for Cluster
	ClusterPool    string // Identificator for Cluster Pool
	ClusterProfile string // Identificator for Cluster Profile
	ClusterRegion  string // Identificator for Cluster Region
	ClusterType    string // Identificator for Cluster Type

	ClusterConfig  map[string]interface{}
	CreateAt       time.Time // When cluster was created
	StartAt        time.Time // When cluster started
	StopAt         time.Time // When cluster started
	MaxCapacity    int       // Maximum Jobs per cluster
	TimeOutStartAt time.Time // Initial time for timeout
	TimeOutAt      time.Time // Initial time for timeout

	TimeOutDuration    time.Duration // Duration after that Cluster is timed out
	LastSyncedAt       time.Time     // When cluster metadata last changed
	LastSyncedDuration time.Duration

	LastActivityAt time.Time // When cluster metadata last changed
	PreviousStatus string    // Previous Status
	Status         string    // Currentl status
	// MaxAttempts    int       // Absoulute max num of attempts.
	// MaxFails       int       // Absolute max number of failures.
	// TTR            uint64    // Time-to-run in Millisecond
	mu           sync.RWMutex
	ExitCode     int // Exit code
	ctx          context.Context
	all          map[string]*Job
	inTransition bool // if cluster is modified
}

// NewCluster returns a new Clustec.
func NewCluster(clusterId string) *Cluster {
	return &Cluster{
		ClusterId:          clusterId,
		all:                make(map[string]*Job),
		ClusterType:        CLUSTER_TYPE_EMR,
		TimeOutDuration:    time.Minute * 120,
		LastSyncedAt:       time.Now(),
		LastSyncedDuration: time.Second * 30,
	}
}

func (c *Cluster) RefreshTimeout() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TimeOutStartAt = time.Now()
	c.TimeOutAt = c.TimeOutStartAt.Add(c.TimeOutDuration)

	return c.TimeOutDuration
}

func (c *Cluster) IsTimeout() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Now().After(c.TimeOutAt)
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
	params["ClusterStatus"] = c.Status
	params["Status"] = c.Status

	return params
}

func (c *Cluster) GetParamsMapString() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	params := make(map[string]string)
	params["ClusterId"] = c.ClusterId
	params["ClusterPool"] = c.ClusterPool
	params["ClusterProfile"] = c.ClusterProfile
	params["ClusterRegion"] = c.ClusterRegion
	params["ClusterType"] = c.ClusterType
	params["ClusterStatus"] = c.Status
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

// Len returns length of Jobs on cluster.
func (c *Cluster) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.all)
}

// IsFull returns true if cluster reched maximum capacity for Jobs.
func (c *Cluster) IsFull() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	n := 0
	for _, j := range c.all {
		if !IsTerminalStatus(j.GetStatus()) {
			n += 1
		}
	}
	return n >= c.MaxCapacity
}

// IsEmpty returns true if cluster has no Jobs.
func (c *Cluster) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.all) < 1
}

// IsFree returns true if cluster has no active Jobs.
func (c *Cluster) IsFree() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, j := range c.all {
		if !IsTerminalStatus(j.GetStatus()) {
			return false
		}
	}
	return true
}

// Delete a job by job ID.
// Return false if record does not exist.
func (c *Cluster) Delete(jid string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.all[jid]
	if !ok {
		return false
	}
	delete(c.all, jid)
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

// Record fetch job by Cluster ID.
// Follows comma ok idiom
func (c *Cluster) All() []*Job {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ret := make([]*Job, len(c.all))
	for _, v := range c.all {
		ret = append(ret, v)
	}
	return ret
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
		"StoreKey":       c.StoreKey(),
	}
}

// func (c *Cluster) GetCommunicators()  ([]Communicator, error) {
// 	c.mu.RLock()
// 	defer c.mu.RUnlock()
//     comms, err:= communicator.GetCommunicatorsFromSection(fmt.Sprintf("%v.%v.%v",
//         config.CFG_PREFIX_CLUSTER,config.CFG_PREFIX_UPDATE, c.ClusterType ))
//
// 	return comms, err
// }

// UseExternaleStatus compare with another cluster status.
// returns true if the cluster need update the status
func (c *Cluster) UseExternaleStatus(ext *Cluster) bool {
	if ext.StoreKey() != c.StoreKey() {
		c.mu.RLock()
		defer c.mu.RUnlock()
		ext.mu.RLock()
		defer ext.mu.RUnlock()

	}
	if GetClusterStatusWeight(ext.Status) > GetClusterStatusWeight(c.Status) {
		return true
	} else if strings.ToLower(c.Status) != strings.ToLower(ext.Status) {
		return true
	}

	return false
}

// UseExternaleStatusString compare with cluster status string.
// returns true if the cluster need update the status
func (c *Cluster) UseExternaleStatusString(ext string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if GetClusterStatusWeight(ext) > GetClusterStatusWeight(c.Status) {
		return true
	} else if strings.ToLower(c.Status) != strings.ToLower(ext) {
		return true
	}

	return false
}

func (c *Cluster) IsInTransition() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.inTransition {
		return true
	}
	return false
}

func (c *Cluster) PutInTransition() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.inTransition {
		c.inTransition = true
		return true
	}
	return false
}
func (c *Cluster) FinishTransition() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.inTransition {
		c.inTransition = false
		return true
	}
	return false
}

func (c *Cluster) SyncedWithExternalAPI() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastSyncedAt = time.Now().Add(c.LastSyncedDuration)
}

// UpdateStatus compare with cluster status string and updates.
// returns true if the cluster need update the status
func (c *Cluster) UpdateStatus(ext string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if GetClusterStatusWeight(ext) > GetClusterStatusWeight(c.Status) {
		c.updateStatus(ext)
		return true
	}
	// else if strings.ToLower(c.Status) != strings.ToLower(ext) {
	// 	c.updateStatus(ext)
	// 	return true
	// }

	return false
}

// updateStatus cluster status
func (c *Cluster) updateStatus(ext string) error {
	c.PreviousStatus = c.Status
	c.Status = ext
	log.Trace(fmt.Sprintf("Cluster %s status %s -> %s", c.ClusterId, c.PreviousStatus, c.Status))
	return nil
}

func (c *Cluster) UpdateStatusInTransition(ext string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.inTransition {
		if GetClusterStatusWeight(ext) > GetClusterStatusWeight(c.Status) {
			c.updateStatus(ext)
			c.inTransition = true
			return true
		} else if c.Status != ext {
			c.updateStatus(ext)
			c.inTransition = true
			return true
		}
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
