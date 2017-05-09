package main

import (
	"log"

	"github.com/appcelerator/amp/amp-ui/server/core"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	err := core.ServerInit(Version, Build)
	if err != nil {
		log.Fatal(err)
	}
}
