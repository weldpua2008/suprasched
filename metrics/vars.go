package metrics

import (
	"errors"
	"github.com/sirupsen/logrus"
	// "net/http"
	"sync"
)

var (
	log                  = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrServerListenError = errors.New("Error HTTP server ListenAndServe:")
	listenServersStore   = make(map[string]*SrvSession, 0)
	mu                   sync.RWMutex
)