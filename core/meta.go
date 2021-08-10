package core

import "time"

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

// Object lets you work with object metadata from any of the versioned or
// internal API objects. Attempting to set or retrieve a field on an object that does
// not support that field (Name, UID, Namespace on lists) will be a no-op and return
// a default value.
type Object interface {
	GetNamespace() string
	SetNamespace(namespace string)
	GetName() string
	SetName(name string)
	GetGenerateName() string
	SetGenerateName(name string)
	GetUID() UID
	SetUID(uid UID)
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetCreationTimestamp() time.Time
	SetCreationTimestamp(timestamp time.Time)
	GetDeletionTimestamp() *time.Time
	SetDeletionTimestamp(timestamp *time.Time)
	GetDeletionGracePeriodSeconds() *int64
	SetDeletionGracePeriodSeconds(*int64)
	GetClusterName() string
	SetClusterName(clusterName string)
}
