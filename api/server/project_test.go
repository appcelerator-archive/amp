package server

import (
	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
)

const (
	address = "localhost:50051"
	port = ":50051"
)

func TestShouldSucceedWhenProvidingAValidCreateRequest(t *testing.T) {
	server.Start(port)
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}

	// Contact the server and print out its response.
	c := project.NewProjectClient(conn)
	r, err := c.Create(context.Background(), &project.CreateRequest{Name: "test-project-" + uuid.NewV4().String()})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	t.Logf("Greeting: %s", r.Message)

	conn.Close()
}
