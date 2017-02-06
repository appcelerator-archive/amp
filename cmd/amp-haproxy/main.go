package main

import (
	"github.com/appcelerator/amp/cmd/amp-haproxy/core"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	core.Run(Version, Build)
}
