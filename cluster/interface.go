package cluster

type ClusterDescriber interface {
	DescribeCluster(map[string]interface{}) (string, error)
}
