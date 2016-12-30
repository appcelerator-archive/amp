package client

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

const (
	//DefaultServerAddress amplifier address + port default
	DefaultServerAddress = "127.0.0.1:8080"
)

var (
	verbose = false
)

func init() {
	grpclog.SetLogger(logger{})
}

type logger struct{}

func (l logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}
func (l logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
func (l logger) Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}
func (l logger) Print(args ...interface{}) {
	if verbose {
		log.Print(args...)
	}
}
func (l logger) Printf(format string, args ...interface{}) {
	if verbose {
		log.Printf(format, args...)
	}
}
func (l logger) Println(args ...interface{}) {
	if verbose {
		log.Println(args...)
	}
}

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose            bool
	Github             string
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
	Conn          *grpc.ClientConn
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
	verbose = c.Verbose
	return &AMP{Configuration: c}
}
