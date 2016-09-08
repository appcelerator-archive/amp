package project

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/golang/protobuf/proto"
)

const (
	etcdDefaultEndpoint = "http://localhost:2379"
)

var (
	etcdEndpoints    = []string{etcdDefaultEndpoint}
	proj             *Proj
	sampleProject    = &ProjectRequest{&Project{RepoId: 12345, OwnerName: "amp", RepoName: "amp-repo", Token: "FakeToken"}}
	updSampleProject = &ProjectRequest{&Project{RepoId: 12345, OwnerName: "amp", RepoName: "amp-repo2", Token: "FakeToken"}}
	sampleProject2   = &ProjectRequest{&Project{RepoId: 12346, OwnerName: "amp", RepoName: "amp-repo", Token: "FakeToken"}}
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")
	proj = createProjectServer()
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	proj.Delete(context.Background(), sampleProject)
	resp, err := proj.Create(context.Background(), sampleProject)
	if err != nil {
		t.Error(err)
	}
	resp, err = proj.Get(context.Background(), sampleProject)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, sampleProject.Project) {
		t.Errorf("expected %v, got %v", sampleProject.Project, resp.Project)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(context.Background(), sampleProject)
	// Create new Entry
	proj.Create(context.Background(), sampleProject)
	// Attempt to create a duplicate
	_, err := proj.Create(context.Background(), sampleProject)
	// Should result in a duplicate entry
	if !strings.Contains(err.Error(), "key already exists") {
		t.Error("Duplicate Key Error was not detected")
	}
}

func TestGet(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(context.Background(), sampleProject)
	// Create new Entry
	proj.Create(context.Background(), sampleProject)
	// Fetch The Entry

	resp, err := proj.Get(context.Background(), sampleProject)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, sampleProject.Project) {
		t.Errorf("expected %v, got %v", sampleProject.Project, resp.Project)
	}
}
func TestUpdate(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(context.Background(), sampleProject)
	// Create new Entry
	proj.Create(context.Background(), sampleProject)
	// Update The Entry
	resp, err := proj.Update(context.Background(), updSampleProject)
	if err != nil {
		t.Error(err)
	}
	resp, _ = proj.Get(context.Background(), sampleProject)
	if !proto.Equal(resp.Project, updSampleProject.Project) {
		t.Errorf("expected %v, got %v", updSampleProject.Project, resp.Project)
	}
}

func TestList(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(context.Background(), sampleProject)
	proj.Delete(context.Background(), sampleProject2)
	// Create new Entries
	proj.Create(context.Background(), sampleProject)
	proj.Create(context.Background(), sampleProject2)
	// Fetch The Entry
	resp, err := proj.List(context.Background(), &Empty{})
	if err != nil {
		t.Error(err)
	} else if len(resp.Projects) != 2 {
		t.Errorf("Expected length of 2, got %v\n", len(resp.Projects))
	}

}

func TestDelete(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(context.Background(), sampleProject)
	// Create new Entry
	proj.Create(context.Background(), sampleProject)
	// Delete The Entry
	resp, err := proj.Delete(context.Background(), sampleProject)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, sampleProject.Project) {
		t.Errorf("expected %v, got %v", sampleProject.Project, resp.Project)
	}
}

// Create the ProjectServer as a local struct that can be excercised directly over the call stack
func createProjectServer() *Proj {
	//Create the config
	var proj = &Proj{}

	if endpoints := os.Getenv("endpoints"); endpoints != "" {
		etcdEndpoints = strings.Split(endpoints, ",")
	}
	proj.Store = etcd.New(etcdEndpoints, "amp")
	if err := proj.Store.Connect(5 * time.Second); err != nil {
		panic(err)
	}
	return proj
}
