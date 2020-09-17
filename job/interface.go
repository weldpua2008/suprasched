package job

import (
	// "context"
	model "github.com/weldpua2008/suprasched/model"
)

// A JobsFetcher is the interface used to communicate with APIs
// that will eventually return metadata. JobsFetcher
// allow you to get information from remote APi, databases, etc.
//
// JobsFetcher must be safe for concurrency, meaning multiple calls to
// any method may be called at the same time.
type JobsFetcher interface {
	// Fetch metadata from remote storage
	Fetch() ([]*model.Job, error)
}
