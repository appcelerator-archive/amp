package cli

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	// amp is a singleton
	amp *AMP
)

// AMP holds the state for the current environment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *Configuration

	// Conn is the gRPC connection to amplifier
	Conn *grpc.ClientConn

	// Log also implements the grpclog.Logger interface
	Log *Logger
}

// Connect to amplifier
func (a *AMP) Connect() error {
	ampAddr := fmt.Sprintf("%s:%s", a.Configuration.AmpAddress, a.Configuration.ServerPort)
	conn, err := grpc.Dial(ampAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		grpc.WithPerRPCCredentials(GetLoginCredentials()),
	)
	if err != nil {
		return fmt.Errorf("Error connecting to amplifier @ %s: %v", ampAddr, err)
	}
	a.Conn = conn
	return nil
}

// Disconnect from amplifier
func (a *AMP) Disconnect() {
	if a.Conn == nil {
		return
	}
	err := a.Conn.Close()
	if err != nil {
		a.Log.Panic(err)
	}
}

// GetAuthorizedContext returns an authorized context
func (a *AMP) GetAuthorizedContext() (ctx context.Context, err error) {
	// TODO: reenable
	// Disabled temporally
	// if a.Configuration.Github == "" {
	// 	return nil, fmt.Errorf("Requires login")
	// }
	md := metadata.Pairs("sessionkey", a.Configuration.GitHub)
	ctx = metadata.NewContext(context.Background(), md)
	return
}

// Verbose returns true if verbose flag is set
func (a *AMP) Verbose() bool {
	return a.Configuration.Verbose
}

// NewAMP creates an AMP singleton instance
// (will only be configured with the first call)
func NewAMP(c *Configuration, l *Logger) *AMP {
	if amp == nil {
		amp = &AMP{Configuration: c, Log: l}
	}
	return amp
}
