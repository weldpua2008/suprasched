package poc

import (
	"context"
	config "github.com/weldpua2008/suprasched/config"
	handlers "github.com/weldpua2008/suprasched/handlers"
	model "github.com/weldpua2008/suprasched/model"
	"go.uber.org/goleak"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func BenchmarkParallelHandler(b *testing.B) {
	b.SkipNow()
	config.JobsRegistry = model.NewRegistry()
	config.ClusterRegistry = model.NewClusterRegistry()
	cluster_types := []string{"type1", "type2", "type3"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are getting the kill signal or exit
	handlers.Start("ParallelHandler", ParallelHandler, config.TOPIC_JOB_CREATED)
	defer handlers.Stop("ParallelHandler")

	if err := InitializeJobs((b.N + 10),
		cluster_types,
		[]string{"PENDING", "RUNNING", "FAILED"},
	); err != nil {
		panic(err)
	}
	jobs := config.JobsRegistry.All()
	assigned_jobs = 0
	empty_cntr = 0
	start := time.Now()
	var wg sync.WaitGroup
	b.ResetTimer()
	// b.Run("create jobs", func(b *testing.B) {
	for _, j := range jobs {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
			}()
			topic := config.TOPIC_JOB_CREATED
			_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
			if err != nil {
				log.Tracef("%v", err)
			}
		}(&wg)
	}
	wg.Wait()
	// })
	if empty_cntr < (b.N + 10) {
		b.Errorf("Expect assigned_jobs %v > %v", empty_cntr, (b.N + 10))

	}

	if empty_cntr > len(jobs) {
		b.Errorf("Expect assigned_jobs %v < %v", empty_cntr, len(jobs))
	}
	if empty_cntr > config.JobsRegistry.Len() {
		b.Errorf("Expect assigned_jobs %v <  config.JobsRegistry.Len() %v", empty_cntr, config.JobsRegistry.Len())
	}
	log.Warningf("BenchmarkParallelHandler jobs %v took %s %s per job", len(jobs), time.Since(start), time.Duration(int(int64(time.Since(start).Nanoseconds())/int64(len(jobs))))*time.Nanosecond)

	// b.TimerStop()
	start = time.Now()

	for _, j := range jobs {
		topic := config.TOPIC_JOB_CREATED
		_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
		if err != nil {
			log.Tracef("%v", err)
		}
	}
	log.Warningf("BenchmarkParallelHandler jobs %v took %s %s per job", len(jobs), time.Since(start), time.Duration(int(int64(time.Since(start).Nanoseconds())/int64(len(jobs))))*time.Nanosecond)

}

func BenchmarkEmptyHandler(b *testing.B) {
	b.SkipNow()
	config.JobsRegistry = model.NewRegistry()
	config.ClusterRegistry = model.NewClusterRegistry()
	cluster_types := []string{"type1", "type2", "type3"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are getting the kill signal or exit
	handlers.Start("EmptyHandler", EmptyHandler, config.TOPIC_JOB_CREATED)
	defer handlers.Stop("EmptyHandler")

	if err := InitializeJobs(b.N,
		cluster_types,
		[]string{"PENDING", "RUNNING", "FAILED"},
	); err != nil {
		panic(err)
	}
	jobs := config.JobsRegistry.All()
	assigned_jobs = 0
	empty_cntr = 0
	start := time.Now()
	var wg sync.WaitGroup
	b.ResetTimer()
	// b.Run("create jobs", func(b *testing.B) {
	for _, j := range jobs {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer func() {
				wg.Done()
			}()
			topic := config.TOPIC_JOB_CREATED
			_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
			if err != nil {
				log.Tracef("%v", err)
			}
		}(&wg)
	}
	wg.Wait()
	// })
	if empty_cntr < (b.N + 1) {
		b.Errorf("Expect assigned_jobs %v > %v", empty_cntr, (b.N + 1))

	}
	// config.JobsRegistry.Len()
	if empty_cntr > len(jobs) {
		b.Errorf("Expect assigned_jobs %v < %v", empty_cntr, len(jobs))
	}
	if empty_cntr > config.JobsRegistry.Len() {
		b.Errorf("Expect assigned_jobs %v <  config.JobsRegistry.Len() %v", empty_cntr, config.JobsRegistry.Len())
	}
	// log.Warningf(" took %s", time.Since(start))
	log.Warningf("Empty took %s %s per job", time.Since(start), time.Duration(int(int64(time.Since(start).Nanoseconds())/int64(len(jobs))))*time.Nanosecond)

}

func BenchmarkMultipleTopicsHandler11(b *testing.B) {
	benchmarkMultipleTopicsHandler(1, 1, b)
}
func BenchmarkMultipleTopicsHandlerJ100C10(b *testing.B) {
	benchmarkMultipleTopicsHandler(100, 10, b)
}

