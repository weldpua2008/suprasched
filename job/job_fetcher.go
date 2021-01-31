package job

import (
	"context"
	"fmt"
	config "github.com/weldpua2008/suprasched/config"
	metrics "github.com/weldpua2008/suprasched/metrics"
	model "github.com/weldpua2008/suprasched/model"
	"strings"
	"time"
)

// GetJobsFetchersFromSection returns multiple JobsFetcher from configuration file.
// By default http communicator will be used.
// Example YAML config for `section` that will return new `RestJobsFetcher`:
//     jobs:
//         fetch:
//             my_communicator:
//                 type: "HTTP"
//             -:
//                 type: "HTTP"
func GetJobsFetchersFromSection(section string) ([]JobsFetcher, error) {

	def := make(map[string]string)
	fetchersCfgs := config.GetSliceStringMapStringTemplatedDefault(section, config.CFG_PREFIX_JOBS_FETCHER, def)
	res := make([]JobsFetcher, 0)
	for _, comm := range fetchersCfgs {
		if comm == nil {
			continue
		}
		describerType := ConstructorsJobsFetcherRest
		if tmpDescriberType, ok := comm["type"]; ok {
			describerType = tmpDescriberType
		}
		k := strings.ToUpper(describerType)
		if typeStruct, ok := Constructors[k]; ok {
			describerInstance, err := typeStruct.constructor(section)
			if err != nil {
				log.Tracef("Can't get fetcher %v", err)
				continue
			}
			res = append(res, describerInstance)

		}

	}
	if len(res) > 0 {

		return res, nil
	}
	return nil, fmt.Errorf("%w in section %s.\n", ErrNoSuitableJobsFetcher, section)

}

// StartFetchJobs goroutine for getting jobs from API with internal
// exists on kill
func StartFetchJobs(ctx context.Context, jobs chan bool, interval time.Duration) error {
	fetchers, err := GetJobsFetchersFromSection(config.CFG_PREFIX_JOBS)
	if err != nil || fetchers == nil || len(fetchers) == 0 {
		return fmt.Errorf("Failed to start StartFetchJobs %v", err)
	}

	doneNumJobs := make(chan int, 1)
	log.Infof("Pulling Jobs Metadata with delay %v", interval)
	tickerPullJobs := time.NewTicker(interval)
	defer func() {
		close(jobs)
		tickerPullJobs.Stop()
	}()

	go func() {
		cntr := 0
		for {
			select {
			case <-ctx.Done():
				doneNumJobs <- cntr
				log.Debug("Jobs fetch finished [ SUCCESSFULLY ]")
				return
			case <-tickerPullJobs.C:
				start := time.Now()

				for _, fetcher := range fetchers {

					jobsSlice, err := fetcher.Fetch()
					if err != nil {
						log.Tracef("Fetch job metadata '%v', but failed with %v", jobsSlice, err)
						continue
					}
					for _, j := range jobsSlice {
						storeKey := j.StoreKey()
						if len(storeKey) < 1 {
							continue
						}
						var topic string
						rec, exist := config.JobsRegistry.Record(storeKey)
						if !exist {
							// Possible broken Job
							if !config.JobsRegistry.Add(j) {
								continue
							} else {
								rec, exist = config.JobsRegistry.Record(storeKey)
								if !exist {
									continue
								}
								topic = config.TOPIC_JOB_CREATED
								cntr += 1
								jobsFetched.Inc()
							}
						} else if rec.IsInTransition() {
							continue
						} else if rec.UpdateStatus(j.Status) {
							topic = strings.ToLower(fmt.Sprintf("job.%v", rec.Status))
							// if model.IsTerminalStatus(rec.GetStatus()) && rec.ClusterType != model.CLUSTER_TYPE_ON_DEMAND {

						}

						if len(topic) > 1 {
							_, err := config.Bus.Emit(ctx, topic, rec.EventMetadata())
							if err != nil {
								log.Tracef("%v", err)
							}

							if (rec != nil) && (rec.ClusterType != model.CLUSTER_TYPE_ON_DEMAND) && (len(rec.ClusterId) > 0) {
								// log.Warningf("Send update for %v",  rec.GetClusterStoreKey())
								clusterEventMetadata := map[string]string{"StoreKey": rec.GetClusterStoreKey()}
								_, err := config.Bus.Emit(ctx, config.TOPIC_CLUSTER_REFRESH_TIMEOUT, clusterEventMetadata)
								if err != nil {
									log.Tracef("%v", err)
								}

								if model.IsTerminalStatus(rec.Status) {
									_, err := config.Bus.Emit(ctx, config.TOPIC_CLUSTER_IS_EMPTY, clusterEventMetadata)
									if err != nil {
										log.Tracef("%v", err)
									}
								}

							}
						}
						metrics.FetchMetadataLatency.WithLabelValues("fetch_jobs",
							"single").Observe(float64(time.Now().Sub(start).Nanoseconds()))

					}

				}
				metrics.FetchMetadataLatency.WithLabelValues("fetch_jobs",
					"whole").Observe(float64(time.Now().Sub(start).Nanoseconds()))
			}
		}
	}()

	numSentClusters := <-doneNumJobs

	log.Infof("Fetched %v jobs", numSentClusters)
	return nil
}
