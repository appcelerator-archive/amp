package etcd

import (
	"encoding/json"
	"github.com/coreos/etcd/client"
	"github.com/satori/go.uuid"
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

// Connect to the elastic search server
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

// NewKey Generates a new key with the given value in the given keyspace and returns the key
func (etcd *Etcd) NewKey(keyPrefix string, value interface{}) (key string, err error) {
	key = keyPrefix + "/" + uuid.NewV4().String()
	json, err := json.Marshal(value)
	if err != nil {
		return
	}

	resp, err := kapi.Set(context.Background(), key, string(json), nil)
	if err != nil {
		return
	}
	log.Printf("Set is done. Metadata is %q\n", resp)
	return
}

// All get all the key/value pairs (nodes) under the given path.
func (etcd *Etcd) All(path string) (nodes client.Nodes, err error) {
	resp, err := kapi.Get(context.Background(), path, &client.GetOptions{Recursive: true, Quorum: true})
	if err != nil {
		return
	}
	log.Printf("All is done. Metadata is %q\n", resp)
	nodes = resp.Node.Nodes
	return
}
