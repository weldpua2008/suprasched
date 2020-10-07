package handlers

import (
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"
	"testing"
	"time"
)

func TestEmptyCluster(t *testing.T) {
	startTestingHandler(EmptyCluster)
	defer stopTestingHandler()
	cases := []struct {
		want bool // true - If we expect start cleanup of empty cluster
		job  map[string]interface{}
	}{
		{
			want: true,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_SUCCESS,
				"ClusterId": "ClusterId1",
				"Id":        "Id1",
			},
		},
		{
			want: true,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_FAILED,
				"ClusterId": "ClusterId2",
				"Id":        "Id2",
			},
		},

		{
			want: false,
			job: map[string]interface{}{
				"Status":    model.JOB_STATUS_PENDING,
				"ClusterId": "ClusterId3",
				"Id":        "Id3",
			},
		},
	}
	for _, tc := range cases {

		job := model.NewJobFromMap(tc.job)
		cls := model.NewCluster(tc.job["ClusterId"].(string))
		if ok := cls.Add(job); !ok {
			t.Errorf("Can't add job %v", job)

		}
		cls.TimeOutDuration = time.Minute * 1

		timeOutDuration := cls.TimeOutDuration
		timeOutStartAt := cls.TimeOutStartAt
		timeoutAt := cls.TimeOutAt

		rec, exist := config.ClusterRegistry.Record(cls.StoreKey())
		if exist {
			t.Errorf("exist %v", rec)
		}

		if ok := config.ClusterRegistry.Add(cls); !ok {
			t.Errorf("Can't add %v", cls)
		}

		if err := emitTestingData(cls.EventMetadata()); err != nil {
			t.Errorf("%v", err)
		}
		if timeOutDuration != cls.TimeOutDuration {
			t.Errorf("want %v got %v", timeOutDuration, cls.TimeOutDuration)
		}

		if tc.want {

			if timeOutStartAt == cls.TimeOutStartAt {
				t.Errorf("want %v <> got %v", timeOutStartAt, cls.TimeOutStartAt)
			}
			if timeoutAt == cls.TimeOutAt {
				t.Errorf("want %v <> got %v", timeoutAt, cls.TimeOutAt)
			}
		} else {
            timeoutAt =cls.TimeOutStartAt.Add(cls.TimeOutDuration)
			// if timeOutStartAt != cls.TimeOutStartAt {
			// 	t.Errorf("want %v == got %v", timeOutStartAt, cls.TimeOutStartAt)
			// }
			if timeoutAt != cls.TimeOutAt {
				t.Errorf("want %v == got %v", timeoutAt, cls.TimeOutAt)
			}
		}
	}
}
