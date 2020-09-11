package communicator
import (
	// "bytes"
	// "context"
	// "encoding/json"
	// "fmt"
	// "github.com/spf13/viper"
	// "io/ioutil"
	// "net/http"
	// "net/http/httptest"
	"testing"
	// "time"
    config "github.com/weldpua2008/suprasched/config"

)

func TestFetch(t *testing.T) {
    config.LoadCfgForTests(t, "fixtures/fetch_http.yml")

	// // startTrace()
	// want := "{\"job_uid\":\"job-testing.(*common).Name-fm\",\"run_uid\":\"1\",\"extra_run_id\":\"1\",\"msg\":\"'S'\\n\"}"
	// var got string
	// CMD := "echo && exit 0"
	// responses := []ApiJobResponse{
	// 	{
	// 		JobId:       "job_id",
	// 		JobStatus:   "PENDING",
	// 		JobName:     "job_name",
	// 		RunUID:      "run_uid",
	// 		ExtraRunUID: "extra_run_id",
	// 		CMD:         CMD,
	// 		Parameters:  []string{},
	// 		CreateDate:  "createDate",
	// 		LastUpdated: "lastUpdated",
	// 		StopDate:    "stopDate",
	// 	},
	// 	{
	// 		JobId:       "job_id",
	// 		JobStatus:   "PENDING",
	// 		JobName:     "job_name",
	// 		RunUID:      "run_uid",
	// 		ExtraRunUID: "extra_run_id",
	// 		CMD:         CMD,
	// 		Parameters:  []string{},
	// 		CreateDate:  "createDate",
	// 		LastUpdated: "lastUpdated",
	// 		StopDate:    "stopDate",
	// 	},
	// }
    //
	// // notifyStdoutSent := make(chan bool)
	// srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //
	// 	var c ApiJobResponse
	// 	if len(responses) > 1 {
	// 		c, responses = responses[0], responses[1:]
	// 	} else if len(responses) == 1 {
	// 		c = responses[0]
	// 	}
	// 	c1 := make([]ApiJobResponse, 0)
	// 	c1 = append(c1, c)
	// 	js, err := json.Marshal(&c1)
	// 	if err != nil {
	// 		log.Tracef("Failed to marshal for '%v' due %v", c, err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	if _, errWrite := w.Write(js); errWrite != nil {
	// 		t.Errorf("Can't w.Write %v due %v\n", js, err)
	// 	}
	// 	b, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		t.Errorf("ReadAll %s", err)
	// 	}
	// 	got = string(fmt.Sprintf("%s", b))
	// 	// notifyStdoutSent <- true
	// }))
	// defer func() {
	// 	srv.Close()
	// 	model.FetchNewJobAPIURL = ""
	// 	// restoreLevel()
	// }()
	// viper.SetConfigType("yaml")
	// var yamlExample = []byte(`
    // logs:
    //   update:
    //     method: GET
    // jobs:
    //   get:
    //     url: "` + srv.URL + `"
    //     method: POST
    // `)
    //
	// _ = viper.ReadConfig(bytes.NewBuffer(yamlExample))
    //
	// model.FetchNewJobAPIURL = srv.URL
	// log.Trace(fmt.Sprintf("model.FetchNewJobAPIURL  %s", model.FetchNewJobAPIURL))
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel() // cancel when we are getting the kill signal or exit
	// jobs := make(chan *model.Job, 1)
    //
	// go func() {
	// 	if err := StartGenerateJobs(ctx, jobs, time.Duration(5)*time.Second); err != nil {
	// 		log.Infof("StartGenerateJobs failed %v", err)
	// 	}
	// }()
	// for job := range jobs {
	// 	if job.Status != model.JOB_STATUS_PENDING {
	// 		t.Errorf("Expected %s, got %s", model.JOB_STATUS_PENDING, job.Status)
	// 	}
    //
	// 	if job.CMD != CMD {
	// 		t.Errorf("want %s, got %v", want, got)
	// 	}
	// 	job.Status = model.JOB_STATUS_CANCELED
	// 	// stop loop
	// 	if len(responses) == 1 {
	// 		cancel()
	// 	}
    //
	// }
}