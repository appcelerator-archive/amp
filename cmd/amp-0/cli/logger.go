package cli

import (
	"log"
)

// Logger is a simple logger for the cli that also implements grpclog.Logger
type Logger struct {
	Verbose bool
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l Logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
	log.Panic()
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l Logger) Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}

// Print delegates to log.Print
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Print(args ...interface{}) {
	if l.Verbose {
		log.Print(args...)
	}
}

// Printf delegates to log.Printf
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, args ...interface{}) {
	if l.Verbose {
		log.Printf(format, args...)
	}
}

// Println delegates to log.Println
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Println(args ...interface{}) {
	if l.Verbose {
		log.Println(args...)
	}
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	log.Panic(v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	log.Panicln(v...)
}

// NewLogger creates a CLI Logger instance
func NewLogger(verbose bool) *Logger {
	return &Logger{Verbose: verbose}
}
