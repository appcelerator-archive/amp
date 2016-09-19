package client

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

const (
	//DefaultServerAddress amplifier address + port default
	DefaultServerAddress = "localhost:50101"
)

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose       bool
	Github        string
	Target        string
	Images        []string
	Port          string
	ServerAddress string
}

// AMP holds the state for the current environment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *Configuration
	Conn          *grpc.ClientConn
}

// Connect to amplifier
func (a *AMP) Connect() *grpc.ClientConn {
	conn, err := grpc.Dial(a.Configuration.ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}
	a.Conn = conn
	return conn
}

// Disconnect from amplifier
func (a *AMP) Disconnect() {
	err := a.Conn.Close()
	if err != nil {
		log.Panic(err)
	}
}

// GetAuthorizedContext returns an authorized context
func (a *AMP) GetAuthorizedContext() (ctx context.Context, err error) {
	// TODO: reenable
	// Disabled temporally
	// if a.Configuration.Github == "" {
	// 	return nil, fmt.Errorf("Requires login")
	// }
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
