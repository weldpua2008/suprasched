package job

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoSuitableJobsFetcher = errors.New("No suitable communicator found")
	// internal
	jobsFetched = promauto.NewCounter(prometheus.CounterOpts{
		Name: "suprasched_jobs_fetch_total",
		Help: "The total number of fetched jobs",
	})

	log = logrus.WithFields(logrus.Fields{"package": "communicator"})
)
