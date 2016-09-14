package server

import (
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	defaultPort             = ":50101"
	etcdDefaultEndpoints    = "http://localhost:2379"
	serverAddress           = "localhost" + defaultPort
	elasticsearchDefaultURL = "http://localhost:9200"
	kafkaDefaultURL         = "localhost:9092"
	influxDefaultURL        = "http://localhost:8086"
)

var (
	config           Config
	port             string
	etcdEndpoints    string
	elasticsearchURL string
	kafkaURL         string
	influxURL        string
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
	kafkaURL = os.Getenv("kafkaURL")
	if kafkaURL == "" {
		kafkaURL = kafkaDefaultURL
	}
	influxURL = os.Getenv("influxURL")
	if influxURL == "" {
		influxURL = influxDefaultURL
	}

	// update config
	config.Port = port
	for _, s := range strings.Split(etcdEndpoints, ",") {
		config.EtcdEndpoints = append(config.EtcdEndpoints, s)
	}
	config.ElasticsearchURL = elasticsearchURL
	config.KafkaURL = kafkaURL
	config.InfluxURL = influxURL
}

// StartTestServer start a server for test
func StartTestServer() (Config, *grpc.ClientConn) {
	parseEnv()

	go Start(config)

	// there is no event when the server starts listening, so we just wait a second
	time.Sleep(1 * time.Second)
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Println("connection failure")
		os.Exit(1)
	}
	return config, conn
}
