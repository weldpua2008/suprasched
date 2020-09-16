package cluster

import (
	"context"
	"fmt"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	utils "github.com/weldpua2008/suprasched/utils"

	"sync"
	"time"
)

func init() {
	FetcherConstructors[ConstructorsFetcherTypeRest] = FetcherTypeSpec{
		instance:    NewFetchClusterHttp,
		constructor: NewFetchClustersDefault,
		Summary: `
FetchClustersDefault is the default implementation of ClustersFetcher and is
used by Default.`,
		Description: `
It supports the following params:
- ` + "`ClusterId`" + ` Cluster Identificator
- ` + "`ClusterPool`" + ` To differentiate clusters by Pools
- ` + "`ClusterProfile`" + ` To differentiate clusters by Accounts.`,
	}
}

type FetchClustersDefault struct {
	ClustersFetcher
	mu    sync.RWMutex
	comm  communicator.Communicator
	comms []communicator.Communicator
	t     string
}

// NewFetchClustersDefault prepare struct FetchClustersDefault
func NewFetchClustersDefault(section string) (ClustersFetcher, error) {
	comms, err := communicator.GetCommunicatorsFromSection(section)
	if err == nil {
		return &FetchClustersDefault{comms: comms, t: "FetchClustersDefault"}, nil
	} else {
		comm, err := communicator.GetSectionCommunicator(section)
		if err == nil {
			comms := make([]communicator.Communicator, 0)
			comms = append(comms, comm)
			return &FetchClustersDefault{comm: comm, comms: comms, t: "FetchClustersDefault"}, nil

		}
	}
	return nil, fmt.Errorf("Can't initialize FetchClusters '%s': %v", config.CFG_PREFIX_CLUSTER, err)

}

// NewFetchClustersDefault prepare struct FetchClustersDefault
func NewFetchClusterHttp() ClustersFetcher {

	return &FetchClustersDefault{}

}

func (f *FetchClustersDefault) Fetch() ([]*model.Cluster, error) {
	var results []*model.Cluster

	var ctx context.Context
	var fetchCtx context.Context
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ttr := 30
	fetchCtx, cancel = context.WithTimeout(ctx, time.Duration(ttr)*time.Second)
	defer cancel() // cancel when we are getting the kill signal or exit
	params := make(map[string]interface{})
	// f.mu.Lock()
	// defer f.mu.Unlock()
	for _, comm := range f.comms {
		res, err := comm.Fetch(fetchCtx, params)
		if err != nil {
			return nil, fmt.Errorf("Can't fetch more clusters: %v", err)
		}
		for _, v := range res {
			if v == nil {
				continue
			}
			var cl *model.Cluster
			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"ClusterId", "clusterID", "ClusterID", "clusterId",
				"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"}); ok {
				cl = model.NewCluster(found_val)
			} else {
				continue
			}

			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"ClusterStatus", "clusterStatus", "Cluster_Status", "Status", "status"}); ok {
				cl.Status = found_val
			}
			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"clusterPool", "ClusterPool", "Pool", "pool"}); ok {
				cl.ClusterPool = found_val
			}

			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"clusterProfile", "ClusterProfile", "AWSProfile", "AWS_Profile", "AWS_PROFILE"}); ok {
				cl.ClusterProfile = found_val
			}
			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"ClusterRegion", "clusterRegion", "AWSRegion", "AWS_Region", "AWS_REGION"}); ok {
				cl.ClusterRegion = found_val
			}

			if found_val, ok := utils.GetFirstStringFromMap(v, []string{"ClusterType", "clusterType"}); ok {
				cl.ClusterType = found_val
			}

			if found_val, ok := utils.GetFirstTimeFromMap(v, []string{"CreateAt", "createAt", "Created", "createDate", "CreateDate"}); ok {
				cl.CreateAt = found_val
			}
			if found_val, ok := utils.GetFirstTimeFromMap(v, []string{"StartAt", "startAt", "StartDate", "startDate"}); ok {
				cl.StartAt = found_val
			}
			if found_val, ok := utils.GetFirstTimeFromMap(v, []string{"lastUpdated", "lastUpdated", "LastActivityAt", "lastActivityAt"}); ok {
				cl.LastActivityAt = found_val
			}

			for _, k := range []string{"jobs_info", "job_info", "job_ids"} {
				if value_of_slice, ok := v[k].([]interface{}); ok {
					for _, elem := range value_of_slice {
						if value_map, ok1 := elem.(map[string]interface{}); ok1 {

							j := model.NewEmptyJob()
							if found_val, ok := utils.GetFirstTimeFromMap(v, []string{"StartAt", "startAt", "StartDate", "startDate"}); ok {
								j.StartAt = found_val
							}
							for _, sub_key := range []string{"JobStatus", "jobStatus", "Job_Status", "Status", "status"} {
								if _, ok := value_map[sub_key]; ok {
									j.Status = value_map[sub_key].(string)
									break
								}
							}

							for _, sub_key := range []string{"JobId", "jobId", "Job_ID", "Job_Id", "job_Id", "job_id"} {
								if _, ok := value_map[sub_key]; ok {
									j.Id = value_map[sub_key].(string)
									break
								}
							}

							for _, sub_key := range []string{"JobRunId", "jobRunId", "Job_RUN_ID", "Job_Run_Id", "job_run_id", "run_id",
								"run_uid", "RunId", "RunUID"} {
								if _, ok := value_map[sub_key]; ok {
									j.RunUID = value_map[sub_key].(string)
									break
								}
							}
							for _, sub_key := range []string{"JobExtraRunId", "jobExtraRunId", "JOB_EXTRA_RUN_ID", "Job_Extra_Run_Id", "job_extra_run_id", "extra_run_id",
								"job_extra_run_uid", "extra_run_uid"} {
								if _, ok := value_map[sub_key]; ok {
									j.ExtraRunUID = value_map[sub_key].(string)
									break
								}
							}
							j.ClusterId = cl.ClusterId
							if len(j.Id) < 1 {
								continue
							}
							var topic string
							if config.JobsRegistry.Add(j) {
								topic = config.TOPIC_JOB_CREATED
								log.Tracef("Job %v added ", j.StoreKey())
							}
							if job_on_cluster, ok := config.JobsRegistry.Record(j.StoreKey()); ok {
								cl.Add(job_on_cluster)
								if len(topic) > 0 {
									_, err := config.Bus.Emit(ctx, topic, j.EventMetadata())
									if err != nil {
										log.Tracef("%v", err)
									}

								}
							} else {
								log.Tracef("Can't add job %v %v", j.StoreKey(), j)

							}

							// if len(topic) > 0 {
							//     _, err := config.Bus.Emit(ctx, topic, cls.EventMetadata())
							//     if err != nil {
							//         log.Tracef("%v", err)
							//     }
							//
							// }
							// if !config.ClusterRegistry.Add(cls) {
							//     if rec, exist := config.ClusterRegistry.Record(cls.StoreKey()); exist {
							//         if rec.UseExternaleStatus(cls) {
							//             topic = strings.ToLower(fmt.Sprintf("cluster.%v", cls.Status))
							//         }
							//
							//     }
							// } else {
							//     topic = config.TOPIC_CLUSTER_CREATED
							// }
							// if len(topic) > 0 {
							//     _, err := config.Bus.Emit(ctx, topic, cls.EventMetadata())
							//     if err != nil {
							//         log.Tracef("%v", err)
							//     }
							//
							// }
						}

					}

				}
			}
			results = append(results, cl)
		}
	}
	return results, nil

}
