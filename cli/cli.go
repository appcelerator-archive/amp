package cli

import (
	"io"

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
	ClientConn() (*grpc.ClientConn, error)

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
}

// NewCLI returns a new CLI instance.
func NewCLI(in io.ReadCloser, out, err io.Writer, config *Configuration) Interface {
	c := &cli{
		Configuration: *config,
		in:            NewInStream(in),
		out:           NewOutStream(out),
		err:           err,
	}
	c.console = NewConsole(c.Out(), config.Verbose)
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
	c.Configuration.Server = server
	c.clientConn = nil
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

// ClientConn returns the grpc connection to the API.
func (c *cli) ClientConn() (*grpc.ClientConn, error) {
	if c.clientConn == nil {
		var err error
		c.clientConn, err = NewClientConn(c.Server(), GetToken())
		if err != nil {
			return nil, err
		}
	}
	return c.clientConn, nil
}

func (c *cli) ShowHelp(cmd *cobra.Command, args []string) error {
	cmd.SetOutput(c.Err())
	cmd.HelpFunc()(cmd, args)
	return nil
}
