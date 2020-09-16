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

// GetSectionClustersFetcher returns ClustersFetcher from configuration file.
// By default http ClustersFetcher will be used.
// Example YAML config for `section` that will return new `RestClustersFetcher`:
//     section:
//         type: "HTTP"
func GetSectionClustersFetcher(section string) (ClustersFetcher, error) {
	fetcher_type := config.GetStringDefault(fmt.Sprintf("%s.type", section), ConstructorsFetcherTypeRest)
	k := strings.ToUpper(fetcher_type)
	if type_struct, ok := FetcherConstructors[k]; ok {
		if comm, err := type_struct.constructor(fmt.Sprintf("%s.%s", section, config.CFG_PREFIX_FETCHER)); err == nil {
			// log.Infof("ClustersFetcher %v", comm)
			return comm, nil
		} else {
			return nil, err
		}

	}
	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableClustersFetcher, fetcher_type)
}

// StartGenerateClusters goroutine for getting clusters from API with internal
// exists on kill
func StartGenerateClusters(ctx context.Context, clusters chan *model.Cluster, interval time.Duration) error {
	single_fetcher, err := GetSectionClustersFetcher(config.CFG_PREFIX_CLUSTER)

	fetchers := make([]ClustersFetcher, 0)
	fetchers = append(fetchers, single_fetcher)
	// comms, err:=communicator.GetCommunicatorsFromSection(fmt.Sprintf("%s.fetch", config.CFG_PREFIX_CLUSTER))
	if err != nil || fetchers == nil || len(fetchers) == 0 {
		close(clusters)
		return fmt.Errorf("Failed to start StartGenerateClusters %v", err)
	}

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting fetching Clusters with delay %v", interval)
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
				log.Debug("Clusters fetch finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				for _, fetcher := range fetchers {

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
						log.Tracef("Fetch cluster metadata '%v', but failed with %v", clusters_slice, err)

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
