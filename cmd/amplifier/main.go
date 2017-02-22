package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/appcelerator/amp/cmd/amplifier/server"
	"github.com/appcelerator/amp/pkg/config"
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
	config           *amp.Config
	serverAddress    string
	port             string
	etcdEndpoints    string
	elasticsearchURL string
	clientID         string
	clientSecret     string
	natsURL          string
	influxURL        string
	dockerURL        string
)

func parseFlags(config *amp.Config) {
	var displayVersion bool

	// set up flags
	flag.StringVarP(&serverAddress, "server-address", "a", amp.AmplifierDefaultAddress, "Amplifier server address")
	flag.StringVarP(&port, "port", "p", defaultPort, "Server port (default '"+defaultPort+"')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", amp.EtcdDefaultEndpoint, "Etcd comma-separated endpoints")
	flag.StringVarP(&elasticsearchURL, "elasticsearch-url", "s", amp.ElasticsearchDefaultURL, "Elasticsearch URL (default '"+amp.ElasticsearchDefaultURL+"')")
	flag.StringVarP(&clientID, "clientid", "i", defaultClientID, "GitHub app clientid (default '"+defaultClientID+"')")
	flag.StringVarP(&clientSecret, "clientsecret", "c", defaultClientSecret, "GitHub app clientsecret (default '"+defaultClientSecret+"')")
	flag.StringVarP(&natsURL, "nats-url", "", amp.NatsDefaultURL, "Nats URL (default '"+amp.NatsDefaultURL+"')")
	flag.StringVarP(&influxURL, "influx-url", "", amp.InfluxDefaultURL, "InfluxDB URL (default '"+amp.InfluxDefaultURL+"')")
	flag.StringVar(&dockerURL, "docker-url", amp.DockerDefaultURL, "Docker URL (default '"+amp.DockerDefaultURL+"')")
	flag.BoolVarP(&displayVersion, "version", "v", false, "Print version information and quit")

	// parse command line flags
	flag.Parse()

	// Check if command if version
	for _, arg := range os.Args {
		if arg == "version" {
			displayVersion = true
			break
		}
	}

	if displayVersion {
		os.Exit(0)
	}

	// update config
	config.Version = Version
	config.ServerAddress = serverAddress
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
	config = amp.GetConfig()
	amp.InitConfig(config)
	parseFlags(config)
	server.Start(config)
}
