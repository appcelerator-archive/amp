package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/server"
	flag "github.com/spf13/pflag"
	"strings"
)

const (
	defaultPort         = ":50101"
	defaultClientID     = ""
	defaultClientSecret = ""
)

var (
	etcdDefaultEndpoints    = "http://127.0.0.1:2379"
	elasticsearchDefaultURL = "http://127.0.0.1:9200"
	natsDefaultURL          = "nats://127.0.0.1:4222"
	influxDefaultURL        = "http://127.0.0.1:8086"
	dockerDefaultURL        = "unix:///var/run/docker.sock"
	dockerDefaultVersion    = "1.24"
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
	isService        bool
)

func parseFlags() {
	// set up flags
	flag.BoolVar(&isService, "service", false, "Launched as a service into swarm network")
	flag.StringVarP(&port, "port", "p", defaultPort, "server port (default '"+defaultPort+"')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", etcdDefaultEndpoints, "etcd comma-separated endpoints")
	flag.StringVarP(&elasticsearchURL, "elasticsearchURL", "s", elasticsearchDefaultURL, "elasticsearch URL (default '"+elasticsearchDefaultURL+"')")
	flag.StringVarP(&clientID, "clientid", "i", defaultClientID, "github app clientid (default '"+defaultClientID+"')")
	flag.StringVarP(&clientSecret, "clientsecret", "c", defaultClientSecret, "github app clientsecret (default '"+defaultClientSecret+"')")
	flag.StringVarP(&natsURL, "natsURL", "", natsDefaultURL, "Nats URL (default '"+natsDefaultURL+"')")
	flag.StringVarP(&influxURL, "influxURL", "", influxDefaultURL, "InfluxDB URL (default '"+influxDefaultURL+"')")
	flag.StringVar(&dockerURL, "dockerURL", dockerDefaultURL, "Docker URL (default '"+dockerDefaultURL+"')")

	// parse command line flags
	flag.Parse()

	//Update url is service usage
	if isService {
		etcdEndpoints = "http://etcd:2379"
		elasticsearchURL = "http://elasticsearch:9200"
		natsURL = "nats://nats:4222"
		influxURL = "http://influxdb:8086"
	}

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
	config.DockerVersion = dockerDefaultVersion
}

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	parseFlags()
	server.Start(config)
}
