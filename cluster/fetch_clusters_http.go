package cluster

import (
	"context"
	"fmt"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"

	"strconv"
	"sync"
	"time"
)

type FetchClustersHttp struct {
	ClustersFetcher
	mu   sync.RWMutex
	comm communicator.Communicator
}

// NewFetchEMR prepare struct communicator for EMR
func NewFetchClustersHttp() (*FetchClustersHttp, error) {

	if comm, err := communicator.GetSectionCommunicator(fmt.Sprintf("%s.fetch", config.CFG_PREFIX_CLUSTER)); err == nil {
		// log.Tracef("Getting GetSectionCommunicator %s", config.CFG_PREFIX_CLUSTER)

		return &FetchClustersHttp{comm: comm}, nil
	} else {
		return nil, fmt.Errorf("Can't initialize FetchClusters '%s': %v", config.CFG_PREFIX_CLUSTER, err)
	}
}

func (f *FetchClustersHttp) Fetch() ([]*model.Cluster, error) {
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

	res, err := f.comm.Fetch(fetchCtx, params)
	if err == nil {
		for _, v := range res {
			if v == nil {
				continue
			}
			var clusterId string
			var cl *model.Cluster

			for _, k := range []string{"ClusterId", "clusterID", "ClusterID", "clusterId",
				"clusterid", "JobFlowID", "JobFlowId", "JobflowID", "jobFlowId"} {
				if _, ok := v[k]; ok {

					clusterId = v[k].(string)
					break
				}
			}
			if len(clusterId) < 1 {
				continue
			}
			cl = model.NewCluster(clusterId)

			for _, k := range []string{"ClusterStatus", "clusterStatus", "Cluster_Status", "Status", "status"} {
				if _, ok := v[k]; ok {
					cl.Status = v[k].(string)
					break
				}
			}
			for _, k := range []string{"clusterPool", "ClusterPool", "Pool", "pool"} {
				if _, ok := v[k]; ok {
					cl.ClusterPool = v[k].(string)
					break
				}
			}
			for _, k := range []string{"clusterProfile", "ClusterProfile", "AWSProfile", "AWS_Profile", "AWS_PROFILE"} {
				if _, ok := v[k]; ok {
					cl.ClusterProfile = v[k].(string)
					break
				}
			}
			for _, k := range []string{"ClusterRegion", "clusterRegion", "AWSRegion", "AWS_Region", "AWS_REGION"} {
				if _, ok := v[k]; ok {
					cl.ClusterRegion = v[k].(string)
					break
				}
			}

			for _, k := range []string{"CreateAt", "createAt", "Created", "createDate", "CreateDate"} {
				if _, ok := v[k]; ok {
					switch t := v[k].(type) {
					case string:
						if i, err := strconv.Atoi(t); err == nil {
							cl.CreateAt = time.Unix(int64(i), 0)
							break
						}
					case int:
						cl.CreateAt = time.Unix(int64(t), 0)
						break
					case float64:
						cl.CreateAt = time.Unix(int64(int(t)), 0)
						break

					}
				}
			}
			// cl.StartAt = cl.CreateAt
			for _, k := range []string{"StartAt", "startAt", "StartDate", "startDate"} {
				if _, ok := v[k]; ok {
					switch t := v[k].(type) {
					case string:
						if i, err := strconv.Atoi(t); err == nil {
							cl.StartAt = time.Unix(int64(i), 0)
							break
						}
					case int:
						cl.StartAt = time.Unix(int64(t), 0)
						break
					case float64:
						cl.StartAt = time.Unix(int64(int(t)), 0)
						break

					}
				}
			}
			for _, k := range []string{"lastUpdated", "lastUpdated", "LastActivityAt", "lastActivityAt"} {
				if _, ok := v[k]; ok {
					switch t := v[k].(type) {
					case string:
						if i, err := strconv.Atoi(t); err == nil {
							cl.LastActivityAt = time.Unix(int64(i), 0)
							break
						}
					case int:
						cl.LastActivityAt = time.Unix(int64(t), 0)
						break
					case float64:
						cl.LastActivityAt = time.Unix(int64(int(t)), 0)
						break

					}
				}
			}
			for _, k := range []string{"jobs_info", "job_info", "job_ids"} {
				if value_of_slice, ok := v[k].([]interface{}); ok {
                    for _, elem := range value_of_slice {
                        if value_map, ok1:=elem.(map[string]interface{}); ok1{

                            j := model.NewEmptyJob()
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
                                log.Tracef("Can't add job %v %v", j.StoreKey(),j)

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
		return results, nil
	}
	return results, fmt.Errorf("Can't fetch more clusters: %v", err)
}
