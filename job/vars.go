package job

import (
	"errors"
	"github.com/sirupsen/logrus"
	// model "github.com/weldpua2008/suprasched/model"
)

var (
	log                      = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrNoSuitableJobsFetcher = errors.New("No suitable communicator found")
)
