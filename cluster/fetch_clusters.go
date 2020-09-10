package cluster

import (
	"context"
	"fmt"
	// "github.com/sirupsen/logrus"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	// "strconv"
	"strings"
	"time"
)

// StartGenerateClusters goroutine for getting clusters from API with internal
// exists on kill
func StartGenerateClusters(ctx context.Context, clusters chan *model.Cluster, interval time.Duration) error {
	fetcher, err := NewFetchClustersHttp()
	if err != nil {
		close(clusters)
		return fmt.Errorf("Failed to start StartGenerateClusters %v", err)
	}

	doneNumClusters := make(chan int, 1)
	log.Info(fmt.Sprintf("Starting fetching Clusters with delay %v", interval))
	tickerGenerateClusters := time.NewTicker(interval)
	defer func() {
		tickerGenerateClusters.Stop()
	}()

	go func() {
		j := 0
		for {
			select {
			case <-ctx.Done():
				close(clusters)
				doneNumClusters <- j
				log.Debug("Clusters generation finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				clusters_slice, err := fetcher.Fetch()
				if err == nil {

					for _, cls := range clusters_slice {
						var topic string
						if !config.ClusterRegistry.Add(cls) {
							if rec, exist := config.ClusterRegistry.Record(cls.StoreKey()); exist {
								if rec.UseExternaleStatus(cls) {
									topic = strings.ToLower(fmt.Sprintf("cluster.%v", cls.Status))
								}

							}
						} else {
							topic = config.TOPIC_CLUSTER_CREATED
						}
						if len(topic) > 0 {
							_, err := config.Bus.Emit(ctx, topic, cls.EventMetadata())
							if err != nil {
								log.Tracef("%v", err)
							}

						}

					}

				} else {
					log.Tracef("%v, %v", clusters_slice, err)

				}

			}
		}
	}()

	numSentClusters := <-doneNumClusters

	log.Info(fmt.Sprintf("Sent %v clusters", numSentClusters))
	return nil
}

// // GracefullShutdown cancel all running clusters
// // returns error in case any job failed to cancel
// func GracefullShutdown(clusters <-chan *model.Job) bool {
// 	// empty clusters channel
// 	if len(clusters) > 0 {
// 		log.Trace(fmt.Sprintf("clusters chan still has size %v, empty it", len(clusters)))
// 		for len(clusters) > 0 {
// 			<-clusters
// 		}
// 	}
// 	ClustersRegistry.GracefullShutdown()
// 	if ClustersRegistry.Len() > 0 {
// 		log.Trace(fmt.Sprintf("GracefullShutdown failed, '%v' clusters left ", ClustersRegistry.Len()))
// 		return false
// 	}
// 	return true
//
// }
