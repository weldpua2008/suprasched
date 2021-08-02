package job

import (
	"context"
	config "github.com/weldpua2008/suprasched/config"
	metrics "github.com/weldpua2008/suprasched/metrics"
	model "github.com/weldpua2008/suprasched/model"
	"strings"
	"time"
)

func CancelTimeoutJobs(ctx context.Context, jobs chan bool, interval time.Duration, pendingTimeOut time.Duration) error {
	doneCancelJobs := make(chan int, 1)
	log.Infof("Checks Jobs Metadata for Timeout with delay %v Job pendingTimeOut %v", interval, pendingTimeOut)
	tickerPullJobs := time.NewTicker(interval)
	defer func() {
		close(jobs)
		tickerPullJobs.Stop()
	}()

	go func() {
		counter := 0
		for {
			select {
			case <-ctx.Done():
				doneCancelJobs <- counter
				log.Debug("Jobs timeout cancellation finished [ SUCCESSFULLY ]")
				return
			case <-tickerPullJobs.C:
				start := time.Now()

				for _, j := range config.JobsRegistry.All() {
					if len(j.StoreKey()) < 1 {
						continue
					}
					if j.IsInTransition() {
						continue
					}
					if strings.EqualFold(j.Status, model.JOB_STATUS_PENDING) {
						counter += 1
						pendingDelay := j.CreateAt.Add(pendingTimeOut)
						if time.Now().After(pendingDelay) {
							_, err := config.Bus.Emit(ctx, config.TOPIC_JOB_FORCE_TIMEOUT, j.EventMetadata())
							if err != nil {
								log.Tracef("%v", err)
							}
							// log.Tracef("Job  %v timeout %v", j.StoreKey(),pending_delay)

						}
					}

					metrics.FetchMetadataLatency.WithLabelValues("timeout_jobs",
						"single").Observe(float64(time.Since(start).Nanoseconds()))
				}

				metrics.FetchMetadataLatency.WithLabelValues("timeout_jobs",
					"whole").Observe(float64(time.Since(start).Nanoseconds()))
			}
		}

	}()

	numSentClusters := <-doneCancelJobs

	log.Infof("Canceled %v jobs", numSentClusters)
	return nil
}
