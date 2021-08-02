package handlers

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	log                   = logrus.WithFields(logrus.Fields{"package": "handlers"})
	ErrFailedSendRequest  = errors.New("Failed to send request")
	ErrFailedUpdateStatus = errors.New("Failed to update status")
	ErrNoClusterFound     = errors.New("No Cluster found")
)
