package etcd_test

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack"
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

	err := store.Delete(ctx, key, false, out)
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

	// confirm
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

	// cleanup
	ctx3, cancel3 := newContext()
	err = store.Delete(ctx3, key, false, out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel3()
	if err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	// generic context
	ctx := context.Background()

	// store everything under amp/foo/
	key := "foo"

	// this is a "template" object that provides a concrete type for list to unmarshal into
	obj := &storage.Project{}

	// unlimited ttl
	ttl := int64(0)

	// will store values that we store, which we will use to compare list results against
	vals := []*storage.Project{}

	// will store the results of calling list
	var out []proto.Message

	// this will create a bunch of storage.Project items to store in etcd
	// it will also save them to vals so we can compare them against what
	// we get back when we call the list method to ensure actual matches expected
	for i := 0; i < 5; i++ {
		id := strconv.Itoa(i)
		name := fmt.Sprintf("bar%d", i)
		subkey := path.Join(key, id)
		vals = append(vals, &storage.Project{Id: id, Name: name})

		err := store.Create(ctx, subkey, vals[i], obj, ttl)
		if err != nil {
			t.Error(err)
		}
	}

	// use list with the Everything filter to fetch all the items we stored
	err := store.List(ctx, key, storage.Everything, obj, &out)
	if err != nil {
		t.Error(err)
	}

	// compare the list that we got back (out) with the items we created (vals)
	for i := 0; i < len(out); i++ {
		if !proto.Equal(vals[i], out[i]) {
			t.Errorf("expected %v, got %v", vals[i], out[i])
		}

		// actually inspect individual message contents
		// Unfortunately, without generics in Go, this requires a type assertion
		m, ok := out[i].(*storage.Project)
		if !ok {
			t.Errorf("value is not the right type (expected storage.Project): %T", out[i])
		}
		name := fmt.Sprintf("bar%d", i)
		if m.Name != name {
			t.Errorf("id: %d, name: %s does not match expected name: %s", m.Id, m.Name, name)
		}

		// clean up after ourselves -- delete the item
		id := strconv.Itoa(i)
		subkey := path.Join(key, id)
		err := store.Delete(ctx, subkey, false, obj)
		if err != nil {
			t.Error(err)
		}
		// confirm the deleted object matches the original object
		if !proto.Equal(vals[i], obj) {
			t.Errorf("expected %v, deleted %v", vals[i], obj)
		}
	}
}

func TestCompareAndSet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "state"
	expect := &stack.State{Value: stack.StackState_Stopped}
	update := &stack.State{Value: stack.StackState_Running}

	err := store.Delete(ctx, key, false, &stack.State{})
	err = store.Create(ctx, key, expect, nil, 0)
	err = store.CompareAndSet(ctx, key, expect, update)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	actual := &stack.State{}
	err = store.Get(ctx, key, actual, false)
	if !proto.Equal(update, actual) {
		t.Errorf("expected %v, got %v", update, actual)
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defTimeout)
}
