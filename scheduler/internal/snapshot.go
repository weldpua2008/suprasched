package cache

import (
	"container/list"
	"github.com/weldpua2008/suprasched/core"
)

// Snapshot is a snapshot of cache NodeInfo and NodeTree order. The scheduler takes a
// snapshot at the beginning of each scheduling cycle and uses it for its operations in that cycle.
type Snapshot struct {
	clusters   map[core.Namespace]map[core.UID]*core.Cluster // a map from namespace to an array of clusters.
	jobs       map[core.Namespace]*list.List                 // a map from namespace to an array of jobs.
	generation int64
}

// NumNodes returns the number of nodes in the snapshot.
func (s *Snapshot) NumClusters(ns core.Namespace) int {
	if val, ok := s.clusters[ns]; ok && val != nil {
		return len(val)
		//return val.Len()
	}
	return 0
}

// GetClustersFromNs returns the clusters in the snapshot.
func (s *Snapshot) GetClustersFromNs(ns core.Namespace) (ret []core.Cluster) {
	if s.clusters == nil {
		return ret
	}

	if val, ok := s.clusters[ns]; ok && val != nil {
		for _, cl := range val {
			ret = append(ret, *cl)
		}
		//	for e := val.Front(); e != nil; e = e.Next() {
		//		val := e.Value
		//		if cl, ok := val.(core.Cluster); ok {
		//			ret = append(ret, cl)
		//		} else if cl, ok := val.(*core.Cluster); ok {
		//			ret = append(ret, *cl)
		//		}
		//	}
	}
	return ret
}
