package cluster

import (
	model "github.com/weldpua2008/suprasched/model"
)

// ClusterDescriber fetch metadata for specific cluster
type ClustersDescriber interface {
	ClusterStatus(map[string]interface{}) (string, error)
	SupportedClusters() []*model.Cluster
}

// ClustersFetcher fetch Cluster's list
type ClustersFetcher interface {
	Fetch() ([]*model.Cluster, error)
}

// ClustersTerminator Terminates Cluster
type ClustersTerminator interface {
	Terminate(map[string]interface{}) error
	SupportedClusters() []*model.Cluster
}
