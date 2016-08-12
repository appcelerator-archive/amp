package cli

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/logs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

const (
	defaultPort   = ":50101"
	serverAddress = "localhost" + defaultPort
)

// AMP holds the state for the current envirionment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Config *Config
}

func (a *AMP) verbose() bool {
	return a.Config.Verbose
}

// NewAMP creates a new AMP instance
func NewAMP(c *Config) *AMP {
	return &AMP{Config: c}
}

// Create a new swarm
func (a *AMP) Create() {
	if a.verbose() {
		fmt.Println("Create")
	}
}

// Start the swarm
func (a *AMP) Start() {
	if a.verbose() {
		fmt.Println("Start")
	}
}

// Update the swarm
func (a *AMP) Update() {
	if a.verbose() {
		fmt.Println("Update")
	}
}

// Stop the swarm
func (a *AMP) Stop() {
	if a.verbose() {
		fmt.Println("Stop")
	}
}

// Status returns the current status
func (a *AMP) Status() {
	if a.verbose() {
		fmt.Println("Status")
	}
}

// Logs returns the logs
func (a *AMP) Logs() {
	if a.verbose() {
		fmt.Println("Logs")
	}
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// Contact the server and print out its response.
	c := logs.NewLogsClient(conn)
	r, err := c.Get(context.Background(), &logs.GetRequest{})
	if err != nil {
		panic(err)
	}
	for _, entry := range r.Entries {
		log.Printf("%+v", entry)
	}
	conn.Close()
}
