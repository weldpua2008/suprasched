package job

import (
	"context"
	"fmt"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"

	"sync"
	"time"
)

func init() {
	Constructors[ConstructorsJobsFetcherRest] = TypeSpec{
		instance:    NewFetchJobsRest,
		constructor: NewFetchJobsDefault,
		Summary: `
FetchJobsDefault is the default implementation of JobsFetcher and is
used by Default.`,
		Description: `
It supports the following params:
- ` + "`ClusterId`" + ` Cluster Identificator
- ` + "`ClusterPool`" + ` To differentiate jobs by Pools
- ` + "`ClusterProfile`" + ` To differentiate jobs by Accounts.`,
	}
}

type FetchJobsDefault struct {
	JobsFetcher
	mu    sync.RWMutex
	comm  communicator.Communicator
	comms []communicator.Communicator
	t     string
}

// NewFetchJobsDefault prepare struct FetchJobsDefault
func NewFetchJobsDefault(section string) (JobsFetcher, error) {
	comms, err := communicator.GetCommunicatorsFromSection(fmt.Sprintf("%v.%v", section, config.CFG_PREFIX_JOBS_FETCHER))
	if err == nil {
		return &FetchJobsDefault{comms: comms, t: "FetchJobsDefault"}, nil
	} else {
		comm, err := communicator.GetSectionCommunicator(section)
		if err == nil {
			comms := make([]communicator.Communicator, 0)
			comms = append(comms, comm)
			return &FetchJobsDefault{comm: comm, comms: comms, t: "FetchJobsDefault"}, nil

		}
	}
	return nil, fmt.Errorf("Can't initialize FetchJobs '%s': %v", config.CFG_PREFIX_JOBS, err)

}

// NewFetchJobsDefault prepare struct FetchJobsDefault
func NewFetchJobsRest() JobsFetcher {

	return &FetchJobsDefault{t: "NewFetchJobsRest"}

}

func (f *FetchJobsDefault) Fetch() ([]*model.Job, error) {

	var results []*model.Job

	//var ctx context.Context
	var fetchCtx context.Context
	var cancel context.CancelFunc
	//if ctx == nil {
	var ctx = context.Background()
	//}
	ttr := 30
	fetchCtx, cancel = context.WithTimeout(ctx, time.Duration(ttr)*time.Second)
	defer cancel() // cancel when we are getting the kill signal or exit
	params := make(map[string]interface{})
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, comm := range f.comms {
		res, err := comm.Fetch(fetchCtx, params)
		if err != nil {
			return nil, fmt.Errorf("Can't fetch more jobs: %v, %v", err, comm)
		}
		for _, v := range res {
			if v == nil {
				continue
			}

			j := model.NewJobFromMap(v)

			if len(j.Id) < 1 {
				continue
			}

			results = append(results, j)
		}
	}
	return results, nil

}
