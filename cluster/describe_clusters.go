package cluster

import (
	"context"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	"strings"
	"time"
)

// GetSectionClustersFetcher returns ClustersFetcher from configuration file.
// By default http ClustersFetcher will be used.
// Example YAML config for `section` that will return new `RestClustersFetcher`:
//     section:
//         type: "HTTP"
func GetSectionClustersDescriber(section string) ([]ClustersDescriber, error) {

	def := make(map[string]string)
	describers_cfgs := config.GetSliceStringMapStringTemplatedDefault(section, config.CFG_PREFIX_DESCRIBERS, def)
	res := make([]ClustersDescriber, 0)
	for _, comm := range describers_cfgs {
		if comm == nil {
			continue
		}
		describer_type := ConstructorsDescriberTypeRest
		if descr_type, ok := comm["type"]; ok {
			describer_type = descr_type
		}
		k := strings.ToUpper(describer_type)
		if type_struct, ok := DescriberConstructors[k]; ok {
			describer_instance, err := type_struct.constructor(section)
			if err != nil {
				log.Tracef("Can't get describer %v", err)
				continue
			}
			// if err1 := describer_instance.Configure(config.ConvertMapStringToInterface(comm)); err1 != nil {
			//     log.Tracef("Can't configure %v communicator, got %v", describer_type, comm)
			//     return nil, err1
			// }
			res = append(res, describer_instance)

		}

	}
	if len(res) > 0 {
		return res, nil
	}

	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableClustersDescriber, section)
}

// StartUpdateClustersMetadata goroutine for getting clusters from API with internal
// exists on kill
func StartUpdateClustersMetadata(ctx context.Context, clusters chan *model.Cluster, interval time.Duration) error {
	describers_instances, err := GetSectionClustersDescriber(config.CFG_PREFIX_CLUSTER)
	if err != nil || describers_instances == nil || len(describers_instances) == 0 {
		close(clusters)
		return fmt.Errorf("Failed to start StartUpdateClustersMetadata %v", err)
	}

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting update Clusters with delay %v", interval)
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
				log.Debug("Clusters description finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				for _, describer := range describers_instances {

					supported_cluster := describer.SupportedClusters()
					for _, claster := range supported_cluster {

						params := claster.GetParams()

						cluster_status, err := describer.ClusterStatus(params)
						if err == nil {

							var topic string
							rec, ok := config.ClusterRegistry.Record(claster.StoreKey())
							if !ok {
								continue
							}

							if rec.UpdateStatus(cluster_status) {
								topic = strings.ToLower(fmt.Sprintf("cluster.%v", cluster_status))
								_, err := config.Bus.Emit(ctx, topic, rec.EventMetadata())
								if err != nil {
									log.Tracef("%v", err)
								}

							}
						} else {
							log.Tracef("Failed to describe cluster status '%v', but failed with %v", cluster_status, err)

						}
					}
				}
			}
		}
	}()

	numSentClusters := <-doneNumClusters

	log.Infof("Described %v clusters", numSentClusters)
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
