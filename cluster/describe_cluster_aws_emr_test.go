package cluster

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/emr/emriface"
	"github.com/stretchr/testify/mock"
	"regexp"
	"strings"
	"testing"
)

type mockEMR struct {
	emriface.EMRAPI
	mock.Mock
}

func (m *mockEMR) DescribeClusterWithContext(ctx aws.Context, input *emr.DescribeClusterInput, opts ...request.Option) (*emr.DescribeClusterOutput, error) {

	return m.DescribeCluster(input)
}

// Mock using the cluster id of input to set the cluster state
// ClusterId = "j-RUNNING" will result in a cluster with the RUNNING state
func (m *mockEMR) DescribeCluster(input *emr.DescribeClusterInput) (*emr.DescribeClusterOutput, error) {
	if !strings.HasPrefix(*input.ClusterId, "j-") {
		return nil, fmt.Errorf("DescribeCluster failed for ClusterId %v ", *input.ClusterId)
	}
	var state string
	r, _ := regexp.Compile("j-([[:alpha:]]+)")
	res := r.FindStringSubmatch(fmt.Sprintf("%v", *input.ClusterId))
	if len(res) > 1 {
		state = res[1]
	}
	if state == "" {
		return nil, fmt.Errorf("DescribeCluster failed")
	}
	if state == "TERMINATED" {
		return &emr.DescribeClusterOutput{
			Cluster: &emr.Cluster{
				Status: &emr.ClusterStatus{
					State: aws.String(state),
					StateChangeReason: &emr.ClusterStateChangeReason{
						Code: aws.String("BOOTSTRAP_FAILURE"),
					},
				},
			},
		}, nil
	}
	return &emr.DescribeClusterOutput{
		Cluster: &emr.Cluster{
			Status: &emr.ClusterStatus{
				State: aws.String(state),
			},
		},
	}, nil

}

// DefaultGetEMR implements EMR Api wrapper for tests.
func MockGetEMR(sess *session.Session) emriface.EMRAPI {
	return new(mockEMR)
}

func TestEMRClusterStatus(t *testing.T) {
	d := &DescribeEMR{
		awsSessions: make(map[string]*session.Session),
		t:           "DescribeEMR",
		getEmrSvc:   MockGetEMR,
	}
	cases := []struct {
		clusterId string
		want      string
		err       error
	}{
		{
			clusterId: "j-RUNNING",
			want:      "RUNNING",
			err:       nil,
		},
	}
	for _, tc := range cases {
		params := map[string]interface{}{
			"ClusterId": tc.clusterId,
		}
		status, err := d.ClusterStatus(params)
		if err != tc.err {
			t.Errorf("want %v, got %v ", tc.err, err)
		}
		if status != tc.want {
			t.Errorf("want %v, got %v ", tc.want, status)

		}
	}
}
