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
	"github.com/golang/protobuf/proto"
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
func (s *etcd) Create(ctx context.Context, key string, val proto.Message, out proto.Message, ttl int64) error {
	key = s.prefix(key)

	opts, err := s.options(ctx, int64(ttl))
	if err != nil {
		return err
	}

	data, err := proto.Marshal(val)

	txn, err := s.client.KV.Txn(ctx).
		If(notFound(key)).
		Then(clientv3.OpPut(key, string(data), opts...)).
		Commit()

	if err != nil {
		return err
	}

	if !txn.Succeeded {
		return fmt.Errorf("key already exists: %q", key)
	}

	if out != nil {
		// TODO: out will be the encoded message, revision comes from resp header
		putResp := txn.Responses[0].GetResponsePut()
		kv := putResp.PrevKv
		if kv != nil && kv.Value != nil {
			proto.Unmarshal(kv.Value, out)
		}
	}

	return nil
}

// Get implements storage.Interface.Get.
// TODO: out will be a protocol buffer message
func (s *etcd) Get(ctx context.Context, key string, out proto.Message, ignoreNotFound bool) error {
	if out == nil {
		return fmt.Errorf("`out` param must not be nil")
	}

	key = s.prefix(key)

	getResp, err := s.client.KV.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(getResp.Kvs) == 0 {
		if ignoreNotFound {
			// TODO: do we want to ignore out or set it to an empty message
			// if out != nil {
			//	*out.Reset()
			// }
			return nil
		}
		return fmt.Errorf("key not found: %q", key)
	}

	kv := getResp.Kvs[0]
	data := []byte(kv.Value)
	return proto.Unmarshal(data, out)
}

// Delete implements storage.Interface.Delete
func (s *etcd) Delete(ctx context.Context, key string, out proto.Message) error {
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
	data := []byte(kv.Value)
	return proto.Unmarshal(data, out)
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
