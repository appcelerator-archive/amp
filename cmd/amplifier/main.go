package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/server"
)

const (
	port = ":50051"
)

var (
	// Version is set with a linker flag (see Makefile)
	Version string
	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	server.Start(port)
}
