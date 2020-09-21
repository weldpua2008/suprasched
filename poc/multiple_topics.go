package poc

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	"sync"
	"time"
)

var (
	log        = logrus.WithFields(logrus.Fields{"package": "PoC"})
	job_cntr   int
	empty_cntr int

	assigned_jobs    int
	new_cluster_cntr int
	mu               sync.RWMutex
)

// config.InitEvenBus()

func InitializeClusters(num int, types []string, statuses []string) error {
	for i := 0; i < (num); i++ {
		for t := 0; t < len(types); t++ {

			for z := 0; z < len(statuses); z++ {
				if config.ClusterRegistry.Len() > num {
					break
				}
				cl := model.NewCluster(fmt.Sprintf("Cluster-%v.%v.%v", i, t, z))
				cl.Status = statuses[z]
				cl.ClusterType = types[t]
				cl.MaxCapacity = 5
				if !config.ClusterRegistry.Add(cl) {
					return fmt.Errorf("Can't add cluster N%v cluster %v", i, cl)
				}
			}
		}
	}
	return nil
}

func InitializeJobs(num int, cluster_types []string, statuses []string) error {

	for i := 0; i < (num); i++ {
		for t := 0; t < len(cluster_types); t++ {

			for z := 0; z < len(statuses); z++ {
				if config.JobsRegistry.Len() > num {
					break
				}
				j := model.NewJob(fmt.Sprintf("Job-%v.%v.%v", i, t, z))
				j.Status = statuses[z]
				j.ClusterType = cluster_types[t]
				if !config.JobsRegistry.Add(j) {
					return fmt.Errorf("Can't add cluster N%v", i)
				}
			}
		}
	}
	return nil

}

func EmptyHandler(e *bus.Event) {
	mu.Lock()
	empty_cntr += 1
	mu.Unlock()
}

func ParallelHandler(e *bus.Event) {
	mu.Lock()
	empty_cntr += 1
	mu.Unlock()
	time.Sleep(500 * time.Millisecond)
}

func AssignNewFreeCluster(e *bus.Event) {

	clusterType := e.Data.(map[string]string)["ClusterType"]
	storeKey := e.Data.(map[string]string)["StoreKey"]
	if rec, ok := config.JobsRegistry.Record(storeKey); ok {
		if rec == nil || rec.StoreKey() == "" {
			log.Warningf("Empty Job Event in %s: %+v", e.Topic, e)

		} else if len(rec.GetClusterStoreKey()) == 0 {
			assigned := false
			cluster_types := []string{clusterType}
			for _, cl := range config.ClusterRegistry.Filter(cluster_types) {
				if cl.IsFull() {
					continue
				}
				if !cl.Add(rec) {
					log.Warningf("Can't add Job to cluster %v Event in %s: %+v", cl, e.Topic, e)
					cl.Delete(rec.StoreKey())
					continue
				}
				rec.ChangeClusterStoreKey(cl.StoreKey())

				assigned = true
				break
			}
			if !assigned {
				cl := model.NewCluster(fmt.Sprintf("%v.%v", e.ID, config.ClusterRegistry.Len()))
				cl.Status = "STARTING"
				cl.ClusterType = clusterType
				if !config.ClusterRegistry.Add(cl) {
					panic(fmt.Sprintf("Can't create cluster %v", cl))
				}
				rec.ChangeClusterStoreKey(cl.StoreKey())
				assigned = true

			}
			if assigned {
				mu.Lock()
				defer mu.Unlock()
				assigned_jobs += 1
			}
		}

	}

	// log.Infof("Event for %s: %+v", e.Topic, e)
}

func FreeJobFromCluster(e *bus.Event) {
	// log.Infof("Event for %s: %+v", e.Topic, e)
	// return
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are getting the kill signal or exit

	storeKey := e.Data.(map[string]string)["StoreKey"]
	if rec, ok := config.JobsRegistry.Record(storeKey); ok {
		if rec == nil || storeKey == "" {
			log.Warningf("Empty Job Event in %s: %+v", e.Topic, e)

		}
		clusterStoreKey := rec.GetClusterStoreKey()

		// ev:=rec.EventMetadata()
		if !config.JobsRegistry.Delete(storeKey) {
			log.Warningf("Can't delete Job from JobsRegistry Event in %s: %+v", e.Topic, e)
		}

		if cl, ok := config.ClusterRegistry.Record(clusterStoreKey); ok {
			// if !cl.Delete(rec.StoreKey()) {
			// 	log.Warningf("Can't delete Job from cl %v Event in %s: %+v", cl, e.Topic, e)
			// }
			// log.Infof("Remove Job from cl %v Event in %s: %+v", cl, e.Topic, e)

			cl.Delete(storeKey)
			topic := config.TOPIC_CLUSTER_TERMINATED
			_, err := config.Bus.Emit(ctx, topic, map[string]string{"StoreKey": clusterStoreKey})
			if err != nil {
				log.Tracef("%v", err)
			}

			// if cl.Len() < 1 {
			// 	topic := config.TOPIC_CLUSTER_TERMINATED
			// 	_, err := config.Bus.Emit(ctx, topic, map[string]string{"StoreKey":clusterStoreKey})
			// 	if err != nil {
			// 		log.Tracef("%v", err)
			// 	}
			// }

			// topic := config.TOPIC_CLUSTER_TERMINATED
			// _, err := config.Bus.Emit(nil, topic, ev)
			// if err != nil {
			// 	log.Tracef("%v", err)
			// }
		}
		// topic := config.TOPIC_CLUSTER_TERMINATED
		// _, err := config.Bus.Emit(ctx, topic, map[string]string{"StoreKey": clusterStoreKey})
		// if err != nil {
		// 	log.Tracef("%v", err)
		// }

	}
	// else {
	// 	log.Warningf("Can't find Job %v", storeKey)
	// }

	// log.Infof("Event for %s: %+v", e.Topic, e)
}

func DeleteCluster(e *bus.Event) {
	storeKey := e.Data.(map[string]string)["StoreKey"]
	// log.Warningf("config.ClusterRegistry %v", config.ClusterRegistry.Len())
	if cl, ok := config.ClusterRegistry.Record(storeKey); ok {
		// log.Warningf("c %v Event in %s: %+v", cl, e.Topic, e)

		if cl.Len() < 1 {
			if !config.ClusterRegistry.Delete(storeKey) {
				log.Warningf("Can't delete cl %v Event in %s: %+v", cl, e.Topic, e)

			}
		} else {
			log.Warningf("Len %v", cl.Len())
		}
	} else {
		log.Infof("Can't find %v", storeKey)
	}

}
