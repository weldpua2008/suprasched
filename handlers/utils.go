package handlers

import (
	"context"
	"fmt"
	"github.com/mustafaturan/bus"
	communicator "github.com/weldpua2008/suprasched/communicator"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"

	"time"
)

func eventDataStringMapString(e *bus.Event) (res map[string]string) {
	if d, ok := e.Data.(map[string]string); ok {
		res = d
	}
	return res
}

func eventGetCLuster(e *bus.Event) (*model.Cluster, error) {
	res := eventDataStringMapString(e)
	if storeKey, ok := res["StoreKey"]; ok {
		if rec, ok := config.ClusterRegistry.Record(storeKey); ok {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("%w", ErrNoClusterFound)
}

func eventCLusterRunComms(e *bus.Event) error {
	if rec, err := eventGetCLuster(e); err == nil {
		section := fmt.Sprintf("%v.%v.%v",
			config.CFG_PREFIX_CLUSTER, config.CFG_PREFIX_UPDATE, rec.ClusterType)
		if comms, err := communicator.GetCommunicatorsFromSection(section); err == nil {
			fetchCtx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
			params := rec.GetParams()
			// from := rec.GetParamsMapString()

			defer cancel() // cancel when we are getting the kill signal or exit
			for _, comm := range comms {
				// var cfg_params map[string]interface{}
				// cfg_params = config.ConvertMapStringToInterface(
				// 	config.GetStringMapStringTemplatedFromMap(section, config.CFG_COMMUNICATOR_PARAMS_KEY, from))
				res, err := comm.Fetch(fetchCtx, params)
				if err != nil {
					return fmt.Errorf("%w for %v got %v error %v", ErrFailedUpdateStatus, rec.ClusterId, res, err)
				}
				// else {
				//         log.Tracef("Performed cfg_params %v %v for %v because %v ", cfg_params, res, rec.ClusterId, err)
				//     }

			}
		} else {
			return fmt.Errorf("%w for %v got error %v", communicator.ErrNoSuitableCommunicator, rec.ClusterId, err)
		}
	}

	return nil
}

func eventJobRunComms(j *model.Job, eData map[string]string) error {
	section := fmt.Sprintf("%v.%v",
		config.CFG_PREFIX_JOBS, config.CFG_PREFIX_UPDATE)
	if comms, err := communicator.GetCommunicatorsFromSection(section); err == nil {
		fetchCtx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
		params := j.GetParams()

		defer cancel() // cancel when we are getting the kill signal or exit
		for _, comm := range comms {
			res, err := comm.Fetch(fetchCtx, params)
			if err != nil {
				return fmt.Errorf("%w for %v got %v error %v", ErrFailedUpdateStatus, j.Id, res, err)
			}
		}

	}
	return nil
}

// emitTestingData sends one event.
func emitTestingData(in map[string]string) error {
	topic := "testing"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are getting the kill signal or exit
	_, err := config.Bus.Emit(ctx, topic, in)
	return err

}
