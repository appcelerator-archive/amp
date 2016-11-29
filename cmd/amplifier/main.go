package main

import (
	"fmt"
	"strings"

	"github.com/appcelerator/amp/api/server"
	"github.com/appcelerator/amp/config"
	flag "github.com/spf13/pflag"
)

const (
	defaultPort         = ":50101"
	defaultClientID     = ""
	defaultClientSecret = ""
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
	clientID         string
	clientSecret     string
	natsURL          string
	influxURL        string
	dockerURL        string
	dockerVersion    string
)

func parseFlags() {
	// set up flags
	flag.StringVarP(&port, "port", "p", defaultPort, "server port (default '"+defaultPort+"')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", amp.EtcdDefaultEndpoint, "etcd comma-separated endpoints")
	flag.StringVarP(&elasticsearchURL, "elasticsearchURL", "s", amp.ElasticsearchDefaultURL, "elasticsearch URL (default '"+amp.ElasticsearchDefaultURL+"')")
	flag.StringVarP(&clientID, "clientid", "i", defaultClientID, "github app clientid (default '"+defaultClientID+"')")
	flag.StringVarP(&clientSecret, "clientsecret", "c", defaultClientSecret, "github app clientsecret (default '"+defaultClientSecret+"')")
	flag.StringVarP(&natsURL, "natsURL", "", amp.NatsDefaultURL, "Nats URL (default '"+amp.NatsDefaultURL+"')")
	flag.StringVarP(&influxURL, "influxURL", "", amp.InfluxDefaultURL, "InfluxDB URL (default '"+amp.InfluxDefaultURL+"')")
	flag.StringVar(&dockerURL, "dockerURL", amp.DockerDefaultURL, "Docker URL (default '"+amp.DockerDefaultURL+"')")

	// parse command line flags
	flag.Parse()

	// update config
	config.Port = port
	config.ClientID = clientID
	config.ClientSecret = clientSecret
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
	config.NatsURL = natsURL
	config.InfluxURL = influxURL
	config.DockerURL = dockerURL
	config.DockerVersion = amp.DockerDefaultVersion
}

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	parseFlags()
	server.Start(config)
}
