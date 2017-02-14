package client

import (
	"fmt"
	"time"

	conf "github.com/appcelerator/amp/pkg/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

var (
	// amp is a singleton
	amp *AMP
)

func init() {
	grpclog.SetLogger(logger{})
}

// Logger is a simple log interface that also implements grpclog.Logger
type Logger interface {
	grpclog.Logger
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

// logger implements grpclog.Logger
type logger struct{}

func (l logger) Fatal(args ...interface{}) {
	if amp != nil {
		amp.Log.Fatal(args...)
	}
}
func (l logger) Fatalf(format string, args ...interface{}) {
	if amp != nil {
		amp.Log.Fatalf(format, args...)
	}
}
func (l logger) Fatalln(args ...interface{}) {
	if amp != nil {
		amp.Log.Fatalln(args...)
	}
}
func (l logger) Print(args ...interface{}) {
	if amp != nil {
		amp.Log.Print(args...)
	}
}
func (l logger) Printf(format string, args ...interface{}) {
	if amp != nil {
		amp.Log.Printf(format, args...)
	}
}
func (l logger) Println(args ...interface{}) {
	if amp != nil {
		amp.Log.Println(args...)
	}
}

// AMP holds the state for the current environment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *conf.Configuration

	// Conn is the gRPC connection to amplifier
	Conn *grpc.ClientConn

	// Log also implements the grpclog.Logger interface
	Log Logger
}

// Connect to amplifier
func (a *AMP) Connect() error {
	ampAddr := fmt.Sprintf("%s:%s", a.Configuration.AmpAddress, a.Configuration.ServerPort)
	conn, err := grpc.Dial(ampAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second))
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
func NewAMP(c *conf.Configuration, l Logger) *AMP {
	if amp == nil {
		amp = &AMP{Configuration: c, Log: l}
	}
	return amp
}
