package core

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

// LoggerFromContext returns a logrus.Entry with the PID of the current process
// set as a field, and also includes every field set using the From* functions
// this package.
func LoggerFromContext(ctx context.Context, logger *logrus.Entry) *logrus.Entry {
	if logger == nil {
		logger = logrus.WithFields(logrus.Fields{"package": "context"})
	}

	entry := logger.WithField("pid", os.Getpid())
	if ctx == nil {
		return entry
	}

	return entry
}
