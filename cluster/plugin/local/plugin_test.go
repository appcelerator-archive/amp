package local

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/docker/docker/client"
)

const (
	defaultURL     = "unix:///var/run/docker.sock"
	defaultVersion = "1.30"
)

var (
	testClient *client.Client
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	testClient, _ = client.NewClient(defaultURL, defaultVersion, nil, nil)
}

func teardown() {
	// TODO destroy test stacks
}

func TestCreate(t *testing.T) {

	opts := &RequestOptions{}

	log.Println("starting test...")

	// create cluster
	// ============
	ctxCreate := context.Background()
	err := EnsureSwarmExists(ctxCreate, testClient, opts)
	if err != nil {
		t.Fatal(err)
	}

	// describe cluster
	// ============
	ctxDescribe := context.Background()
	output, err := InfoCluster(ctxDescribe, testClient)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("-------------------------------------------------------")
	log.Println(output)
	log.Println("-------------------------------------------------------")

	// delete cluster
	// ============
	ctxDelete := context.Background()
	err = DeleteCluster(ctxDelete, testClient, opts)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("cluster deleted\n")
}
