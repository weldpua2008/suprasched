package metrics

import (
	"errors"
	"github.com/sirupsen/logrus"
	// "net/http"
	"sync"
)

var (
	log                  = logrus.WithFields(logrus.Fields{"package": "communicator"})
	ErrServerListenError = errors.New("HTTP server ListenAndServe:")
	listenServersStore   = make(map[string]*SrvSession)
	mu                   sync.RWMutex
)
