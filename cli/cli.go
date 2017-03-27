package cli

import (
	"io"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// Interface for the CLI's functionality.
type Interface interface {
	In() *InStream
	Out() *OutStream
	Err() io.Writer
	Console() *Console

	ClientConn() *grpc.ClientConn

	OnInitialize(initializers ...func())
}

// c implements cli.Interface
type cli struct {
	in         *InStream
	out        *OutStream
	err        io.Writer
	console    *Console
	clientConn *grpc.ClientConn
}

// NewCLI returns a new CLI instance.
func NewCLI(in io.ReadCloser, out, err io.Writer, verbose bool) Interface {
	c := &cli{
		in:  NewInStream(in),
		out: NewOutStream(out),
		err: err,
	}
	c.console = NewConsole(c.Out(), verbose)
	return c
}

// In returns the reader used for stdin.
func (c cli) In() *InStream {
	return c.in
}

// Out returns the writer used for stdout.
func (c cli) Out() *OutStream {
	return c.out
}

// Err returns the writer used for stderr.
func (c cli) Err() io.Writer {
	return c.err
}

// Console returns the console for formatted printing.
func (c cli) Console() *Console {
	return c.console
}

// OnInitialize runs initializer functions before executing the command.
func (c cli) OnInitialize(initializers ...func()) {
	cobra.OnInitialize(initializers...)
}

// ClientConn returns the grpc connection to the API.
func (c cli) ClientConn() *grpc.ClientConn {
	return c.clientConn
}
