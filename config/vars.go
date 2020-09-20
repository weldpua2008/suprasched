package config
const (
	// ProjectName defines project name
	ProjectName = "suprasched"
	// CFG_PREFIX_JOB for the config
	CFG_PREFIX_JOBS         = "jobs"
	CFG_PREFIX_JOBS_FETCHER = "fetch"

	//CFG_PREFIX_COMMUNICATOR defines parameter in the config for Communicators
	CFG_PREFIX_COMMUNICATOR            = "communicator"
	CFG_PREFIX_COMMUNICATORS           = "communicators"
	CFG_PREFIX_CLUSTER_SUPPORTED_TYPES = "supported"

	CFG_PREFIX_CLUSTER          = "cluster"
    CFG_PREFIX_UPDATE           = "update"
	CFG_PREFIX_FETCHER          = "fetch"
	CFG_PREFIX_DESCRIBERS       = "describe"
	CFG_COMMUNICATOR_PARAMS_KEY = "params"
	CFG_INTERVAL_PARAMETER      = "interval"

    // Event Matchers.
    MATCHER_CLUSTER_TERMINATING = "cluster.term.*"

    // Cluster related topics.
	TOPIC_CLUSTER_CREATED       = "cluster.created"
	TOPIC_CLUSTER_STARTING      = "cluster.starting"
	TOPIC_CLUSTER_BOOTSTRAPPING = "cluster.bootstraping"
	TOPIC_CLUSTER_RUNNING       = "cluster.running"
	TOPIC_CLUSTER_WAITING       = "cluster.waiting"
	TOPIC_CLUSTER_TERMINATING            = "cluster.terminating"
	TOPIC_CLUSTER_TERMINATED             = "cluster.terminated"
	TOPIC_CLUSTER_TERMINATED_WITH_ERRORS = "cluster.terminated_with_errors"
    // Jobs related topics.
	TOPIC_JOB_CANCELED = "job.canceled"
	TOPIC_JOB_CREATED  = "job.created"
	TOPIC_JOB_STARTING = "job.starting"
	TOPIC_JOB_PENDING  = "job.pending"
	TOPIC_JOB_RUNNING  = "job.running"
	TOPIC_JOB_FAILED                 = "job.failed"
	TOPIC_JOB_SUCCEEDED              = "job.succeeded"
	TOPIC_JOB_SUCCESS                = "job.success"
	TOPIC_JOB_TERMINATING            = "job.terminating"
	TOPIC_JOB_TERMINATED             = "job.terminated"
	TOPIC_JOB_TERMINATED_WITH_ERRORS = "job.terminated_with_errors"

)
