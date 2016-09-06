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
	defaultClientID         = ""
	defaultClientSecret     = ""
	kafkaDefaultURL         = "localhost:9092"
	influxDefaultURL        = "http://localhost:8086"
	natsDefaultURL          = "nats://localhost:4222"
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
	kafkaURL         string
	influxURL        string
	natsURL          string
)

func parseFlags() {
	// set up flags
	flag.StringVarP(&port, "port", "p", defaultPort, "server port (default '"+defaultPort+"')")
	flag.StringVarP(&etcdEndpoints, "endpoints", "e", etcdDefaultEndpoints, "etcd comma-separated endpoints")
	flag.StringVarP(&elasticsearchURL, "elasticsearchURL", "s", elasticsearchDefaultURL, "elasticsearch URL (default '"+elasticsearchDefaultURL+"')")
	flag.StringVarP(&clientID, "clientid", "i", defaultClientID, "github app clientid (default '"+defaultClientID+"')")
	flag.StringVarP(&clientSecret, "clientsecret", "c", defaultClientSecret, "github app clientsecret (default '"+defaultClientSecret+"')")
	flag.StringVarP(&kafkaURL, "kafkaURL", "k", kafkaDefaultURL, "kafka URL (default '"+kafkaDefaultURL+"')")
	flag.StringVarP(&influxURL, "influxURL", "", influxDefaultURL, "InfluxDB URL (default '"+influxDefaultURL+"')")
	flag.StringVarP(&natsURL, "natsURL", "", natsDefaultURL, "Nats URL (default '"+natsDefaultURL+"')")

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
	config.KafkaURL = kafkaURL
	config.InfluxURL = influxURL
	config.NatsURL = natsURL
}

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)
	parseFlags()
	server.Start(config)
}
