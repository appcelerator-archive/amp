package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Interface for the CLI's functionality.
type Interface interface {
	Version() string
	Build() string

	Console() *Console
	In() *InStream
	Out() *OutStream
	Err() io.Writer
	ShowHelp(cmd *cobra.Command, args []string) error

	Server() string
	SetServer(server string)
	SetSkipVerify(skipVerify bool)
	Connect() (*grpc.ClientConn, error)
	ClientConn() *grpc.ClientConn

	OnInitialize(initializers ...func())
}

// cli implements cli.Interface
type cli struct {
	Configuration
	version    string // nolint: structcheck, unused
	build      string // nolint: structcheck, unused
	console    *Console
	in         *InStream
	out        *OutStream
	err        io.Writer
	clientConn *grpc.ClientConn
	skipVerify bool
}

const (
	// DefaultAddress for amp connection
	DefaultAddress = "127.0.0.1"

	// DefaultPort for amp connection
	DefaultPort = configuration.DefaultPort
)

// NewCLI returns a new CLI instance.
func NewCLI(in io.ReadCloser, out, err io.Writer, config *Configuration) Interface {
	c := &cli{
		Configuration: *config,
		in:            NewInStream(in),
		out:           NewOutStream(out),
		err:           err,
	}
	c.console = NewConsole(c.Out(), config.Verbose)
	c.SetServer(config.Server)
	return c
}

// Version returns the version of the CLI process that supplied this value at initialization.
func (c *cli) Version() string {
	return c.Configuration.Version
}

// Build returns the build of the CLI process that supplied this value at initialization.
func (c *cli) Build() string {
	return c.Configuration.Build
}

// Server returns the address of the grpc api (host:port) used for the client connection.
func (c *cli) Server() string {
	return c.Configuration.Server
}

// SetServer sets the address of the grpc api (host:port) used for the client connection.
func (c *cli) SetServer(server string) {
	if !strings.Contains(server, ":") {
		server += DefaultPort
	}
	c.Configuration.Server = server
	c.clientConn = nil
}

// SetSkipVerify controls whether a client verifies the
// server's certificate chain and host name.
// If SetSkipVerify is set to true, TLS accepts any certificate
// presented by the server and any host name in that certificate.
// In this mode, TLS is susceptible to man-in-the-middle attacks.
// This should be used only for testing.
func (c *cli) SetSkipVerify(skipVerify bool) {
	c.skipVerify = skipVerify
}

// In returns the reader used for stdin.
func (c *cli) In() *InStream {
	return c.in
}

// Out returns the writer used for stdout.
func (c *cli) Out() *OutStream {
	return c.out
}

// Err returns the writer used for stderr.
func (c *cli) Err() io.Writer {
	return c.err
}

// Console returns the console for formatted printing.
func (c *cli) Console() *Console {
	return c.console
}

// OnInitialize runs initializer functions before executing the command.
func (c *cli) OnInitialize(initializers ...func()) {
	cobra.OnInitialize(initializers...)
}

// Connect opens a connection to the grpc API server (if not already connected).
func (c *cli) Connect() (*grpc.ClientConn, error) {
	if c.clientConn == nil {
		var err error
		c.clientConn, err = NewClientConn(c.Server(), GetToken(c.Server()), c.skipVerify)
		if err != nil {
			// extra newline helpful for grpc errors
			return nil, fmt.Errorf("\nunable to establish grpc connection: %s", err)
		}
	}
	return c.clientConn, nil
}

// ClientConn opens a connection if necessary, then returns the grpc connection to the API.
// If there is an error, this will exit the application with a fatal error.
func (c *cli) ClientConn() *grpc.ClientConn {
	conn, err := c.Connect()
	if err != nil {
		c.Console().Fatalln(err)
	}

	c.clientConn = conn
	return c.clientConn
}

func (c *cli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.SetOutput(c.Err())
	cmd.HelpFunc()(cmd, args)
	return nil
}
