package etcd_test

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const (
	defTimeout              = 5 * time.Second
	defaultPort             = ":50101"
	etcdDefaultEndpoint     = "http://localhost:2379"
)

var (
	store            storage.Interface
	port             string
	etcdEndpoints    = []string{etcdDefaultEndpoint}
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")

	if endpoints := os.Getenv("endpoints"); endpoints != "" {
		etcdEndpoints = strings.Split(endpoints, ",")
	}

	store = etcd.New(etcdEndpoints, "amp")
	if err := store.Connect(5 * time.Second); err != nil {
		panic(err)
	}
	log.Printf("connected to etcd at %v", strings.Join(store.Endpoints(), ","))

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := &storage.Project{Id: "100", Name: "AMP"}
	out := &storage.Project{}
	ttl := int64(0)

	err := store.Create(ctx, key, val, out, ttl)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	// on creation there should be no previous value
	if proto.Equal(val, out) {
		t.Errorf("expected %v, got %v", val, out)
	}
}

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := &storage.Project{Id: "100", Name: "AMP"}
	out := &storage.Project{}
	ignoreNotFound := false

	err := store.Get(ctx, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	if !proto.Equal(val, out) {
		t.Errorf("expected %v, got %v", val, out)
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := &storage.Project{Id: "100", Name: "AMP"}
	out := &storage.Project{}

	err := store.Delete(ctx, key, out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	if !proto.Equal(val, out) {
		t.Errorf("expected %v, got %v", val, out)
	}
}
