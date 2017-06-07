package etcd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/state"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const (
	defTimeout = 5 * time.Second
)

var (
	store storage.Interface
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")

	etcdEndpoints := []string{etcd.DefaultEndpoint}
	log.Printf("connecting to etcd at %s", strings.Join(etcdEndpoints, ","))
	store = etcd.New(etcdEndpoints, "amp", defTimeout)
	if err := store.Connect(); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", etcdEndpoints, err)
	}
	log.Printf("connected to etcd at %v", strings.Join(store.Endpoints(), ","))
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := &TestMessage{Id: "100", Name: "AMP"}
	out := &TestMessage{}
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
	out := &TestMessage{}
	ignoreNotFound := false

	err := store.Get(ctx, key, out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	expected := &TestMessage{Id: "100", Name: "AMP"}
	if !proto.Equal(expected, out) {
		t.Errorf("expected %v, got %v", expected, out)
	}
}

func TestGetWithError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foobar"
	out := &TestMessage{}
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
	out := &TestMessage{}
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
	out := &TestMessage{}

	err := store.Delete(ctx, key, false, out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	expected := &TestMessage{Id: "100", Name: "AMP"}
	if !proto.Equal(expected, out) {
		t.Errorf("expected %v, got %v", expected, out)
	}
}

func TestConcurrentUpdates(t *testing.T) {
	ctx := context.Background()
	key := "entry"

	// cleanup
	if err := store.Delete(ctx, key, false, nil); err != nil {
		if err != storage.NotFound {
			t.Error(err)
		}
	}

	// Create an empty entry
	value := &TestMessage{}
	store.Create(ctx, key, value, nil, 0)

	// Perform updates concurrently
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func(i int) {
			defer wg.Done()
			uf := func(current proto.Message) (proto.Message, error) {
				msg, ok := current.(*TestMessage)
				if !ok {
					return nil, fmt.Errorf("value is not the right type (expected TestMessage): %T", msg)
				}

				msg.Names = append(msg.Names, fmt.Sprintf("%d", i))
				return msg, nil
			}
			if err := store.Update(ctx, key, uf, &TestMessage{}); err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

	current := &TestMessage{}
	if err := store.Get(ctx, key, current, false); err != nil {
		t.Error(err)
	}
	assert.Equal(t, len(current.Names), iterations)

	// cleanup
	if err := store.Delete(ctx, key, false, nil); err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	// generic context
	ctx := context.Background()

	// store everything under amp/foo/
	key := "foo"

	// this is a "template" object that provides a concrete type for list to unmarshal into
	obj := &TestMessage{}

	// unlimited ttl
	ttl := int64(0)

	// will store values that we store, which we will use to compare list results against
	vals := []*TestMessage{}

	// will store the results of calling list
	var out []proto.Message

	// this will create a bunch of TestMessage items to store in etcd
	// it will also save them to vals so we can compare them against what
	// we get back when we call the list method to ensure actual matches expected
	for i := 0; i < 5; i++ {
		id := strconv.Itoa(i)
		name := fmt.Sprintf("bar%d", i)
		subkey := path.Join(key, id)
		vals = append(vals, &TestMessage{Id: id, Name: name})

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
		m, ok := out[i].(*TestMessage)
		if !ok {
			t.Errorf("value is not the right type (expected TestMessage): %T", out[i])
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
	ctx := context.Background()

	key := "state"
	expect := &state.State{Value: "hello"}
	update := &state.State{Value: "world"}

	store.Delete(ctx, key, false, &state.State{})
	store.Create(ctx, key, expect, nil, 0)
	if err := store.CompareAndSet(ctx, key, expect, update); err != nil {
		t.Error(err)
	}
	actual := &state.State{}
	if err := store.Get(ctx, key, actual, false); err != nil {
		t.Error(err)
	}
	if !proto.Equal(update, actual) {
		t.Errorf("expected %v, got %v", update, actual)
	}
}
