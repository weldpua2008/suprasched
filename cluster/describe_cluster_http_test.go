package cluster

import (
	config "github.com/weldpua2008/suprasched/config"
	model "github.com/weldpua2008/suprasched/model"

	"fmt"
	"go.uber.org/goleak"
	"testing"
)

type Response struct {
	ClusterId     string `json:"ClusterId"`
	ClusterStatus string `json:"ClusterStatus"`
}

var (
	globalGot string
	responses []Response
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
	globalGot = in
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestSupportedClusters(t *testing.T) {
	C, tmpC := config.LoadCfgForTests(t, "fixtures/describe_cluster_http.yml")
	config.C = C
	defer func() {
		config.C = tmpC
		config.RefreshRegistries()

	}()
	cases := []struct {
		section          string
		unsupported_type []string
		supported_type   []string
	}{

		{
			// defined
			section:          fmt.Sprintf("%v.%v", config.CFG_PREFIX_CLUSTER, config.CFG_PREFIX_DESCRIBERS),
			unsupported_type: []string{"EMR", "emr", "unsupported_type", "unsupported"},
			supported_type:   []string{"HTTP", "on-demand"},
		},
		{
			// default
			section:          config.CFG_PREFIX_CLUSTER,
			unsupported_type: []string{"EMR", "emr", "unsupported_type", "unsupported", "on-demand"},
			supported_type:   []string{"HTTP"},
		},
	}

	for _, tc := range cases {

		config.RefreshRegistries()
		descr, err := NewDescribeClusterHttpBySection(tc.section)
		if err != nil {
			t.Errorf("want %v, got %v", nil, err)
		}

		for idx, unsupported_type := range tc.unsupported_type {
			cls := model.NewCluster(fmt.Sprintf("Cluster-unsupported_type-%v", idx))
			cls.ClusterType = unsupported_type

			if ok := config.ClusterRegistry.Add(cls); !ok {
				t.Errorf("Can't add Cluster %v", cls)
			}
		}
		if len(descr.SupportedClusters()) > 0 {
			t.Errorf("want %v, got %v %v", 0, len(descr.SupportedClusters()), descr.SupportedClusters())
		}
		for idx, supported_type := range tc.supported_type {
			cls := model.NewCluster(fmt.Sprintf("Cluster-supported_type-%v", idx))
			cls.ClusterType = supported_type

			if ok := config.ClusterRegistry.Add(cls); !ok {
				t.Errorf("Can't add Cluster %v", cls)
			}
		}
		if len(descr.SupportedClusters()) < len(tc.supported_type) {
			t.Errorf("want %v, got %v %v", len(tc.supported_type), len(descr.SupportedClusters()), descr.SupportedClusters())
		}
	}

}

func TestClusterStatus(t *testing.T) {
	srv := config.NewTestServer(t, in, out)

	C, tmpC := config.LoadCfgForTests(t, "fixtures/describe_cluster_http.yml")
	config.C = C
	config.C.URL = srv.URL
	defer func() {
		globalGot = ""
		srv.Close()
		responses = []Response{}

		config.C = tmpC
		config.RefreshRegistries()

	}()

	cases := []struct {
		section          string
		unsupported_type []string
		supported_type   []string
		responses        []Response
	}{
		{
			responses: []Response{
				{
					ClusterId:     "id-1",
					ClusterStatus: "RUNNING",
				},
			},
			section:          fmt.Sprintf("%v.%v", config.CFG_PREFIX_CLUSTER, config.CFG_PREFIX_DESCRIBERS),
			unsupported_type: []string{"EMR", "emr", "unsupported_type", "unsupported"},
			supported_type:   []string{"HTTP", "on-demand"},
		},
		{
			// non-exists
			section:          "empty",
			unsupported_type: []string{"EMR", "emr", "unsupported_type", "unsupported", "on-demand"},
			supported_type:   []string{}, // default value
		},
	}

	for _, tc := range cases {
		responses = tc.responses
		descr, err := NewDescribeClusterHttpBySection(tc.section)
		if err != nil {
			t.Errorf("want %v, got %v", nil, err)
		}
		for _, supported_type := range tc.supported_type {

			for _, resp := range tc.responses {
				config.RefreshRegistries()

				cls := model.NewCluster(resp.ClusterId)
				cls.ClusterType = supported_type
				if ok := config.ClusterRegistry.Add(cls); !ok {
					t.Errorf("Can't add Cluster %v", cls)
				}
				if len(descr.SupportedClusters()) < len(tc.responses) {
					t.Errorf("want %v, got %v %v", len(tc.responses), len(descr.SupportedClusters()), descr.SupportedClusters())
				}
				params := map[string]interface{}{
					"ClusterId": resp.ClusterId,
				}
				status, err := descr.ClusterStatus(params)
				if status != resp.ClusterStatus {
					t.Errorf("want %v, got %v %v", resp.ClusterStatus, status, resp.ClusterId)

				}
				if err != nil {
					t.Errorf("want %v, got %v %v", nil, err, resp.ClusterId)

				}
			}
		}
	}

}
