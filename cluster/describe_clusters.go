package cluster

import (
	"context"
	// "math/rand"
	"errors"
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
func GetSectionClustersDescriber(section string) ([]ClustersDescriber, error) {

	def := make(map[string]string)

	describers_cfgs := config.GetMapStringMapStringTemplatedDefault(section, config.CFG_PREFIX_DESCRIBERS, def)
	res := make([]ClustersDescriber, 0)
	for subsection, comm := range describers_cfgs {
		if comm == nil {
			continue
		}
		describer_type := ConstructorsDescriberTypeRest
		if descr_type, ok := comm["type"]; ok {
			describer_type = descr_type
		}
		k := strings.ToUpper(describer_type)
		if type_struct, ok := DescriberConstructors[k]; ok {
			describer_instance, err := type_struct.constructor(fmt.Sprintf("%v", subsection))
			if err != nil {
				log.Tracef("Can't get describer %v", err)
				continue
			}
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
	notValidClusterIds := make(map[string]struct{}, 0)

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting update Clusters with delay %v", interval)
	tickerGenerateClusters := time.NewTicker(interval)
	defer func() {
		tickerGenerateClusters.Stop()
	}()

	go func() {
		cntr := 0
		for {
			select {
			case <-ctx.Done():
				close(clusters)
				doneNumClusters <- cntr
				log.Debug("Clusters description finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				isDelayed := utils.RandomBoolean()
				for _, describer := range describers_instances {
					if isDelayed {
						break
					}

					for _, cls := range describer.SupportedClusters() {

						rec, ok := config.ClusterRegistry.Record(cls.StoreKey())
						if !ok {
							continue
						}
						if !time.Now().After(rec.LastSyncedAt) {
							continue
						}

						if _, ok := notValidClusterIds[rec.ClusterId]; ok {
							// log.Tracef("Skip %v", rec.ClusterId)
							continue
						}

						params := rec.GetParams()
						cluster_status, err := describer.ClusterStatus(params)
						if err == nil {
							reqClustersDescribed.Inc()
							var topic string
							if rec.IsInTransition() {
								continue
							}

							if rec.UpdateStatus(cluster_status) {
								// log.Tracef("=> %v %v", rec.ClusterId, cluster_status)
								if model.IsTerminalStatus(cluster_status) {
									rec.PutInTransition()
								}
								clustersDescribed.Inc()
								cntr += 1
								topic = strings.ToLower(fmt.Sprintf("cluster.%v", cluster_status))
								_, err := config.Bus.Emit(ctx, topic, rec.EventMetadata())
								if err != nil {
									log.Tracef("%v", err)
								}
								if rec.Status != cluster_status {
									log.Tracef("rec.Status %v != %v", rec.Status, cluster_status)

								}
							}
						} else if errors.Is(err, ErrClusterIdIsNotValid) {
							/*
							   TODO: It's better to remove such cluster and log once
							*/
							clusterIdsAreNotValid.Set(float64(len(notValidClusterIds)))
							if len(notValidClusterIds) > 4096 {
								notValidClusterIds = make(map[string]struct{}, 0)
							}
							notValidClusterIds[rec.ClusterId] = struct{}{}
							if config.ClusterRegistry.Delete(cls.StoreKey()) {
								cls = nil
							}
							continue

						} else {
							reqClustersFailDescribed.Inc()
							log.Tracef("Failed to describe cluster status '%v', failed with %v", cluster_status, err)
						}
					}
					config.ClusterRegistry.DumpMetrics(clusterStatuses)
					// clusterStatuses
				}
			}
		}
	}()

	numSentClusters := <-doneNumClusters

	log.Infof("Described %v clusters", numSentClusters)
	return nil
}
