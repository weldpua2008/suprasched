package model

import (
	"github.com/sirupsen/logrus"
)

var (
	log           = logrus.WithFields(logrus.Fields{"package": "model"})
	previousLevel logrus.Level
)

const (
	JOB_STATUS_PENDING     = "pending"
	JOB_STATUS_IN_PROGRESS = "in_progress"
	JOB_STATUS_SUCCESS     = "success"
	JOB_STATUS_ERROR       = "error"
	JOB_STATUS_CANCELED    = "canceled"
	JOB_STATUS_QUEUED      = "queued"

	CLUSTER_STATUS_STARTING               = "starting"
	CLUSTER_STATUS_BOOTSTRAPPING          = "bootstraping"
	CLUSTER_STATUS_RUNNING                = "running"
	CLUSTER_STATUS_WAITING                = "waiting"
	CLUSTER_STATUS_TERMINATING            = "terminating"
	CLUSTER_STATUS_TERMINATED             = "terminated"
	CLUSTER_STATUS_TERMINATED_WITH_ERRORS = "terminated_with_errors"

	CLUSTER_TYPE_ON_DEMAND = "on-demand"
	CLUSTER_TYPE_EMR       = "EMR"
)

func init() {
	previousLevel = logrus.GetLevel()
}
