package handlers

import (
	// "github.com/sirupsen/logrus"
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"

	"strings"
	"testing"
)

type Response struct {
	Number int    `json:"number"`
	Str    string `json:"str"`
}

var (
	globalGotCluster string
	globalGotJob     string
	responses        []Response
)

func in() interface{} {
	var c Response
	if len(responses) > 1 {
		c, responses = responses[0], responses[1:]
	} else if len(responses) == 1 {
		c = responses[0]
	}
	c1 := make([]Response, 0)
	c1 = append(c1, c)
	return c1
}

func out(in string) {
	if strings.Contains(in, "typejob") {
		globalGotJob = in
	} else {
		globalGotCluster = in
	}

}

func TestClusterTermination(t *testing.T) {
	// logrus.SetLevel(logrus.TraceLevel)

	C, tmpC := config.LoadCfgForTests(t, "fixtures/cluster_termination.yml")
	config.C = C

	srv := config.NewTestServer(t, in, out)

	defer func() {
		globalGotCluster = ""
		globalGotJob = ""
		srv.Close()

		config.C = tmpC
	}()

	startTestingHandler(ClusterTermination)
	defer stopTestingHandler()
	cases := []struct {
		want   bool // true - If we expect termination process
		status string
		job    map[string]interface{}
	}{
		{
			want:   true,
			status: model.CLUSTER_STATUS_TERMINATING,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_PENDING,
				"ClusterId": "ClusterId1",
				"Id":        "Id1",
			},
		},
		{
			want:   true,
			status: model.CLUSTER_STATUS_TERMINATING,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_PENDING,
				"ClusterId": "ClusterId2",
				"Id":        "Id2",
			},
		},
		{
			want:   false,
			status: model.CLUSTER_STATUS_RUNNING,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_PENDING,
				"ClusterId": "ClusterId3",
				"Id":        "Id3",
			},
		},
	}
	for _, tc := range cases {
		globalGotCluster = ""
		globalGotJob = ""
		responses = []Response{
			{
				Number: 1,
				Str:    "Str",
			},
			{
				Number: 2,
				Str:    "Str1",
			},
		}

		job := model.NewJobFromMap(tc.job)
		cls := model.NewCluster(tc.job["ClusterId"].(string))
		cls.ClusterType = model.CLUSTER_TYPE_ON_DEMAND
		if ok := cls.Add(job); !ok {
			t.Errorf("Can't add job %v", job)

		}
		rec_job, exist := config.JobsRegistry.Record(job.StoreKey())
		if exist {
			t.Errorf("Job exist %v", rec_job)
		}
		job_status := job.GetStatus()
		if ok := config.JobsRegistry.Add(job); !ok {
			t.Errorf("Can't add  job %v", job)
		}

		rec, exist := config.ClusterRegistry.Record(cls.StoreKey())
		if exist {
			t.Errorf("Cluster exist %v", rec)
		}

		cls.Status = model.CLUSTER_STATUS_STARTING
		if ok := config.ClusterRegistry.Add(cls); !ok {
			t.Errorf("Can't add %v", cls)
		}
		cls.Status = tc.status
		config.C.URL = srv.URL
		job.ExtraSendParams = map[string]string{
			"url": srv.URL,
		}

		if err := emitTestingData(cls.EventMetadata()); err != nil {
			t.Errorf("Failed to send event %v", err)
		}

		if tc.want {
			if !model.IsTerminalStatus(tc.job["Status"].(string)) {
				if !strings.Contains(globalGotJob, job.GetStatus()) {
					t.Errorf("want %v = got %v [%v]", job.GetStatus(), globalGotJob, job.GetStatus())
				}

				if !strings.Contains(globalGotCluster, cls.Status) {
					t.Errorf("want %v = got %v [%v]", job.Status, globalGotCluster, cls.Status)
				}
				if globalGotCluster == "" {
					t.Errorf("want %v != got %v", "", globalGotCluster)
				}

			}
		} else {
			if globalGotCluster != "" {
				t.Errorf("want %v = got %v", "", globalGotCluster)
			}
			if globalGotJob != "" {
				t.Errorf("want %v = got %v", "", globalGotCluster)
			}

			if job.GetStatus() != job_status {
				t.Errorf("want %v = got %v", job_status, job.GetStatus())
			}
		}
	}
}
