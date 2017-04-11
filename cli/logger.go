package cli

import (
	"fmt"
	"os"

	"google.golang.org/grpc/grpclog"
)

// TODO - needs to be integrated with the CLI output stream (*cli.OutStream)
// to suppress any terminal formatting codes when not attached to a tty;
// then add terminal formatting (color) support for CLI info/warn/success/fail functions.

// Logger is a simple logger for the Atomiq CLI that also implements grpclog.Logger
type Logger struct {
	out     *OutStream
	verbose bool
}

func init() {
	// Creates an instance of Logger for grpc logging
	// WARNING: the grpc logger can only be set during init()
	// https://godoc.org/google.golang.org/grpc/grpclog#SetLogger
	// TODO: set verbose to false after testing
	grpclog.SetLogger(Logger{out: NewOutStream(os.Stdout), verbose: true})
}

// NewLogger creates a CLI Logger instance that writes to the provided stream.
func NewLogger(out *OutStream, verbose bool) *Logger {
	return &Logger{out: out, verbose: verbose}
}

// Verbose returns whether the logger is verbose
func (l Logger) Verbose() bool {
	return l.verbose
}

// OutStream return the underlying output stream
func (l Logger) OutStream() *OutStream {
	return l.out
}

// Fatal is equivalent to fmt.Print() followed by a call to os.Exit(1).
func (l Logger) Fatal(args ...interface{}) {
	l.Print(args)
	os.Exit(1)
}

// Fatalf is equivalent to fmt.Printf() followed by a call to os.Exit(1).
func (l Logger) Fatalf(format string, args ...interface{}) {
	l.Printf(format, args)
	os.Exit(1)
}

// Fatalln is equivalent to fmt.Println() followed by a call to os.Exit(1).
func (l Logger) Fatalln(args ...interface{}) {
	l.Println(args)
	os.Exit(1)
}

// Print is equivalent to fmt.Print() if verbose mode.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Print(args ...interface{}) {
	if l.verbose {
		fmt.Fprint(l.out, args)
	}
}

// Printf is equivalent to fmt.Printf() if verbose mode.
func (l Logger) Printf(format string, args ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.out, format, args)
	}
}

// Println is equivalent to fmt.Println() if verbose mode.
func (l Logger) Println(args ...interface{}) {
	if l.verbose {
		fmt.Fprintln(l.out, args)
	}
}
