package server

import (
	"github.com/appcelerator/amp/api/rpc/project"
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

	// Contact the server and print out its response.
	c := project.NewProjectClient(conn)
	r, err := c.Create(context.Background(), &project.CreateRequest{Name: "test-project"})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	t.Logf("Greeting: %s", r.Message)

	conn.Close()
}
