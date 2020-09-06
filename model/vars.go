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
)

func init() {
	previousLevel = logrus.GetLevel()
}
