package cluster

import (
	"context"
	"errors"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	metrics "github.com/weldpua2008/suprasched/metrics"
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
func StartTerminateClusters(ctx context.Context, clusters chan bool, interval time.Duration, delay time.Duration) error {
	terminators_instances, err := GetSectionClustersTerminator(config.CFG_PREFIX_CLUSTER)

	if err != nil || terminators_instances == nil || len(terminators_instances) == 0 {
		return fmt.Errorf("Failed to start StartTerminateClusters %v", err)
	}
	notValidClusterIds := make(map[string]struct{}, 0)

	doneNumClusters := make(chan int, 1)
	log.Infof("Starting terminate Clusters every %v after %v", interval, delay)
	tickerTerminateClusters := time.NewTicker(interval)
	defer func() {
		tickerTerminateClusters.Stop()
		close(clusters)
	}()

    time.Sleep(delay)
	go func() {
		cntr := 0
		for {
			select {
			case <-ctx.Done():
				doneNumClusters <- cntr
				log.Debug("Clusters termination finished [ SUCCESSFULLY ]")
				return
			case <-tickerTerminateClusters.C:
				isDelayed := utils.RandomBoolean()
				start := time.Now()
				for _, terminator := range terminators_instances {
					if isDelayed {
						break
					}

					for _, cls := range terminator.SupportedClusters() {

						rec, ok := config.ClusterRegistry.Record(cls.StoreKey())
						if !ok {
							continue
						}
						if !rec.IsFree() {
                            // log.Warningf(" !IsFree UnMarkFree cluster %v", rec.StoreKey())

							config.ClusterRegistry.UnMarkFree(rec.StoreKey())
							rec.RefreshTimeout()
							continue
						}
						if !rec.IsTimeout() {
                            if rec.TimeOutAt.Sub(time.Now()).Seconds() < (time.Duration(interval * 10).Seconds()+60) {
                                log.Tracef(" Cluster %v will be terminated in %vs",  rec.StoreKey(), rec.TimeOutAt.Sub(time.Now()).Seconds())
                            }
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
							// reqClustersFailDescribed.Inc()
							log.Tracef("Failed to terminate cluster status '%v', failed with %v", rec.ClusterId, err)
						}
						metrics.FetchMetadataLatency.WithLabelValues("terminate_clusters",
							"single").Observe(float64(time.Now().Sub(start).Nanoseconds()))

					}
				}
				if !isDelayed {
					metrics.FetchMetadataLatency.WithLabelValues("terminate_clusters",
						"whole").Observe(float64(time.Now().Sub(start).Nanoseconds()))
				}

			}
		}
	}()

	numTermClusters := <-doneNumClusters

	log.Infof("Terminated %v clusters", numTermClusters)
	return nil
}
