package main

import (
	"github.com/appcelerator/amp/api/rpc/project"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

const (
	address = "localhost:50051"
)

func TestShouldGetAHundredLogEntriesByDefault(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := project.NewProjectClient(conn)

	// Contact the server and print out its response.
	name := "test-project"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.CreateProject(context.Background(), &project.CreateRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	conn.Close()
}
