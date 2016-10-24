package server

import (
	"os"
	"strings"

	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://127.0.0.1:2379"
	serverAddress           = "127.0.0.1" + defaultPort
	elasticsearchDefaultURL = "http://127.0.0.1:9200"
	natsDefaultURL          = "nats://127.0.0.1:4222"
	influxDefaultURL        = "http://127.0.0.1:8086"
	dockerDefaultURL        = "unix:///var/run/docker.sock"
	dockerDefaultVersion    = "1.24"
)

var (
	config           Config
	port             string
	etcdEndpoints    string
	elasticsearchURL string
	natsURL          string
	influxURL        string
	dockerURL        string
	dockerVersion    string
)

func parseEnv() {
	port = os.Getenv("port")
	if port == "" {
		port = defaultPort
	}
	etcdEndpoints = os.Getenv("endpoints")
	if etcdEndpoints == "" {
		etcdEndpoints = etcdDefaultEndpoints
	}
	elasticsearchURL = os.Getenv("elasticsearchURL")
	if elasticsearchURL == "" {
		elasticsearchURL = elasticsearchDefaultURL
	}
	natsURL = os.Getenv("natsURL")
	if natsURL == "" {
		natsURL = natsDefaultURL
	}
	influxURL = os.Getenv("influxURL")
	if influxURL == "" {
		influxURL = influxDefaultURL
	}
	dockerURL = os.Getenv("DOCKER_HOST")
	if dockerURL == "" {
		dockerURL = dockerDefaultURL
	}
	dockerVersion = os.Getenv("DOCKER_API_VERSION")
	if dockerVersion == "" {
		dockerVersion = dockerDefaultVersion
	}
	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
	config.NatsURL = natsURL
	config.InfluxURL = influxURL
	config.DockerURL = dockerURL
	config.DockerVersion = dockerVersion
}

// StartTestServer start a server for test
func StartTestServer() (Config, *grpc.ClientConn) {
	parseEnv()

	go Start(config)

	// Wait for swarm to be ready
	log.Println("Waiting for swarm to be ready")
	if err := initDependencies(config); err != nil {
		log.Panicln("Dependencies are not ready", err)
	}

	// Connect to amplifier
	log.Println("Connecting to amplifier")
	conn, err := grpc.Dial(serverAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second))
	if err != nil {
		log.Panicln("Cannot connect to amplifier", err)
	}
	log.Println("Connected to amplifier")

	return config, conn
}
