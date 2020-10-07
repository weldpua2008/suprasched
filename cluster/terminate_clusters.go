package cluster

import (
	"context"
	// "math/rand"
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
func GetSectionClustersTerminator(section string) ([]ClustersTerminator, error) {

	def := make(map[string]string)

	terminators_cfgs := config.GetMapStringMapStringTemplatedDefault(section, config.CFG_PREFIX_TERMINATORS, def)
	res := make([]ClustersTerminator, 0)
	for subsection, comm := range terminators_cfgs {
		if comm == nil {
			continue
		}
		terminator_type := ConstructorsTerminaterTypeEMR
		if term_type, ok := comm["type"]; ok {
			terminator_type = term_type
		}
		k := strings.ToUpper(terminator_type)
		if type_struct, ok := TerminatorConstructors[k]; ok {
			terminator_instance, err := type_struct.constructor(fmt.Sprintf("%v", subsection))
			if err != nil {
				log.Tracef("Can't get terminator %v", err)
				continue
			}
			res = append(res, terminator_instance)

		}

	}
	if len(res) > 0 {
		return res, nil
	}

	return nil, fmt.Errorf("%w for %s.\n", ErrNoSuitableClustersTerminator, section)
}

// StartTerminateClusters goroutine for terminatting clusters & updating API with internal
func StartTerminateClusters(ctx context.Context, clusters chan *model.Cluster, interval time.Duration) error {
	terminators_instances, err := GetSectionClustersTerminator(config.CFG_PREFIX_CLUSTER)

	if err != nil || terminators_instances == nil || len(terminators_instances) == 0 {
		close(clusters)
		return fmt.Errorf("Failed to start StartTerminateClusters %v", err)
	}
	notValidClusterIds := make(map[string]struct{}, 0)

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting terminate Clusters with delay %v", interval)
	tickerTerminateClusters := time.NewTicker(interval)
	defer func() {
		tickerTerminateClusters.Stop()
	}()

	go func() {
		cntr := 0
		for {
			select {
			case <-ctx.Done():
				close(clusters)
				doneNumClusters <- cntr
				log.Debug("Clusters termination finished [ SUCCESSFULLY ]")
				return
			case <-tickerTerminateClusters.C:
				isDelayed := utils.RandomBoolean()
				for _, terminator := range terminators_instances {
					if isDelayed {
						break
					}

					for _, cls := range terminator.SupportedClusters() {

						rec, ok := config.ClusterRegistry.Record(cls.StoreKey())
						if !ok {
							continue
						}

						if !rec.IsTimeout() {
							continue
						}

						if _, ok := notValidClusterIds[rec.ClusterId]; ok {
							continue
						}

						params := rec.GetParams()
						err := terminator.Terminate(params)
						if err == nil {
							cntr += 1
							metrics.ReqClustersTerminated.WithLabelValues(
								"aws",
								strings.ToLower(fmt.Sprintf("%v.%v", rec.ClusterProfile, rec.ClusterRegion)),
								"emr",
							).Inc()
							if rec.IsInTransition() {
								continue
							}
							rec.SyncedWithExternalAPI()
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
							log.Tracef("Failed to terminate cluster status '%v', failed with %v", rec.ClusterId, err)
						}
					}
					config.ClusterRegistry.DumpMetrics(clusterStatuses)
					// clusterStatuses
				}
			}
		}
	}()

	numTermClusters := <-doneNumClusters

	log.Infof("Terminated %v clusters", numTermClusters)
	return nil
}
