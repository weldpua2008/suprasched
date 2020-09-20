package handlers

import (
    "fmt"
    "time"
    "context"
	"github.com/mustafaturan/bus"
    config "github.com/weldpua2008/suprasched/config"
    communicator "github.com/weldpua2008/suprasched/communicator"


)

// ClusterTermination handler
// 1. Send Cluster termination (external)
// 2. Find new cluster / Create a new cluster
// 3. Lock Cluster & sub jobs
// 4. Reassign Jobs (ext)
// 5. Unlock Cluster & sub jobs
func ClusterTermination(e *bus.Event) {
    if cl,ok:= e.Data.(map[string]string);ok {
        if storeKey, ok := cl["StoreKey"];ok {
            if rec, ok := config.ClusterRegistry.Record(storeKey); ok {


                section:= fmt.Sprintf("%v.%v.%v",
                    config.CFG_PREFIX_CLUSTER,config.CFG_PREFIX_UPDATE, rec.ClusterType )
                if comms, err:=  communicator.GetCommunicatorsFromSection(section);err == nil {
                    fetchCtx, cancel := context.WithTimeout(context.Background(), time.Duration(15)*time.Second)
                    params:=rec.GetParams()
                    from := rec.GetParamsMapString()

                	defer cancel() // cancel when we are getting the kill signal or exit
                    for _,comm:= range comms {
                        log.Tracef("Comm %v", comm)
                        var cfg_params map[string]interface{}
                        cfg_params = config.ConvertMapStringToInterface(
                            config.GetStringMapStringTemplatedFromMap(section, config.CFG_COMMUNICATOR_PARAMS_KEY, from))

                        // if err := comm.Configure(cfg_params); err != nil {
                        //     log.Tracef("Can't Configure communicator with %v because %v ", cfg_params, err)
                        // }

                        res, err := comm.Fetch(fetchCtx, params)
                		if err != nil {
                            log.Tracef("Can't change status %v for %v because %v ", res, rec.ClusterId, err)
                		}else {
                            log.Tracef("status changed cfg_params %v %v for %v because %v ", cfg_params, res, rec.ClusterId, err)
                        }

                    }
                }else {
                    log.Tracef("Can't get communicators because %v for %v Event for %s: %+v", err, rec, e.Topic, e)

                }

                log.Tracef("Terminated for %v Event for %s: %+v", rec, e.Topic, e)
            }else {
                log.Tracef("No such StoreKey in ClusterRegistry e.Data %v forEvent for %s: %+v", e.Data, e.Topic, e)

            }
        }else {
            log.Tracef("No StoreKey in e.Data %v for Event for %s: %+v", e.Data, e.Topic, e)

        }
    }else {
        log.Tracef("Can't read e.Data %v for Event for %s: %+v", e.Data, e.Topic, e)
    }
}
