package cluster

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr/emriface"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	utils "github.com/weldpua2008/suprasched/utils"

	"sync"
	"time"
)

func init() {
	TerminatorConstructors[ConstructorsTerminaterTypeEMR] = TerminatorTypeSpec{
		instance:    NewTerminateClusterEMR,
		constructor: NewTerminateClusterEMRBySection,
		Summary: `
TerminateEMR is an implementation of ClustersTerminator for Amazon EMR clusters.`,
		Description: `
It supports the following params:
- ` + "`ClusterId`" + ` Cluster Identificator
- ` + "`ClusterPool`" + ` To differentiate clusters by Pools
- ` + "`ClusterProfile`" + ` To differentiate clusters by Accounts.`,
	}
}

type TerminateClusterEMR struct {
	ClustersTerminator
	section string
	mu      sync.RWMutex
	comm    communicator.Communicator
	comms   []communicator.Communicator
	t       string
	getemr  func(*session.Session) emriface.EMRAPI
}

// NewTerminateEMR prepare struct communicator for EMR
func NewTerminateClusterEMR() ClustersTerminator {
	return &TerminateClusterEMR{
		getemr: DefaultGetEMR,
	}
}

// NewTerminateClustersDefault prepare struct TerminateClustersDefault
func NewTerminateClusterEMRBySection(section string) (ClustersTerminator, error) {
	comms, err := communicator.GetCommunicatorsFromSection(section)
	if err == nil {
		return &TerminateClusterEMR{comms: comms, t: "TerminateClusterEMR", section: section}, nil
	} else {
		comm, err := communicator.GetSectionCommunicator(section)
		if err == nil {
			comms := make([]communicator.Communicator, 0)
			comms = append(comms, comm)
			return &TerminateClusterEMR{
				comm:    comm,
				comms:   comms,
				t:       "TerminateClusterEMR",
				section: section,
				getemr:  DefaultGetEMR,
			}, nil

		}
	}
	return nil, fmt.Errorf("Can't initialize TerminateClusterEMR '%s': %v", config.CFG_PREFIX_CLUSTER, err)
}

// SupportedClusters returns all supported in Cluster Registry defined by configuration(e.g. type).
// For example support on-demand and HTTP types in config:
//     cluster:
//         terminate:
//             supported:
//                 - "EMR"
func (d *TerminateClusterEMR) SupportedClusters() []*model.Cluster {
	def := []string{ConstructorsFetcherTypeRest}
	cluster_types := config.GetGetStringSliceDefault(fmt.Sprintf("%v.%v", d.section, config.CFG_PREFIX_CLUSTER_SUPPORTED_TYPES), def)
	return config.ClusterRegistry.FilterFree(cluster_types)
}

// ClusterStatus by the Cluster Id from HTTP rest API.
func (d *TerminateClusterEMR) Terminate(params map[string]interface{}) error {
	var ClusterId string
	var ctx context.Context
	var clusterCtx context.Context
	var cancel context.CancelFunc
	ttr := 30

	for _, k := range []string{"ClusterId", "Clusterid", "clusterID", "ClusterID", "clusterId",
		"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"} {
		if _, ok := params[k]; ok {
			ClusterId = params[k].(string)
			break
		}
	}
	if len(ClusterId) < 1 {
		return ErrEmptyClusterId
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
	// param := make(map[string]interface{})
	param := config.ConvertMapStringToInterface(
		config.GetStringMapStringTemplated(d.section, config.CFG_PREFIX_COMMUNICATORS))
	for k, v := range params {
		if k == "context" || k == "ctx" {
			continue
		}

		param[k] = v
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	utils.EmrClusterTerminate(params, d.getemr)
	log.Tracef("Terminate Cluster %v", ClusterId)

	for _, comm := range d.comms {
		comm.Configure(params)
		res, err := comm.Fetch(clusterCtx, param)
		if err != nil {
			log.Tracef("Can't Terminate %v %v %v", ClusterId, err, res)
			continue
		}

	}
	return nil
}
