package model

import (
	"fmt"
	"testing"
	// "time"
)

func BenchmarkClusterRegistryAdd(b *testing.B) {
	r := NewClusterRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", b.N))
		r.Add(cluster)
	}
}

func BenchmarkClusterRegistryAddByType(b *testing.B) {
	r := NewClusterRegistry()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", b.N))
		cluster.ClusterType = "ClusterType"
		r.Add(cluster)

		if len(r.byType[cluster.ClusterType]) != i {
			b.Errorf("Expect r.byType[%v] to add cluster", cluster.ClusterType)

		}
	}
}

func BenchmarkClusterRegistryDeleteByType(b *testing.B) {
	r := NewClusterRegistry()

	ClusterType := "ClusterType"
	for i := 0; i < b.N; i++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", b.N))
		cluster.ClusterType = ClusterType
		r.Add(cluster)

		if len(r.byType[cluster.ClusterType]) != i {
			b.Errorf("Expect r.byType[%v] to add cluster", cluster.ClusterType)

		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Delete(fmt.Sprintf("cluster-%v", b.N))
		if len(r.byType[ClusterType]) == b.N {
			b.Errorf("Expect r.byType[%v] to delete cluster", ClusterType)

		}
	}
}

// func BenchmarkClusterRegistryCleanUp(b *testing.B) {
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		r := NewClusterRegistry()
// 		for ii := 0; ii < 100; ii++ {
// 			cluster := NewCluster(fmt.Sprintf("cluster-%v", b.N))
// 			r.Add(cluster)
// 			r.Cleanup()
// 		}
// 	}
// }

func TestClusterRegistryAddNoDuplicateCluster(t *testing.T) {
	r := NewClusterRegistry()
	for ii := 0; ii < 100; ii++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", ii))
		if !r.Add(cluster) {
			t.Errorf("Expect to add cluster")
		}
		for j := 0; j < 10; j++ {
			if r.Add(cluster) {
				t.Errorf("Expect not to add cluster")

			}
		}

	}
}

func TestClusterRegistryLen(t *testing.T) {
	r := NewClusterRegistry()
	num := 100
	for ii := 0; ii < num; ii++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", ii))
		if !r.Add(cluster) {
			t.Errorf("Expect to add cluster")
		}
	}
	if r.Len() != num {
		t.Errorf("Expect %v got length %v", num, r.Len())
	}
}

func TestClusterRegistryDelete(t *testing.T) {
	r := NewClusterRegistry()
	num := 100
	for ii := 0; ii < num; ii++ {
		cluster := NewCluster(fmt.Sprintf("cluster-%v", ii))
		if !r.Add(cluster) {
			t.Errorf("Expect to add cluster")
		}
		if !r.Delete(cluster.StoreKey()) {
			t.Errorf("Expect to delete cluster")
		}
		if r.Delete(cluster.StoreKey()) {
			t.Errorf("Expect the cluster to be already deleted")
		}

	}
	if r.Len() != 0 {
		t.Errorf("Expect %v got length %v", num, r.Len())
	}
}

// func TestClusterRegistryCleanup(t *testing.T) {
// 	r := NewClusterRegistry()
// 	num := 101
// 	for ii := 0; ii < num; ii++ {
// 		cluster := NewCluster(fmt.Sprintf("cluster-%v", ii))
// 		cluster.TTR = 10000
// 		// no cancelation flow on cleanup
// 		// right now it won't execute something
// 		cluster.Status = JOB_STATUS_CANCELED
// 		if len(cluster.StoreKey()) == 0 {
// 			t.Errorf("cluster.StoreKey size %v > 0", len(cluster.StoreKey()))
// 		}
// 		if r.Len() > 0 {
// 			t.Errorf("Expect registry size %v == 0", r.Len())
// 		}
// 		if !r.Add(cluster) {
// 			t.Errorf("Expect to add cluster")
// 		}
// 		if r.Len() == 0 {
// 			t.Errorf("Expect registry size %v > 0", r.Len())
// 		}
// 		n := r.Len()
//
// 		if (r.Cleanup() > 0) || (r.Len() != n) {
// 			t.Errorf("Expect no cluster to be already deleted by Cleanup")
// 		}
// 		cluster.StartAt = time.Now().Add(time.Duration(-10001) * time.Millisecond)
// 		if (r.Cleanup() == 0) || (r.Len() == n) {
// 			t.Errorf("Expect Cluster to be deleted by Cleanup due to TTR")
// 		}
//
// 	}
// 	if r.Len() != 0 {
// 		t.Errorf("Expect %v got length %v", num, r.Len())
// 	}
// }
