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
	ctx              context.Context
	sampleProject    = ProjectEntry{Id: 12345, OwnerName: "amp", RepoName: "amp-repo", Token: "FakeToken"}
	updSampleProject = ProjectEntry{Id: 12345, OwnerName: "amp", RepoName: "amp-repo2", Token: "FakeToken"}
	sampleProject2   = ProjectEntry{Id: 12346, OwnerName: "amp", RepoName: "amp-repo", Token: "FakeToken"}
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")
	proj = createProjectServer()
	ctx = context.Background()
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	_, err := proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	if err != nil {
		t.Error(err)
	}
	resp, err := proj.Get(ctx, &GetRequest{Id: sampleProject.Id})
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, &sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	// Create new Entry
	proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	// Attempt to create a duplicate
	_, err := proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	// Should result in a duplicate entry
	if !strings.Contains(err.Error(), "key already exists") {
		t.Error("Duplicate Key Error was not detected")
	}
}

func TestGet(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	// Create new Entry
	proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	// Fetch The Entry
	resp, err := proj.Get(ctx, &GetRequest{Id: sampleProject.Id})
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, &sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
	}
}
func TestUpdate(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	// Create new Entry
	proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	// Update The Entry
	_, err := proj.Update(ctx, &UpdateRequest{Project: &updSampleProject})
	if err != nil {
		t.Error(err)
	}
	resp, _ := proj.Get(ctx, &GetRequest{Id: sampleProject.Id})
	if !proto.Equal(resp.Project, &updSampleProject) {
		t.Errorf("expected %v, got %v", updSampleProject, resp.Project)
	}
}

func TestList(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject2.Id})
	// Create new Entries
	proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	proj.Create(ctx, &CreateRequest{Project: &sampleProject2})
	// Fetch The Entry
	resp, err := proj.List(ctx, &ListRequest{})
	if err != nil {
		t.Error(err)
	} else if len(resp.Projects) != 2 {
		t.Errorf("Expected length of 2, got %v\n", len(resp.Projects))
	}

}

func TestDelete(t *testing.T) {
	// Guarantee no data exists
	proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	// Create new Entry
	proj.Create(ctx, &CreateRequest{Project: &sampleProject})
	// Delete The Entry
	resp, err := proj.Delete(ctx, &DeleteRequest{Id: sampleProject.Id})
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(resp.Project, &sampleProject) {
		t.Errorf("expected %v, got %v", sampleProject, resp.Project)
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
