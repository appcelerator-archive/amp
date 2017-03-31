package main

import (
	"os"

	"github.com/appcelerator/amp/cmd/ampbeat/beater"
	"github.com/elastic/beats/libbeat/beat"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func main() {
	if beat.Run("ampbeat", Version, beater.New) != nil {
		os.Exit(1)
	}
}
