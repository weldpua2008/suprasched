package core

// Cluster is a worker cluster in Kubernetes
// The name of the cluster according to etcd is in ObjectMeta.Name.
type Cluster struct {
	TypeMeta
	// +optional
	ObjectMeta

	// Spec defines the behavior of a cluster.
	// +optional
	Spec ClusterSpec

	// Status describes the current status of a Cluster
	// +optional
	Status ClusterStatus
}

// ClusterSpec describes the attributes that a cluster is created with.
type ClusterSpec struct {
	// ID of the cluster assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>
	// +optional
	ProviderID string `json:"providerID,omitempty" protobuf:"bytes,1,opt,name=providerID"`

	// Unschedulable controls node schedulability of new pods. By default node is schedulable.
	// +optional
	Unschedulable bool `json:"unschedulable,omitempty" protobuf:"varint,2,opt,name=unschedulable"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]interface{}

// ClusterStatusType defines the condition of pod.
type ClusterStatusType string

// These are valid conditions of cluster.
const (
	ClusterStatusStarting             ClusterStatusType = "STARTING"
	ClusterStatusBootstraping         ClusterStatusType = "BOOTSTRAPING"
	ClusterStatusRunning              ClusterStatusType = "RUNNING"
	ClusterStatusWaiting              ClusterStatusType = "WAITING"
	ClusterStatusTerminating          ClusterStatusType = "TERMINATING"
	ClusterStatusTerminated           ClusterStatusType = "TERMINATED"
	ClusterStatusTerminatedWithErrors ClusterStatusType = "TERMINATED_WITH_ERRORS"
	ClusterStatusPending              ClusterStatusType = "PENDING"
	ClusterStatusInProgress           ClusterStatusType = "RUNNING"
	ClusterStatusError                ClusterStatusType = "ERROR"
	ClusterStatusCanceled             ClusterStatusType = "CANCELED"
	ClusterStatusTimeout              ClusterStatusType = "TIMEOUT"
)

// ClusterStatus is information about the current status of a cluster.
type ClusterStatus struct {
	// Capacity represents the total resources of a cluster.
	// +optional
	Capacity ResourceList
	// Allocatable represents the resources of a cluster that are available for scheduling.
	// +optional
	Allocatable ResourceList
	// Status is the current lifecycle phase of the cluster.
	// +optional
	Status ClusterStatusType
}

// NewCluster returns a new cluster
func NewCluster(name string, ns Namespace, spec ClusterSpec, status ClusterStatus, uid UID) Cluster {
	return Cluster{
		Status: status,
		Spec:   spec,
		ObjectMeta: ObjectMeta{
			Name:      name,
			Namespace: ns,
			UID:       uid,
		},
		TypeMeta: TypeMeta{
			Kind:       "cluster",
			APIVersion: LatestApi,
		},
	}
}
