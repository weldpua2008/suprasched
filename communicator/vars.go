package communicator

import (
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	log                       = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrNoSuitableCommunicator = errors.New("No suitable communicator found")
)
