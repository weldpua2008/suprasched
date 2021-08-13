package cache

import "github.com/weldpua2008/suprasched/core"

type Cache interface {
	// ClusterCount returns the number of clusters in the cache.
	// DO NOT use outside of tests.
	//ClusterCount() int

	// JobCount returns the number of jobs in the cache.
	// DO NOT use outside of tests.
	//JobCount() (int, error)

	// AddJob either confirms a job if it's assumed, or adds it back if it's expired.
	// If added back, the job's information would be added again.
	//AddJob(job *core.Job) error

	// UpdateJob removes oldJob's information and adds newJob's information.
	//UpdateJob(oldJob, newJob *core.Job) error

	// RemoveJob removes a job. The job's information would be subtracted from assigned cluster.
	//RemoveJob(job *core.Job) error

	// GetJob returns the job from the cache with the same namespace and the
	// same name of the specified job.
	//GetJob(job *core.Job) (*core.Job, error)

	// IsAssumedJob returns true if the job is assumed and not expired.
	//IsAssumedJob(job *core.Job) (bool, error)

	// AddCluster adds overall information about cluster.
	AddCluster(cluster *core.Cluster) error

	// UpdateCluster updates overall information about cluster.
	//UpdateCluster(oldCluster, newCluster *core.Cluster) error

	// RemoveCluster removes overall information about cluster.
	RemoveCluster(cluster *core.Cluster) error
	// UpdateSnapshot updates the passed currSnapshot to the current contents of Cache.
	// The node info contains aggregated information of pods scheduled (including assumed to be)
	// on this node.
	// The snapshot only includes Nodes that are not deleted at the time this function is called.
	// nodeinfo.Node() is guaranteed to be not nil for all the nodes in the snapshot.
	UpdateSnapshot(currSnapshot *Snapshot) error

	// Dump produces a dump of the current cache.
	Dump() *Dump
}
