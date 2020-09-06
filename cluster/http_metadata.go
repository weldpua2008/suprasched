package cluster
import (
    "github.com/weldpua2008/suprasched/communicator"
    "sync"
    "context"
    "fmt"
    "time"

)

type DescribeHttpCluster struct {
	ClusterDescriber
	mu           sync.RWMutex
    comm     communicator.Communicator
}

// NewDescribeEMR prepare struct communicator for EMR
func NewDescribeHttp(comm communicator.Communicator) *DescribeHttpCluster {
	return &DescribeHttpCluster{ comm: comm	}
}


func (d *DescribeHttpCluster) DescribeCluster(params map[string]interface{}) (string, error) {
	var ClusterId string
	var ctx context.Context
	var clusterCtx context.Context
	var cancel context.CancelFunc
	ttr := 30

	for _, k := range []string{"ClusterId", "clusterID", "ClusterID", "clusterId",
		"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"} {
		if _, ok := params[k]; ok {
			ClusterId = params[k].(string)
			break
		}
	}
	for _, k := range []string{"context", "ctx"} {
		if _, ok := params[k]; ok {
			if v, ok := params[k].(context.Context); ok {
				ctx = v
				break
			}
		}
	}
	if ctx == nil {
		ctx = context.Background()
	}
	clusterCtx, cancel = context.WithTimeout(ctx, time.Duration(ttr)*time.Second)
	defer cancel() // cancel when we are getting the kill signal or exit
    param:=make(map[string]interface{})
    d.mu.Lock()
	defer d.mu.Unlock()

    res, err:=d.comm.Fetch(clusterCtx, param)
    result := "UNKNOWN"
    for _, v := range res {
        if v == nil {
            continue
        }
        // if v, ok1 := v.(map[string]interface{}); ok1 {
                for _, k := range []string{"ClusterStatus", "Cluster_Status", "Status", "status"} {
                    if _, ok := v[k]; ok {


                        return v[k].(string), nil
                    }

                }

		// }
    }

	// status := cl.Cluster.Status.State
	// result := *status
	// switch *status {
	// case emr.ClusterStateStarting, emr.ClusterStateBootstrapping:
	// 	result = "STARTING"
	// case emr.ClusterStateRunning, emr.ClusterStateWaiting:
	// 	result = "RUNNING"
	// case emr.ClusterStateTerminated, emr.ClusterStateTerminatedWithErrors:
	// 	result = "TERMINATED"
	// }

	return result, fmt.Errorf("Can't find ClusterId: %s %v", ClusterId, err)
}
