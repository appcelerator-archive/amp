package client

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

const (
	defaultPort = ":50101"
	// ServerAddress the server address
	ServerAddress = "localhost" + defaultPort
)

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose bool
	Github  string
	Target  string
	Images  []string
}

// AMP holds the state for the current environment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *Configuration
	Conn          *grpc.ClientConn
}

func (a *AMP) connect() {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	a.Conn = conn
}

func (a *AMP) disconnect() {
	err := a.Conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GetAuthorizedContext returns an authorized context
func (a *AMP) GetAuthorizedContext() (ctx context.Context, err error) {
	if a.Configuration.Github == "" {
		return nil, fmt.Errorf("Requires login")
	}
	md := metadata.Pairs("sessionkey", a.Configuration.Github)
	ctx = metadata.NewContext(context.Background(), md)
	return
}

// Verbose returns true if verbose flag is set
func (a *AMP) Verbose() bool {
	return a.Configuration.Verbose
}

// NewAMP creates a new AMP instance
func NewAMP(c *Configuration) *AMP {
	return &AMP{Configuration: c}
}

// Create a new swarm
func (a *AMP) Create() {
	if a.Verbose() {
		fmt.Println("Create")
	}
}

// Start the swarm
func (a *AMP) Start() {
	if a.Verbose() {
		fmt.Println("Start")
	}
}

// Update the swarm
func (a *AMP) Update() {
	if a.Verbose() {
		fmt.Println("Update")
	}
}

// Stop the swarm
func (a *AMP) Stop() {
	if a.Verbose() {
		fmt.Println("Stop")
	}
}

// Status returns the current status
func (a *AMP) Status() {
	if a.Verbose() {
		fmt.Println("Status")
	}
}