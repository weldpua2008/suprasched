package cache

import (
	"container/list"
	"fmt"
	"github.com/weldpua2008/suprasched/core"
	"sync"
	"time"
)

var ErrClusterAlreadyExists = fmt.Errorf("cluster already exists in the cache")

// New returns a Cache implementation.
func New(ttl time.Duration) Cache {
	return newSchedulerCache(ttl, 1*time.Second)
}

func newSchedulerCache(ttl, period time.Duration) *schedulerCache {
	return &schedulerCache{
		ttl:        ttl,
		period:     period,
		clusters:   make(map[core.Namespace]*list.List),
		jobs:       make(map[core.Namespace]*list.List),
		namespaces: make([]core.Namespace, 0),
	}
}

type schedulerCache struct {
	ttl    time.Duration
	period time.Duration

	// This mutex guards all fields within schedulerCache struct.
	mu sync.RWMutex

	clusters   map[core.Namespace]*list.List // a map from namespace to an array of clusters.
	jobs       map[core.Namespace]*list.List // a map from namespace to an array of jobs.
	namespaces []core.Namespace
}

func (cache *schedulerCache) ClusterCount() int {
	panic("implement me")
}

func (cache *schedulerCache) JobCount() (int, error) {
	panic("implement me")
}

func (cache *schedulerCache) AddJob(job *core.Job) error {
	panic("implement me")
}

func (cache *schedulerCache) UpdateJob(oldJob, newJob *core.Job) error {
	panic("implement me")
}

func (cache *schedulerCache) RemoveJob(job *core.Job) error {
	panic("implement me")
}

func (cache *schedulerCache) GetJob(job *core.Job) (*core.Job, error) {
	panic("implement me")
}

func (cache *schedulerCache) IsAssumedJob(job *core.Job) (bool, error) {
	panic("implement me")
}

func (cache *schedulerCache) AddCluster(cluster *core.Cluster) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if l, ok := cache.clusters[cluster.Namespace]; ok {
		for e := l.Front(); e != nil; e = e.Next() {
			cl := e.Value.(core.Cluster)
			if cl.Name == cluster.Name {
				return ErrClusterAlreadyExists
			}
		}
		l.PushBack(cluster)
	} else {
		var l *list.List
		l = list.New()
		l.PushBack(cluster)
		cache.clusters[cluster.Namespace] = l
	}

	return nil
}

func (cache *schedulerCache) UpdateCluster(oldCluster, newCluster *core.Cluster) error {
	panic("implement me")
}

func (cache *schedulerCache) RemoveCluster(cluster *core.Cluster) error {
	panic("implement me")
}

func (cache *schedulerCache) Dump() *Dump {
	panic("implement me")
}

func (cache *schedulerCache) UpdateSnapshot(currSnapshot *Snapshot) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for ns, clustersList := range cache.clusters {
		snapshotClusters, ok := currSnapshot.clusters[ns]
		if !ok {
			var t *list.List
			currSnapshot.clusters[ns] = t
		}
		for e := clustersList.Front(); e != nil; e = e.Next() {
			cl := e.Value.(core.Cluster)
			snapshotClusters.PushBack(cl)
		}

	}
	return nil
}

// Dump is a dump of the cache state.
type Dump struct {
	Jobs     *map[core.Namespace]list.List
	Clusters *map[core.Namespace]list.List
}