func BenchmarkMultipleTopicsHandlerJ1000000C1(b *testing.B) {
	benchmarkMultipleTopicsHandler(1000000, 1, b)
}

func BenchmarkMultipleTopicsHandlerJ1000000C1000(b *testing.B) {
	benchmarkMultipleTopicsHandler(1000000, 1000, b)
}

func benchmarkMultipleTopicsHandler(num_jobs int, num_cluster int, b *testing.B) {
	b.SkipNow()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are getting the kill signal or exit
	// "cluster.term.*"
	handlers.Start("AssignNewFreeCluster", AssignNewFreeCluster, config.TOPIC_JOB_CREATED)
	defer handlers.Stop("AssignNewFreeCluster")
	handlers.Start("FreeJobFromCluster", FreeJobFromCluster, config.TOPIC_JOB_SUCCESS)
	defer handlers.Stop("FreeJobFromCluster")

	handlers.Start("DeleteCluster", DeleteCluster, config.TOPIC_CLUSTER_TERMINATED)
	defer handlers.Stop("DeleteCluster")

	config.JobsRegistry = model.NewRegistry()
	config.ClusterRegistry = model.NewClusterRegistry()
	cluster_types := []string{"type1", "type2", "type3"}
	if err := InitializeClusters(num_cluster,
		cluster_types,
		[]string{"PENDING", "RUNNING"},
	); err != nil {
		panic(err)
	}
	if err := InitializeJobs(num_jobs,
		cluster_types,
		[]string{"PENDING", "RUNNING", "FAILED"},
	); err != nil {
		panic(err)
	}
	jobs := config.JobsRegistry.All()
	n := config.JobsRegistry.Len()
	assigned_jobs = 0
	start := time.Now()
	new_jobs_created := 0
	// b.ResetTimer()
	b.Run("create jobs", func(b *testing.B) {

		for _, j := range config.JobsRegistry.All() {
			new_jobs_created += 1
			topic := config.TOPIC_JOB_CREATED
			_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
			if err != nil {
				log.Tracef("%v", err)
			}
		}
	})

	b.Run("check empty jobs", func(b *testing.B) {

		for _, j := range config.JobsRegistry.All() {
			if len(j.GetClusterStoreKey()) == 0 {
				b.Errorf("Expect ClusterStoreKey len > 0 %v", j)
			}
		}
		for _, empty_cl := range config.ClusterRegistry.AllEmpty() {
			config.ClusterRegistry.Delete(empty_cl.StoreKey())
		}

		if assigned_jobs > n {
			b.Errorf("Expect assigned_jobs %v < %v", assigned_jobs, n)
		}

	})
	// if assigned_jobs >  config.JobsRegistry.Len() {
	//     b.Errorf("Expect assigned_jobs %v <  config.JobsRegistry.Len() %v", assigned_jobs, config.JobsRegistry.Len())
	// }
	log.Warningf("assigned jobs %v took %s %s per job", assigned_jobs, time.Since(start), time.Duration(int(int64(time.Since(start).Nanoseconds())/int64(len(jobs))))*time.Nanosecond)
	start = time.Now()
	// for config.JobsRegistry.Len() > 0 {
	if len(config.JobsRegistry.All()) < 1 {
		b.Errorf("Expect config.JobsRegistry.All() %v > 0", len(config.JobsRegistry.All()))
	}

	for _, j := range jobs {
		log.Tracef("%v", j)
	}
	b.Run("success jobs", func(b *testing.B) {

		for _, j := range jobs {
			topic := config.TOPIC_JOB_SUCCESS
			j.UpdateStatus("SUCCESS")
			_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
			if err != nil {
				log.Tracef("%v", err)
			}
		}
	})
	// for _, empty_cl := range config.ClusterRegistry.AllEmpty() {
	//     log.Warningf("%v", empty_cl)
	// }
	// for _, empty_cl := range config.ClusterRegistry.AllEmpty() {
	//     config.ClusterRegistry.Delete(empty_cl.StoreKey())
	// }

	if config.JobsRegistry.Len() > 0 {
		b.Errorf("Expect config.JobsRegistry.Len()  0 != %v", config.JobsRegistry.Len())
	}

	if config.ClusterRegistry.Len() > 0 {
		b.Errorf("Expect config.ClusterRegistry.Len()  0 != %v", config.ClusterRegistry.Len())

	}
	log.Warningf("cleanup jobs %v took %s %s per job", assigned_jobs, time.Since(start), time.Duration(int(int64(time.Since(start).Nanoseconds())/int64(len(jobs))))*time.Nanosecond)

}
