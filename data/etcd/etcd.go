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
func (etcd *Etcd) Connect(endpoints ...string) {
	if len(endpoints) > 0 {
		cfg.Endpoints = endpoints
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi = client.NewKeysAPI(c)
	log.Printf("Successfully Connected to etcd on: %+v\n", cfg.Endpoints)
}

// SetKey set the value for a given key
func (etcd *Etcd) SetKey(keyPrefix string, value interface{}) (key string) {
	key = keyPrefix + "/" + uuid.NewV4().String()
	json, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := kapi.Set(context.Background(), key, string(json), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
		log.Printf("Set is done. Data is %+v\n", value)
	}
	return
}

// GetKey get the value for a given key
func (etcd *Etcd) GetKey(key string) (value string) {
	resp, err := kapi.Get(context.Background(), key, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	value = resp.Node.Value
	return
}
