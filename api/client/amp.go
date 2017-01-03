package client

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

const (
	//DefaultServerAddress amplifier address + port default
	DefaultServerAddress = "127.0.0.1:8080"
	//DefaultAdminServerAddress adm-server address + port default
	DefaultAdminServerAddress = "127.0.0.1:31315"
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

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose            bool
	GitHub             string
	Target             string
	Port               string
	ServerAddress      string
	AdminServerAddress string
	CmdTheme           string
}

// AMP holds the state for the current environment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *Configuration

	// Conn is the gRPC connection to amplifier
	Conn *grpc.ClientConn

	// Log also implements the grpclog.Logger interface
	Log Logger
}

// Connect to amplifier
func (a *AMP) Connect() error {
	conn, err := grpc.Dial(a.Configuration.ServerAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second))
	if err != nil {
		return fmt.Errorf("Error connecting to amplifier @ %s: %v", a.Configuration.ServerAddress, err)
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
func NewAMP(c *Configuration, l Logger) *AMP {
	if amp == nil {
		amp = &AMP{Configuration: c, Log: l}
	}
	return amp
}
