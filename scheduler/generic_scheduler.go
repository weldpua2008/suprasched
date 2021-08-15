package scheduler

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weldpua2008/suprasched/core"
	internalcache "github.com/weldpua2008/suprasched/scheduler/internal"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

const (
	// minFeasibleClustersToFind is the minimum number of cluster that would be
	// scored in each scheduling cycle. This is a semi-arbitrary value to ensure
	// that a certain minimum of clusters are checked for feasibility. This in turn
	// helps ensure a minimum level of spreading.
	minFeasibleClustersToFind = 10
	// minFeasibleClustersPercentageToFind is the minimum percentage of cluster that
	// would be scored in each scheduling cycle. This is a semi-arbitrary value
	// to ensure that a certain minimum of clusters are checked for feasibility.
	// This in turn helps ensure a minimum level of spreading.
	minFeasibleClustersPercentageToFind = 5
)

// ErrNoClusterAvailable is used to describe the error that no cluster available to
// schedule jos.
var ErrNoClusterAvailable = fmt.Errorf("no clusters available to schedule jobs")

// ScheduleResult represents the result of one job scheduled. It will contain
// the final selected Cluster, along with the selected intermediate information.
type ScheduleResult struct {
	// Name of the scheduler suggest host
	SuggestedCluster core.UID
	// Number of clusters scheduler evaluated on one job scheduled
	EvaluatedClusters int
	// Number of feasible clusters on one job scheduled
	FeasibleClusters int
}

// ScheduleAlgorithm is an interface implemented by things that know how to schedule jobs
// onto machines.
type ScheduleAlgorithm interface {
	Schedule(context.Context, *core.Job) (scheduleResult ScheduleResult, err error)
}

type genericScheduler struct {
	cache                       *internalcache.Cache
	currentSnapshot             *internalcache.Snapshot
	percentageOfClustersToScore int32
	nextStartClusterIndex       int
}

// snapshot snapshots scheduler cache and cluster infos for all fit and priority
// functions.
func (g *genericScheduler) snapshot() error {
	if g.cache == nil {
		return fmt.Errorf("Cache is nil")
	}
	c := *g.cache
	// Used for all fit and priority funcs.
	return c.UpdateSnapshot(g.currentSnapshot)
}

// Schedule tries to schedule the given job to one of the clusters in the cluster list.
// If it succeeds, it will return the name of the cluster.
// If it fails, it will return a FitError error with reasons.
func (g *genericScheduler) Schedule(ctx context.Context, job *core.Job) (result ScheduleResult, err error) {

	trace := core.NewTracer(ctx, "Scheduling", logrus.WithFields(logrus.Fields{"namespace": job.Namespace, "name": job.Name}))
	if err := g.snapshot(); err != nil {
		return result, err
	}
	trace.Step("Snapshotting scheduler cache and cluster infos done")

	feasibleClusters, err := g.findClustersThatFitJob(ctx, job)
	if err != nil {
		return result, err
	}
	trace.Step("Computing predicates done")

	if len(feasibleClusters) == 0 {
		return result, ErrNoClusterAvailable
	}

	// When only one cluster after predicate, just use it.
	if len(feasibleClusters) == 1 {
		return ScheduleResult{
			SuggestedCluster:  feasibleClusters[0].UID,
			EvaluatedClusters: g.currentSnapshot.NumClusters(job.Namespace),
			FeasibleClusters:  1,
		}, nil
	}

	host := feasibleClusters[0]
	trace.Step("Prioritizing done")

	return ScheduleResult{
		SuggestedCluster:  host.UID,
		EvaluatedClusters: g.currentSnapshot.NumClusters(job.Namespace),
		FeasibleClusters:  len(feasibleClusters),
	}, err
}

// Filters the clusters to find the ones that fit the job
// TODO: implement algorithm
func (g *genericScheduler) findClustersThatFitJob(ctx context.Context, job *core.Job) ([]*core.Cluster, error) {
	for _, cl := range g.currentSnapshot.GetClustersFromNs(job.Namespace) {
		fmt.Println(job.Namespace)
		fmt.Println(cl.Name)

		return []*core.Cluster{&cl}, nil
	}

	return nil, ErrNoClusterAvailable
}

// NewGenericScheduler creates a genericScheduler object.
func NewGenericScheduler(
	cache *internalcache.Cache,
	currentSnapshot *internalcache.Snapshot,
	percentageOfClustersToScore int32) ScheduleAlgorithm {
	return &genericScheduler{
		cache:                       cache,
		currentSnapshot:             currentSnapshot,
		percentageOfClustersToScore: percentageOfClustersToScore,
	}
}
