package etcd

import (
	"encoding/json"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

var (
	cfg = client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	kapi client.KeysAPI
)

// Etcd singlet	on
type Etcd struct {
}

// Connect to the etcd server
func (etcd *Etcd) Connect(endpoints ...string) (err error) {
	if len(endpoints) > 0 {
		cfg.Endpoints = endpoints
	}
	c, err := client.New(cfg)
	if err != nil {
		return
	}
	kapi = client.NewKeysAPI(c)
	log.Printf("Successfully Connected to etcd on: %+v\n", cfg.Endpoints)
	return
}

// Put Puts a new key with the given value in the given keyspace and returns the key path
func (etcd *Etcd) Put(prefix string, id string, value interface{}) (path string, err error) {
	path = prefix + "/" + id
	json, err := json.Marshal(value)
	if err != nil {
		return
	}

	_, err = kapi.Set(context.Background(), path, string(json), nil)
	if err != nil {
		return
	}
	log.Printf("Successfully put a new key @ %v\n", path)
	return
}

// List lists all the key/value pairs (nodes) under the given path
func (etcd *Etcd) List(path string) (nodes client.Nodes, err error) {
	resp, err := kapi.Get(context.Background(), path, &client.GetOptions{Recursive: true, Quorum: true})
	if err != nil {
		return
	}
	return resp.Node.Nodes, nil
}
