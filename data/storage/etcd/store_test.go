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
	defTimeout          = 5 * time.Second
	defaultPort         = ":50101"
	etcdDefaultEndpoint = "http://localhost:2379"
)

var (
	store         storage.Interface
	port          string
	etcdEndpoints = []string{etcdDefaultEndpoint}
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
	out := &storage.Project{}
	ignoreNotFound := false

	err := store.Get(ctx, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	expected := &storage.Project{Id: "100", Name: "AMP"}
	if !proto.Equal(expected, out) {
		t.Errorf("expected %v, got %v", expected, out)
	}
}

func TestGetWithError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foobar"
	out := &storage.Project{}
	ignoreNotFound := false

	err := store.Get(ctx, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err == nil {
		t.Errorf("expected an error result")
	}
}

func TestGetIgnoreError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foobar"
	out := &storage.Project{}
	ignoreNotFound := true

	err := store.Get(ctx, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Errorf("error not expected: %s", err)
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	out := &storage.Project{}

	err := store.Delete(ctx, key, out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	expected := &storage.Project{Id: "100", Name: "AMP"}
	if !proto.Equal(expected, out) {
		t.Errorf("expected %v, got %v", expected, out)
	}
}

func TestUpdate(t *testing.T) {
	key := "foo"
	val := &storage.Project{Id: "100", Name: "bar"}
	ttl := int64(0)

	ctx1, cancel1 := newContext()
	err := store.Update(ctx1, key, val, ttl)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel1()
	if err != nil {
		t.Error(err)
	}

	out := &storage.Project{}
	ignoreNotFound := false
	ctx2, cancel2 := newContext()
	err = store.Get(ctx2, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel2()
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val, out) {
		t.Errorf("expected %v, got %v", val, out)
	}

	ctx3, cancel3 := newContext()
	err = store.Delete(ctx3, key, out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel3()
	if err != nil {
		t.Error(err)
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defTimeout)
}
