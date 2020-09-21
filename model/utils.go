package model

import (
	"strings"
)

// IsTerminalStatus returns true if status is terminal:
// Jobs
// - Failed
// - Canceled
// - Successful
// Cluster
// - TERMINATING
// - TERMINATED_WITH_ERRORS
// - TERMINATED

func IsTerminalStatus(status string) bool {
	switch strings.ToLower(status) {
	case CLUSTER_STATUS_TERMINATED, CLUSTER_STATUS_TERMINATED_WITH_ERRORS, CLUSTER_STATUS_TERMINATING:
		// log.Tracef("IsTerminalStatus %s true", status)
		return true
	case JOB_STATUS_ERROR, JOB_STATUS_SUCCESS, JOB_STATUS_FAILED, JOB_STATUS_CANCELED:
		// log.Tracef("IsTerminalStatus %s true", status)
		return true
	}
	// log.Tracef("IsTerminalStatus %s false", status)
	return false
}
