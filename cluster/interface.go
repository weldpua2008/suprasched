package cluster

import (
	model "github.com/weldpua2008/suprasched/model"
)

// ClusterDescriber fetch metadata for specific cluster
type ClusterDescriber interface {
	DescribeCluster(map[string]interface{}) (string, error)
}

// ClustersFetcher fetch Cluster list
type ClustersFetcher interface {
	Fetch() ([]*model.Cluster, error)
}
