package server

import (
	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

const (
	address = "localhost:50051"
)

func TestShouldSucceedWhenProvidingAValidCreateRequest(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	// Create the client
	c := project.NewProjectClient(conn)

	// Create a project
	name := "test-" + uuid.NewV4().String()
	p := project.Project{Name: name}
	reply, err := c.Create(context.Background(), &project.CreateRequest{Project: &p})
	if err != nil {
		t.Fatalf("could not create: %v", err)
	}
	assert.Equal(t, name, reply.Created.Name)
	assert.NotNil(t, name, reply.Created.Id)
	conn.Close()
}

func TestShouldListAllProjects(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	// Create the client
	c := project.NewProjectClient(conn)

	// Create projects
	names := [10]string{}
	for i := 0; i < 10; i++ {
		names[i] = "test-" + uuid.NewV4().String()
		p := project.Project{
			Name: names[i],
		}
		_, err := c.Create(context.Background(), &project.CreateRequest{Project: &p})
		if err != nil {
			t.Fatalf("could not create: %v", err)
		}
	}

	// List projects
	reply, err := c.List(context.Background(), &project.ListRequest{})
	if err != nil {
		t.Fatalf("could not create: %v", err)
	}

	// Construct a list of all names
	allNames := make([]string, len(reply.Projects))
	for i, p := range reply.Projects {
		allNames[i] = p.Name
	}

	// Assert we got ours
	for i := 0; i < 10; i++ {
		assert.Contains(t, allNames, names[i])
	}

	conn.Close()
}
