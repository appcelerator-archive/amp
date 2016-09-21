package storage

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Interface must be implemented for a key/value store
type Interface interface {
	// Endpoints returns an array of endpoints for the storage
	Endpoints() []string

	// Connect to etcd using client v3 api
	Connect(timeout time.Duration) error

	// Close connection to etcd
	Close() error

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, val proto.Message, out proto.Message, ttl int64) error

	// Get unmarshals the protocol buffer message found at key into out, if found.
	// If not found and ignoreNotFound is set, then out will be a zero object, otherwise
	// error will be set to not found. A non-existing node or an empty response are both
	// treated as not found.
	Get(ctx context.Context, key string, out proto.Message, ignoreNotFound bool) error

	// List(ctx context.Context, key string, resourceVersion string, filter FilterFunc, list interface{}) error

	// TODO: will need to add preconditions support
	// if recurse then all the key having the same path under 'key' are going to be deleted
	// if !recurse then only 'key' is going to be deleted
	Delete(ctx context.Context, key string, recurse bool, out proto.Message) error

	// Update performs a guaranteed update, which means it will continue to retry until an update succeeds or the request is canceled.
	// Update(ctx context.Context, key string, type interface, ignoreNotFound bool, precondtions *Preconditions, tryUpdate UpdateFunc) error
	// TODO: the following is a temporary interface
	Update(ctx context.Context, key string, val proto.Message, ttl int64) error

	// List returns all the values that match the filter.
	List(ctx context.Context, key string, filter Filter, obj proto.Message, out *[]proto.Message) error

	// ListRaw returns all the raw entries matching the filter.
	ListRaw(ctx context.Context, key string, filter Filter) ([]*mvccpb.KeyValue, error)

	// Watch(ctx context.Context, key string, resourceVersion string, filter FilterFunc) (watch.Interface, error)

	// WatchList(ctx context.Context, key string, resourceVersion string, filter FilterFunc) (watch.Interface, error)

	Watch(ctx context.Context, key string, recursive bool) clientv3.WatchChan
	StopWatching() error
	// CompareAndSet atomically sets the value to the given updated value if the current value == the expected value
	CompareAndSet(ctx context.Context, key string, expect proto.Message, update proto.Message) error
}

// Filter is the interface used for storage operations that apply to sets (list, find, update).
type Filter interface {
	// Filter is a predicate that inspects a value (protocol buffer message instance) and returns true if and only if the value should remain in the set
	Filter(val proto.Message) bool
}

// Everything is a Filter which accepts every object.
var Everything Filter = everything{}

// everything is an implementation of Everything.
type everything struct {
}

// Filter implements the Filter interface to accept every object.
func (e everything) Filter(val proto.Message) bool {
	return true
}
