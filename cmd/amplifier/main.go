package main

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/server"
	flag "github.com/spf13/pflag"
)

const (
	defaultPort          = ":50051"
	etcdDefaultEndpoints = "http://localhost:2379"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

// config vars - used for generating a config from command line flags
var (
	config        = server.Config{Port: port}
	port          string
	etcdEndpoints string
)

func parseFlags() {
	// set up flags
	flag.StringVarP(&port, "port", "p", defaultPort, "server port (default ':50051')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", etcdDefaultEndpoints, "etcd comma-separated endpoints")

	// parse command line flags
	flag.Parse()

	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
}

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	parseFlags()
	server.Start(config)
}
