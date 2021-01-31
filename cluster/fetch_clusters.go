package cluster

import (
	"context"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	utils "github.com/weldpua2008/suprasched/utils"
	"strings"
	"time"
)

// GetSectionClustersFetcher returns ClustersFetcher from configuration file.
// By default http ClustersFetcher will be used.
// Example YAML config for `section` that will return new `RestClustersFetcher`:
//     section:
//         type: "HTTP"
func GetSectionClustersFetcher(section string) (ClustersFetcher, error) {
	fetcherType := config.GetStringDefault(fmt.Sprintf("%s.type", section), ConstructorsFetcherTypeRest)
	k := strings.ToUpper(fetcherType)
	if typeStruct, ok := FetcherConstructors[k]; ok {
		if comm, err := typeStruct.constructor(fmt.Sprintf("%s.%s", section, config.CFG_PREFIX_FETCHER)); err == nil {
			return comm, nil
		} else {
			return nil, err
		}

	}
	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableClustersFetcher, fetcherType)
}

// StartGenerateClusters goroutine for getting clusters from API with internal
// exists on kill
func StartGenerateClusters(ctx context.Context, clusters chan bool, interval time.Duration) error {
	singleFetcher, err := GetSectionClustersFetcher(config.CFG_PREFIX_CLUSTER)

	fetchers := make([]ClustersFetcher, 0)
	fetchers = append(fetchers, singleFetcher)
	if err != nil || fetchers == nil || len(fetchers) == 0 {
		close(clusters)
		return fmt.Errorf("Failed to start StartGenerateClusters %v", err)
	}

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting fetching Clusters with delay %v", interval)
	tickerGenerateClusters := time.NewTicker(interval)
	defer func() {
		tickerGenerateClusters.Stop()
		close(clusters)
	}()

	go func() {
		j := 0
		for {
			select {
			case <-ctx.Done():
				doneNumClusters <- j
				log.Debug("Clusters fetch finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				isDelayed := utils.RandomBoolean()

				for _, fetcher := range fetchers {
					if isDelayed {
						break
					}
					clustersSlice, err := fetcher.Fetch()
					if err == nil {
						for _, cls := range clustersSlice {
							var topic string
							rec, exist := config.ClusterRegistry.Record(cls.StoreKey())
							// log.Tracef("rec  %v %v %v", cls.ClusterId, cls.ClusterType, cls.Status)
							if !exist {
								if !config.ClusterRegistry.Add(cls) {
									log.Tracef("!config.ClusterRegistry.Add %v", cls)
									continue
								} else {

									// config.ClusterRegistry.MarkFree(cls.StoreKey())
									// log.Tracef("Cluster %v TimeOutAt %v", cls.ClusterId, cls.TimeOutAt)
									topic = config.TOPIC_CLUSTER_CREATED
									rec, exist = config.ClusterRegistry.Record(cls.StoreKey())
									if !exist {
										log.Tracef("Skip !config.ClusterRegistry.Record %v", cls)
										continue
									}
									// if rec.IsFree() {
									//
									// }

									clustersFetched.Inc()
									j += 1
								}
							} else if rec.IsInTransition() {
								// log.Tracef("Skip IsInTransition %v", rec)
								continue
							} else if rec.UpdateStatus(cls.Status) {
								topic = strings.ToLower(fmt.Sprintf("cluster.%v", rec.Status))
							}
							if (rec != nil) && (rec.ClusterType != model.CLUSTER_TYPE_ON_DEMAND) && (len(rec.StoreKey()) > 0) {
								clusterEventMetadata := map[string]string{"StoreKey": rec.StoreKey()}
								_, err := config.Bus.Emit(ctx, config.TOPIC_CLUSTER_IS_EMPTY, clusterEventMetadata)
								if err != nil {
									log.Tracef("%v", err)
								}
							}

							if len(topic) > 1 {
								_, err := config.Bus.Emit(ctx, topic, rec.EventMetadata())
								if err != nil {
									log.Tracef("%v", err)
								}
							}
						}

					} else {
						log.Tracef("Fetch cluster metadata '%v', but failed with %v", clustersSlice, err)

					}
				}
			}
		}
	}()

	numSentClusters := <-doneNumClusters

	log.Infof("Fetched %v clusters", numSentClusters)
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
