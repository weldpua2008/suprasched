package scheduler

import (
	"context"
	"github.com/weldpua2008/suprasched/core"
	internalcache "github.com/weldpua2008/suprasched/scheduler/internal"
	"testing"
	"time"
)

var testNamespace core.Namespace = "testNamespace"
var otherNamespace core.Namespace = "othertNamespace"

func TestGenericScheduler(t *testing.T) {
	cl1 := core.NewCluster(
		"test-cluster",
		testNamespace,
		core.ClusterSpec{},
		core.ClusterStatus{},
	)
	cl2 := core.NewCluster(
		"other--test-cluster",
		otherNamespace,
		core.ClusterSpec{},
		core.ClusterStatus{},
	)
	tests := []struct {
		name                 string
		cache                internalcache.Cache
		job                  core.Job
		clusters             []*core.Cluster
		numFeasibleClusters  int
		numEvaluatedClusters int
		SuggestedCluster     core.UID
		err                  error
	}{
		{
			name:                 "test one job and one cluster",
			cache:                internalcache.New(time.Minute),
			job:                  core.NewJob("test", testNamespace),
			numFeasibleClusters:  1,
			numEvaluatedClusters: 1,
			clusters:             []*core.Cluster{&cl1},
			SuggestedCluster:     cl1.UID,
		},
		{
			name:                 "test one job and two clusters",
			cache:                internalcache.New(time.Minute),
			job:                  core.NewJob("test1", testNamespace),
			numFeasibleClusters:  1,
			numEvaluatedClusters: 1,
			clusters:             []*core.Cluster{&cl1, &cl2},
			SuggestedCluster:     cl1.UID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, cl := range test.clusters {
				_ = test.cache.AddCluster(cl)
			}

			scheduler := NewGenericScheduler(test.cache, new(internalcache.Snapshot), 0)
			ctx := context.Background()
			got, err := scheduler.Schedule(ctx, &test.job)
			if got.EvaluatedClusters != test.numEvaluatedClusters {
				t.Errorf("Expects %v got %v EvaluatedClusters", test.numEvaluatedClusters, got.EvaluatedClusters)
			}
			if got.SuggestedCluster != test.SuggestedCluster {
				t.Errorf("Expects %v got %v SuggestedCluster in %v", test.SuggestedCluster, got.SuggestedCluster, test.clusters)
			}
			if got.FeasibleClusters != test.numFeasibleClusters {
				t.Errorf("Expects %v got %v FeasibleClusters", test.numFeasibleClusters, got.FeasibleClusters)
			}
			if err != test.err {
				t.Errorf("Unexpected error %v != %v ", err, test.err)
			}

			//}
		})
	}
}
