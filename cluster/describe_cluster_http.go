package cluster

import (
	"context"
	"fmt"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"

	"sync"
	"time"
)

func init() {
	DescriberConstructors[ConstructorsDescriberTypeRest] = DescriberTypeSpec{
		instance:    NewDescribeClusterHttp,
		constructor: NewDescribeClusterHttpBySection,
		Summary: `
DescribeEMR is an implementation of ClustersDescriber for Amazon EMR clusters.`,
		Description: `
It supports the following params:
- ` + "`ClusterId`" + ` Cluster Identificator
- ` + "`ClusterPool`" + ` To differentiate clusters by Pools
- ` + "`ClusterProfile`" + ` To differentiate clusters by Accounts.`,
	}
}

type DescribeClusterHttp struct {
	ClustersDescriber
	mu    sync.RWMutex
	comm  communicator.Communicator
	comms []communicator.Communicator
	t     string
}

// NewDescribeEMR prepare struct communicator for EMR
func NewDescribeClusterHttp() ClustersDescriber {
	return &DescribeClusterHttp{}
}

// NewDescribeClustersDefault prepare struct DescribeClustersDefault
func NewDescribeClusterHttpBySection(section string) (ClustersDescriber, error) {
	comms, err := communicator.GetCommunicatorsFromSection(section)
	if err == nil {
		return &DescribeClusterHttp{comms: comms, t: "DescribeClusterHttp"}, nil
	} else {
		comm, err := communicator.GetSectionCommunicator(section)
		if err == nil {
			comms := make([]communicator.Communicator, 0)
			comms = append(comms, comm)
			return &DescribeClusterHttp{comm: comm, comms: comms, t: "DescribeClusterHttp"}, nil

		}
	}
	return nil, fmt.Errorf("Can't initialize DescribeClusterHttp '%s': %v", config.CFG_PREFIX_CLUSTER, err)

}

func (d *DescribeClusterHttp) ClusterStatus(params map[string]interface{}) (string, error) {
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
	param := make(map[string]interface{})
	d.mu.Lock()
	defer d.mu.Unlock()

	res, err := d.comm.Fetch(clusterCtx, param)
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
