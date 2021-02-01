package cluster

import (
	"context"
	"errors"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	metrics "github.com/weldpua2008/suprasched/metrics"
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

	describersCFGs := config.GetMapStringMapStringTemplatedDefault(section, config.CFG_PREFIX_DESCRIBERS, def)
	res := make([]ClustersDescriber, 0)
	for subsection, comm := range describersCFGs {
		if comm == nil {
			continue
		}
		describerType := ConstructorsDescriberTypeRest
		if tmpDescribeType, ok := comm["type"]; ok {
			describerType = tmpDescribeType
		}
		k := strings.ToUpper(describerType)
		if typeStruct, ok := DescriberConstructors[k]; ok {
			describerInstance, err := typeStruct.constructor(fmt.Sprintf("%v", subsection))
			if err != nil {
				log.Tracef("Can't get describer %v", err)
				continue
			}
			res = append(res, describerInstance)

		}

	}
	if len(res) > 0 {
		return res, nil
	}

	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableClustersDescriber, section)
}

// StartUpdateClustersMetadata goroutine for getting clusters from API with internal
// exists on kill
func StartUpdateClustersMetadata(ctx context.Context, clusters chan bool, interval time.Duration) error {
	describersInstances, err := GetSectionClustersDescriber(config.CFG_PREFIX_CLUSTER)

	if err != nil || describersInstances == nil || len(describersInstances) == 0 {
		return fmt.Errorf("Failed to start StartUpdateClustersMetadata %v", err)
	}
	notValidClusterIds := make(map[string]struct{})

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting update Clusters with delay %v", interval)
	tickerGenerateClusters := time.NewTicker(interval)
	defer func() {
		close(clusters)
		tickerGenerateClusters.Stop()
	}()

	go func() {
		cntr := 0
		for {
			select {
			case <-ctx.Done():
				doneNumClusters <- cntr
				log.Debug("Clusters description finished [ SUCCESSFULLY ]")
				return
			case <-tickerGenerateClusters.C:
				start := time.Now()
				isDelayed := utils.RandomBoolean()
				for _, describer := range describersInstances {
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
							continue
						}

						params := rec.GetParams()
						clusterStatus, err := describer.ClusterStatus(params)
						if err == nil {
							reqClustersDescribed.Inc()
							var topic string
							if rec.IsInTransition() {
								continue
							}
							rec.SyncedWithExternalAPI()
							if rec.UpdateStatus(clusterStatus) {
								// log.Tracef("=> %v %v", rec.ClusterId, cluster_status)
								if model.IsTerminalStatus(clusterStatus) {
									rec.PutInTransition()
								}
								clustersDescribed.Inc()
								cntr += 1
								topic = strings.ToLower(fmt.Sprintf("cluster.%v", clusterStatus))
								_, err := config.Bus.Emit(ctx, topic, rec.EventMetadata())
								if err != nil {
									log.Tracef("%v", err)
								}
								if rec.Status != clusterStatus {
									log.Tracef("rec.Status %v != %v", rec.Status, clusterStatus)

								}
							}
						} else if errors.Is(err, ErrClusterIdIsNotValid) {
							/*
							   TODO: It's better to remove such cluster and log once
							*/
							clusterIdsAreNotValid.Set(float64(len(notValidClusterIds)))
							if len(notValidClusterIds) > 4096 {
								notValidClusterIds = make(map[string]struct{})
							}
							notValidClusterIds[rec.ClusterId] = struct{}{}
							if config.ClusterRegistry.Delete(cls.StoreKey()) {
								cls = nil
							}
							continue

						} else {
							reqClustersFailDescribed.Inc()
							log.Tracef("Failed to describe cluster status '%v', failed with %v", clusterStatus, err)
						}
						metrics.FetchMetadataLatency.WithLabelValues("describe_clusters",
							"single").Observe(float64(time.Now().Sub(start).Nanoseconds()))

					}
				}

				if !isDelayed {
					config.ClusterRegistry.DumpMetrics(clusterStatuses)
					metrics.FetchMetadataLatency.WithLabelValues("describe_clusters",
						"whole").Observe(float64(time.Now().Sub(start).Nanoseconds()))

				}

			}
		}
	}()

	numSentClusters := <-doneNumClusters

	log.Infof("Described %v clusters", numSentClusters)
	return nil
}
