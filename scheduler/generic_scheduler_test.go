package scheduler

import (
	"context"
	"github.com/weldpua2008/suprasched/core"
	internalcache "github.com/weldpua2008/suprasched/scheduler/internal"
	"testing"
	"time"
)

var testNamespace core.Namespace = "testNamespace"

func TestGenericScheduler(t *testing.T) {

	tests := []struct {
		name     string
		cache    internalcache.Cache
		job      core.Job
		clusters []core.Cluster
	}{
		{
			name:  "test one job and one cluster",
			cache: internalcache.New(time.Minute),
			job:   core.NewJob("test", testNamespace),
			clusters: []core.Cluster{core.NewCluster(
				"test-cluster", testNamespace,
				core.ClusterSpec{},
				core.ClusterStatus{},
			)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i := 0; i < 10; i++ {
				for _, cl := range test.clusters {
					_ = test.cache.AddCluster(&cl)
				}
				scheduler := NewGenericScheduler(test.cache, new(internalcache.Snapshot), 0)
				ctx := context.Background()
				got, err := scheduler.Schedule(ctx, &test.job)
				t.Logf("got %v, err %v", got, err)

				if err != nil {
					t.Error("Unexpected non-error")
				}
				for _, cl := range test.clusters {
					_ = test.cache.RemoveCluster(&cl)
				}
			}
		})
	}
}
