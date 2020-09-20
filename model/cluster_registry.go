package model

import (
	"sync"
	// "time"
)

// NewClusterRegistry returns a new ClusterRegistry.
func NewClusterRegistry() *ClusterRegistry {
	return &ClusterRegistry{
		all:    make(map[string]*Cluster),
		byType: make(map[string][]*Cluster),
	}
}

// ClusterRegistry holds all Cluster Records.
type ClusterRegistry struct {
	all    map[string]*Cluster
	mu     sync.RWMutex
	byType map[string][]*Cluster
}

// Add a cluster.
// Returns false on duplicate or invalid cluster id.
func (r *ClusterRegistry) Add(rec *Cluster) bool {
	if rec == nil || rec.StoreKey() == "" {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.all[rec.StoreKey()]; ok {
		r.byType[rec.ClusterType] = append(r.byType[rec.ClusterType], rec)
		return false
	}

	r.all[rec.StoreKey()] = rec
	return true
}

// Len returns length of registry.
func (r *ClusterRegistry) Len() int {
	r.mu.RLock()
	c := len(r.all)
	r.mu.RUnlock()
	return c
}

// Delete a cluster by cluster ID.
// Return false if record does not exist.
func (r *ClusterRegistry) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.all[id]
	if !ok {
		return false
	}
	delete(r.all, id)
	idx := -1
	for i, _ := range r.byType[rec.ClusterType] {
		if r.byType[rec.ClusterType][idx] == rec {
			idx = i
			break
		}
	}
	if idx > -1 {
		copy(r.byType[rec.ClusterType][idx:], r.byType[rec.ClusterType][idx+1:])
		r.byType[rec.ClusterType] = r.byType[rec.ClusterType][:len(r.byType[rec.ClusterType])-1]

	}

	return true
}

// Filter a cluster by cluster types.
func (r *ClusterRegistry) Filter(cluster_types []string) []*Cluster {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ret := make([]*Cluster, 0)
	for _, clustertype := range cluster_types {
		if val, ok := r.byType[clustertype]; ok {
			ret = append(ret, val...)
		}
	}
	return ret
}

func (r *ClusterRegistry) AllEmpty() []*Cluster {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ret := make([]*Cluster, 0)
	for _, val := range r.all {
		if val.IsEmpty() {
			ret = append(ret, val)
		}
	}
	return ret
}

// // Cleanup process for the registry with batch only locked.
// // Return number of cleaned clusters.
// func (r *ClusterRegistry) Cleanup() (num int) {
// 	n := r.Len()
// 	slice := make([]string, n)
// 	i := 0
// 	for k, _ := range r.all {
// 		if i > n {
// 			slice = append(slice, k)
// 		} else {
// 			slice[i] = k
//
// 		}
// 		i++
// 	}
//
// 	batch := 20
// 	for i := 0; i < len(slice); i += batch {
// 		j := i + batch
// 		if j > len(slice) {
// 			j = len(slice)
// 		}
//
// 		// fmt.Println(slice[i:j]) // Process the batch.
// 		numBatch := r.CleanupBatch(slice[i:j])
// 		num += numBatch
// 	}
//
// 	return num
// }

// // CleanupBatch by cluster TTR.
// // Return number of cleaned clusters.
// func (r *ClusterRegistry) CleanupBatch(slice []string) (num int) {
// 	now := time.Now()
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	// for k, v := range r.all {
// 	for _, k := range slice {
// 		if v, ok := r.all[k]; ok {
// 			end := v.StartAt.Add(time.Duration(v.TTR) * time.Millisecond)
// 			if (v.TTR > 0) && (now.After(end)) {
// 				if !IsTerminalStatus(v.Status) {
// 					if err := v.Cancel(); err != nil {
// 						log.Debugf("failed cancel cluster %s %v StartAt %v", v.Id, err, v.StartAt)
// 					} else {
// 						log.Tracef("successfully canceled cluster %s StartAt %v, TTR %v msec", v.Id, v.StartAt, v.TTR)
// 					}
// 				}
// 				delete(r.all, k)
// 				num += 1
// 			}
//
// 		}
//
// 	}
// 	return num
// }

// GracefullShutdown is used when we stop the ClusterRegistry.
// cancel all running & pending cluster
// return false if we can't cancel any cluster
// func (r *ClusterRegistry) GracefullShutdown() bool {
// 	r.Cleanup()
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	failed := false
// 	// log.Debug("start GracefullShutdown")
// 	// for k, v := range r.all {
// 	// 	if !IsTerminalStatus(v.Status) {
// 	// 		if err := v.Cancel(); err != nil {
// 	// 			log.Debug(fmt.Sprintf("failed cancel cluster %s %v", v.Id, err))
// 	// 			failed = true
// 	// 		} else {
// 	// 			log.Debug(fmt.Sprintf("successfully canceled cluster %s", v.Id))
// 	// 		}
// 	// 	}
// 	// 	delete(r.all, k)
// 	// }
// 	return failed
// }

// func (r *ClusterRegistry) Filter(ClusterType string) ([]*Cluster, bool) {
//     res := make([]*Cluster,0)
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()
//     	for k, _ := range r.all {
//     		if i > n {
//     			slice = append(slice, k)
//     		} else {
//     			slice[i] = k
//
//     		}
//     		i++
//     	}
//
//
// 	if rec, ok := r.all[clusterId]; ok {
// 		return rec, true
// 	}
//
// 	return nil, false
// }

// Record fetch cluster by Cluster ID.
// Follows comma ok idiom
func (r *ClusterRegistry) Record(clusterId string) (*Cluster, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rec, ok := r.all[clusterId]; ok {
		return rec, true
	}

	return nil, false
}
