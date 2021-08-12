package core

// UID is a type that holds unique ID values, including UUIDs.
type UID string

// Namespace is a type that holds namespaces.
type Namespace string

// TypeMeta describes an individual object  with strings representing the type of the object
// and its API schema version.
// Structures that are versioned or persisted should inline TypeMeta.
//
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Cannot be updated.
	// In CamelCase.
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// that users must create.
type ObjectMeta struct {
	// Name must be unique within a type. Is required when creating resources.
	// Name is primarily intended for creation idempotence and configuration definition.
	// Cannot be updated.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Namespace defines the space within which each name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	//
	// Cannot be updated.
	// +optional
	Namespace Namespace `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`

	// UID is the unique in time and space value for this object. It is typically generated by
	// the server on successful creation of a resource and is not allowed to change on PUT
	// operations.
	//
	// Populated by the system.
	// Read-only.
	// +optional
	UID UID `json:"uid,omitempty" protobuf:"bytes,5,opt,name=uid"`
}

// ClusterSelectorOperator is the set of operators that can be used in
// a node selector requirement.
type ClusterSelectorOperator string

// ClusterSelectorRequirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type ClusterSelectorRequirement struct {
	// The label key that the selector applies to.
	Key string
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	Operator ClusterSelectorOperator
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	Values []string
}

// ClusterSelectorTerm represents expressions and fields required to select nodes.
// A null or empty node selector term matches no objects. The requirements of
// them are ANDed.
// The TopologySelectorTerm type implements a subset of the ClusterSelectorTerm.
type ClusterSelectorTerm struct {
	// A list of node selector requirements by node's labels.
	MatchExpressions []ClusterSelectorRequirement
	// A list of node selector requirements by node's fields.
	MatchFields []ClusterSelectorRequirement
}

// PreferredSchedulingTerm represents an empty preferred scheduling term matches all objects with implicit weight 0
// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
type PreferredSchedulingTerm struct {
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
	Weight int32
	// A node selector term, associated with the corresponding weight.
	Preference ClusterSelectorTerm
}

// ClusterAffinity is a group of node affinity scheduling rules.
type ClusterAffinity struct {
	// If the affinity requirements specified by this field are not met at
	// scheduling time, the job will not be scheduled onto the cluster.
	//
	// The scheduler will prefer to schedule jobs to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster matches the corresponding matchExpressions; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm
}

// JobAffinity is a group of inter pod anti affinity scheduling rules.
type JobAffinity struct {

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the job will not be scheduled onto the cluster.
	//
	// The scheduler will prefer to schedule jobs to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster matches the corresponding matchExpressions; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm
}

// JobAntiAffinity is a group of inter pod anti affinity scheduling rules.
type JobAntiAffinity struct {

	// If the anti-affinity  requirements specified by this field are not met at
	// scheduling time, the job will not be scheduled onto the cluster.
	//
	// The scheduler will prefer to schedule jobs to clusters that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a cluster that violates one or more of the expressions. The cluster that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the cluster matches the corresponding matchExpressions; the
	// cluster(s) with the highest sum are the most preferred.
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm
}

// Affinity is a group of affinity scheduling rules.
type Affinity struct {
	// Describes cluster affinity scheduling rules for the job.
	// +optional
	ClusterAffinity *ClusterAffinity
	// Describes job affinity scheduling rules (e.g. co-locate this job in the same cluster, zone, etc. as some other jobs(s)).
	// +optional
	JobAffinity *JobAffinity
	// Describes job anti-affinity scheduling rules (e.g. avoid putting this job in the same cluster, zone, etc. as some other job(s)).
	// +optional
	JobAntiAffinity *JobAntiAffinity
}
