package cache

import (
	"container/list"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weldpua2008/suprasched/core"
	"sync"
	"time"
)

var ErrClusterAlreadyExists = fmt.Errorf("cluster already exists in the cache")
var ErrClusterNotExists = fmt.Errorf("cluster not exists in the cache")
var ErrEmptyCache = fmt.Errorf("no clusters in the cache")

// New returns a Cache implementation.
func New(ttl time.Duration) Cache {
	return newSchedulerCache(ttl, 1*time.Second)
}

func newSchedulerCache(ttl, period time.Duration) *schedulerCache {
	return &schedulerCache{
		ttl:        ttl,
		period:     period,
		clusters:   make(map[core.Namespace]map[core.UID]*core.Cluster),
		jobs:       make(map[core.Namespace]*list.List),
		namespaces: make([]core.Namespace, 0),
	}
}

type schedulerCache struct {
	ttl    time.Duration
	period time.Duration

	// This mutex guards all fields within schedulerCache struct.
	mu sync.RWMutex

	clusters   map[core.Namespace]map[core.UID]*core.Cluster // a map from namespace to a map of clusters.
	jobs       map[core.Namespace]*list.List                 // a map from namespace to an array of jobs.
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
	logrus.Tracef("Adding Cluster %v to ns %v", cluster.Name, cluster.Namespace)
	_, ok := cache.clusters[cluster.Namespace]
	if !ok {
		cache.clusters[cluster.Namespace] = map[core.UID]*core.Cluster{cluster.UID: cluster}
	}
	cache.clusters[cluster.Namespace][cluster.UID] = cluster
	/* else {
		if _, ok1:= l[cluster.UID]; !ok1 {
			l[cluster.UID] = cluster
		}
	} */
	//if ok {
	//for e := l.Front(); e != nil; e = e.Next() {
	//	val := e.Value
	//	if cl, ok := val.(core.Cluster); ok {
	//		if cl.UID == cluster.UID {
	//			return ErrClusterAlreadyExists
	//		}
	//	} else if cl, ok := val.(*core.Cluster); ok {
	//		if cl.UID == cluster.UID {
	//			return ErrClusterAlreadyExists
	//		}
	//	}
	//}
	//l.PushBack(cluster)
	//} else {
	//l := list.New()
	//l.PushBack(cluster)
	//cache.clusters[cluster.Namespace] = l
	//}

	return nil
}

func (cache *schedulerCache) UpdateCluster(oldCluster, newCluster *core.Cluster) error {
	panic("implement me")
}

func (cache *schedulerCache) RemoveCluster(cluster *core.Cluster) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if _, ok := cache.clusters[cluster.Namespace]; ok {
		if _, ok1 := cache.clusters[cluster.Namespace][cluster.UID]; ok1 {
			delete(cache.clusters[cluster.Namespace], cluster.UID)
			return nil
		}
	}

	//if l, ok := cache.clusters[cluster.Namespace]; ok {
	//	for e := l.Front(); e != nil; e = e.Next() {
	//		val := e.Value
	//		if cl, ok := val.(core.Cluster); ok {
	//			if cl.UID == cluster.UID {
	//				l.Remove(e)
	//				return nil
	//			}
	//		} else if cl, ok := val.(*core.Cluster); ok {
	//			if cl.UID == cluster.UID {
	//				l.Remove(e)
	//				return nil
	//			}
	//		}
	//	}
	//}
	return ErrClusterNotExists
}

func (cache *schedulerCache) Dump() *Dump {
	panic("implement me")
}

func (cache *schedulerCache) UpdateSnapshot(currSnapshot *Snapshot) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.clusters == nil {
		return ErrEmptyCache
	}
	//if currSnapshot.clusters == nil {
	currSnapshot.clusters = make(map[core.Namespace]map[core.UID]*core.Cluster, 0)
	//}
	for _, clustersList := range cache.clusters {
		//if len(clustersList) == 0 {
		//	continue
		//}
		for _, cl := range clustersList {
			logrus.Tracef("==> ns %v adding %v", cl.Namespace, cl.Name)

			if _, ok := currSnapshot.clusters[cl.Namespace]; !ok {
				currSnapshot.clusters[cl.Namespace] = make(map[core.UID]*core.Cluster, 1)
			}
			// TODO: Add check revision
			//if _, ok1 := currSnapshot.clusters[ns][cl.UID]; !ok1 {
			currSnapshot.clusters[cl.Namespace][cl.UID] = cl
			logrus.Tracef("Snapshot ns %v adding %v", cl.Namespace, cl.Name)
			//}

			//
			//if clustersList.Len() == 0 {
			//	continue
			//}
			//for e := clustersList.Front(); e != nil; e = e.Next() {
			//	val := e.Value
			//	var tempCluster *core.Cluster
			//	if cl, ok := val.(core.Cluster); ok {
			//		//logrus.Tracef("Snapshot ns %v adding %v in %v",ns, cl.Name, cl.Namespace )
			//		tempCluster = &cl
			//		//snapshotClusters.PushBack(&cl)
			//	} else if cl, ok := val.(*core.Cluster); ok {
			//		tempCluster = cl
			//		//snapshotClusters.PushBack(cl)
			//		//logrus.Tracef("Snapshot ns %v adding %v in %v",ns, cl.Name, cl.Namespace )
			//	}
			//	if tempCluster == nil {
			//		continue
			//	}
			//	_, ok := currSnapshot.clusters[tempCluster.Namespace]
			//	if !ok {
			//		currSnapshot.clusters[tempCluster.Namespace] = list.New()
			//	}
			//	currSnapshot.clusters[tempCluster.Namespace].PushBack(tempCluster)
			//	logrus.Tracef("Snapshot ns %v adding %v",tempCluster.Namespace, tempCluster.Name )

		}

	}
	return nil
}

// Dump is a dump of the cache state.
type Dump struct {
	Jobs     *map[core.Namespace]list.List
	Clusters *map[core.Namespace]list.List
}
