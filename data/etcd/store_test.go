package etcd_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/data"
)

const (
	defTimeout = 5 * time.Second
)

var (
	config server.Config = server.Config{
		Port:          ":50101",
		EtcdEndpoints: []string{"http://localhost:2379"},
	}

	store data.Store
)

func TestMain(m *testing.M) {
	go server.Start(config)

	// there is no event when the server starts listening, so we just wait a second
	time.Sleep(1 * time.Second)
	store = server.Store

	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := "bar"
	var out string
	ttl := int64(0)

	err := store.Create(ctx, key, val, &out, ttl)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	if out != val {
		t.Errorf("expected %q, got %q", val, out)
	}
}

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := "bar"
	var out string
	ignoreNotFound := false

	err := store.Get(ctx, key, &out, ignoreNotFound)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	if out != val {
		t.Errorf("expected %q, got %q", val, out)
	}
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), defTimeout)
	key := "foo"
	val := "bar"
	var out string

	err := store.Delete(ctx, key, &out)
	// cancel timeout (release resources) if operation completes before timeout
	defer cancel()
	if err != nil {
		t.Error(err)
	}

	if out != val {
		t.Errorf("expected %q, got %q", val, out)
	}
}
