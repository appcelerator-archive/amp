// Packaged etcd was influenced by and borrows helper functions from:
// https://github.com/kubernetes/kubernetes/pkg/storage/etcd3
/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package etcd

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/storage"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

// etcd is used to connect to and query etcd
type etcd struct {
	client     *clientv3.Client
	endpoints  []string
	pathPrefix string
}

// New returns an etcd implementation of storage.Interface
func New(endpoints []string, prefix string) storage.Interface {
	return &etcd{endpoints: endpoints, pathPrefix: prefix}
}

// Endpoints gets the endpoints etcd
func (s *etcd) Endpoints() []string {
	return s.endpoints
}

// Connect to etcd using client v3 api
func (s *etcd) Connect(timeout time.Duration) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   s.endpoints,
		DialTimeout: timeout,
	})
	s.client = cli
	return err
}

// Close connection to etcd
func (s *etcd) Close() error {
	if err := s.client.Close(); err != nil {
		return err
	}
	s.client = nil
	return nil
}

// Create implements storage.Interface.Create
// TODO: val, out will be protocol buffer messages
func (s *etcd) Create(ctx context.Context, key string, val interface{}, out *string, ttl int64) error {
	key = s.prefix(key)

	opts, err := s.options(ctx, int64(ttl))
	if err != nil {
		return err
	}

	// TODO: cast to string is just to get basic mechanics in place before using proto buf messages
	var str string
	var ok bool
	if str, ok = val.(string); !ok {
		return fmt.Errorf("val must be a string but was type (%T): %v", val, val)
	}

	txn, err := s.client.KV.Txn(ctx).
		If(notFound(key)).
		Then(clientv3.OpPut(key, str, opts...)).
		Commit()

	if err != nil {
		return err
	}

	if !txn.Succeeded {
		return fmt.Errorf("key already exists: %q", key)
	}

	if out != nil {
		// TODO: out will be the encoded message, revision comes from resp header
		//putResp := txn.Responses[0].GetResponsePut()
		*out = str
	}

	return nil
}

// Get implements storage.Interface.Get.
// TODO: out will be a protocol buffer message
func (s *etcd) Get(ctx context.Context, key string, out *string, ignoreNotFound bool) error {
	key = s.prefix(key)

	getResp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(getResp.Kvs) == 0 {
		if ignoreNotFound {
			*out = key
			return nil
		}
		return fmt.Errorf("key not found: %q", key)
	}

	kv := getResp.Kvs[0]
	*out = string(kv.Value)
	return nil
}

// Delete implements storage.Interface.Delete
func (s *etcd) Delete(ctx context.Context, key string, out *string) error {
	key = s.prefix(key)

	txn, err := s.client.KV.Txn(ctx).
		If().
		Then(clientv3.OpGet(key), clientv3.OpDelete(key)).
		Commit()
	if err != nil {
		return err
	}

	getResp := txn.Responses[0].GetResponseRange()
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("key not found: %q", key)
	}
	kv := getResp.Kvs[0]
	*out = string(kv.Value)
	return nil
}

// options returns a slice of client options (currently just a lease based on the given ttl).
// ttl: time in seconds that key will exist (0 means forever); if ttl is non-zero, it will attach the key to a lease with ttl of roughly the same length
func (s *etcd) options(ctx context.Context, ttl int64) ([]clientv3.OpOption, error) {
	if ttl == 0 {
		return nil, nil
	}
	// TODO: one lease per key is expensive. Analyze further; it should possible to associate keys with the same lease (eg, a lease pool) that share the same ttl window
	lcr, err := s.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return nil, err
	}
	return []clientv3.OpOption{clientv3.WithLease(clientv3.LeaseID(lcr.ID))}, nil
}

func (s *etcd) prefix(key string) string {
	if strings.HasPrefix(key, s.pathPrefix) {
		return key
	}
	return path.Join(s.pathPrefix, key)
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}
