package main

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/server"
	flag "github.com/spf13/pflag"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://localhost:2379"
	elasticsearchDefaultURL = "http://localhost:9200"
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
	config           server.Config
	port             string
	etcdEndpoints    string
	elasticsearchURL string
)

func parseFlags() {
	// set up flags
	flag.StringVarP(&port, "port", "p", defaultPort, "server port (default '"+defaultPort+"')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", etcdDefaultEndpoints, "etcd comma-separated endpoints")
	flag.StringVarP(&elasticsearchURL, "elasticsearchURL", "s", elasticsearchDefaultURL, "elasticsearch URL (default '"+elasticsearchDefaultURL+"')")

	// parse command line flags
	flag.Parse()

	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
}

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	parseFlags()
	server.Start(config)
}
