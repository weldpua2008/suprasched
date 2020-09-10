package model

import (
	"sync"
	"time"
)

// NewRegistry returns a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		all: make(map[string]*Job),
	}
}

// Registry holds all Job Records.
type Registry struct {
	all map[string]*Job
	mu  sync.RWMutex
}

// Add a job.
// Returns false on duplicate or invalid job id.
func (r *Registry) Add(rec *Job) bool {

	if rec == nil || rec.StoreKey() == "" {
        // log.Tracef("False %v", rec)
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.all[rec.StoreKey()]; ok {
        // log.Tracef("False already exist %v", rec)

		return false
	}

	r.all[rec.StoreKey()] = rec
	return true
}

// Len returns length of registry.
func (r *Registry) Len() int {
	r.mu.RLock()
	c := len(r.all)
	r.mu.RUnlock()
	return c
}

// Delete a job by job ID.
// Return false if record does not exist.
func (r *Registry) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.all[id]
	if !ok {
		return false
	}
	delete(r.all, id)
	return true
}

// Cleanup process for the registry with batch only locked.
// Return number of cleaned jobs.
func (r *Registry) Cleanup() (num int) {
	n := r.Len()
	slice := make([]string, n)
	i := 0
	for k, _ := range r.all {
		if i > n {
			slice = append(slice, k)
		} else {
			slice[i] = k

		}
		i++
	}

	batch := 20
	for i := 0; i < len(slice); i += batch {
		j := i + batch
		if j > len(slice) {
			j = len(slice)
		}

		// fmt.Println(slice[i:j]) // Process the batch.
		numBatch := r.CleanupBatch(slice[i:j])
		num += numBatch
	}

	return num
}

// CleanupBatch by job TTR.
// Return number of cleaned jobs.
func (r *Registry) CleanupBatch(slice []string) (num int) {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	// for k, v := range r.all {
	for _, k := range slice {
		if v, ok := r.all[k]; ok {
			end := v.StartAt.Add(time.Duration(v.TTR) * time.Millisecond)
			if (v.TTR > 0) && (now.After(end)) {
				if !IsTerminalStatus(v.Status) {
					if err := v.Cancel(); err != nil {
						log.Debugf("failed cancel job %s %v StartAt %v", v.Id, err, v.StartAt)
					} else {
						log.Tracef("successfully canceled job %s StartAt %v, TTR %v msec", v.Id, v.StartAt, v.TTR)
					}
				}
				delete(r.all, k)
				num += 1
			}

		}

	}
	return num
}

// GracefullShutdown is used when we stop the Registry.
// cancel all running & pending job
// return false if we can't cancel any job
func (r *Registry) GracefullShutdown() bool {
	r.Cleanup()
	r.mu.Lock()
	defer r.mu.Unlock()
	failed := false
	// log.Debug("start GracefullShutdown")
	// for k, v := range r.all {
	// 	if !IsTerminalStatus(v.Status) {
	// 		if err := v.Cancel(); err != nil {
	// 			log.Debug(fmt.Sprintf("failed cancel job %s %v", v.Id, err))
	// 			failed = true
	// 		} else {
	// 			log.Debug(fmt.Sprintf("successfully canceled job %s", v.Id))
	// 		}
	// 	}
	// 	delete(r.all, k)
	// }
	return failed
}

// Record fetch job by Job ID.
// Follows comma ok idiom
func (r *Registry) Record(jid string) (*Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rec, ok := r.all[jid]; ok {
		return rec, true
	}

	return nil, false
}
