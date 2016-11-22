package server

import (
	"os"
	"strings"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://etcd:2379"
	elasticsearchDefaultURL = "http://elasticsearch:9200"
	natsDefaultURL          = "nats://nats:4222"
	influxDefaultURL        = "http://influxdb:8086"
	dockerDefaultURL        = "unix:///var/run/docker.sock"
	dockerDefaultVersion    = "1.24"
)

var (
	port             string
	etcdEndpoints    string
	elasticsearchURL string
	natsURL          string
	influxURL        string
	dockerURL        string
	dockerVersion    string
)

// ConfigFromEnv returns configuration from environment
func ConfigFromEnv() Config {
	config := Config{}
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
	return config
}
