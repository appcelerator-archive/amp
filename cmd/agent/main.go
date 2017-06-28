package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/appcelerator/amp/agent"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info level or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "healthcheck" {
		if !healthy() {
			os.Exit(1)
		}
		os.Exit(0)
	}
	err := core.AgentInit(Version, Build)
	if err != nil {
		log.Fatal(err)
	}
}

func healthy() bool {
	response, err := http.Get("http://127.0.0.1:3000/api/v1/health")
	if err != nil {
		return false
	}
	if response.StatusCode != 200 {
		return false
	}
	return true
}
