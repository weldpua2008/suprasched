package communicator

import (
	"errors"
	"github.com/sirupsen/logrus"
)

var (
	log                        = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrNoSuitableCommunicator  = errors.New("No suitable communicator found")
	ErrFailedSendRequest       = errors.New("Failed to send request")
	ErrFailedReadResponseBody  = errors.New("Failed to read response body")
	ErrFailedUnmarshalResponse = errors.New("Cannot unmarshal response")
	ErrFailedMarshalRequest    = errors.New("Cannot marshal request")
)
