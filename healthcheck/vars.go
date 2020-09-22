package healthcheck

import (
	"errors"
	"github.com/sirupsen/logrus"
	// model "github.com/weldpua2008/suprasched/model"
)

var (
	log                      = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrServerListenError = errors.New("Error HTTP server ListenAndServe:")
)
