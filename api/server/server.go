package server

import (
	"log"
	"net"
	"strings"
	"time"

	// "github.com/appcelerator/amp/api/rpc/build"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/api/rpc/stats"
	"github.com/appcelerator/amp/api/runtime"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

// Start starts the server
func Start(config Config) {
	// ensure all initialization code fails fast on errors; there is no point in
	// attempting to continue in a degraded state if there are problems at start up
	initEtcd(config)
	initElasticsearch(config)
	initKafka(config)
	initInfluxDB(config)
	initDocker(config)

	// register services
	s := grpc.NewServer()
	// project.RegisterProjectServer(s, &project.Service{})
	logs.RegisterLogsServer(s, &logs.Logs{
		Es:    runtime.Elasticsearch,
		Store: runtime.Store,
		Kafka: runtime.Kafka,
	})
	stats.RegisterStatsServer(s, &stats.Stats{
		Influx: runtime.Influx,
	})
	oauth.RegisterGithubServer(s, &oauth.Oauth{
		Store:        runtime.Store,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	})
	// build.RegisterAmpBuildServer(s, &build.Proxy{})
	service.RegisterServiceServer(s, &service.Service{
		Docker: runtime.Docker,
	})
	stack.RegisterStackServiceServer(s, &stack.Server{
		Store:  runtime.Store,
		Docker: runtime.Docker,
	})

	// start listening
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("amplifer is unable to listen on: %s\n%v", config.Port[1:], err)
	}
	log.Printf("amplifier is listening on port %s\n", config.Port[1:])
	s.Serve(lis)
}

func initEtcd(config Config) {
	log.Printf("connecting to etcd at %v", strings.Join(config.EtcdEndpoints, ","))
	runtime.Store = etcd.New(config.EtcdEndpoints, "amp")
	if err := runtime.Store.Connect(5 * time.Second); err != nil {
		panic(err)
	}
	log.Printf("connected to etcd at %v", strings.Join(runtime.Store.Endpoints(), ","))
}

func initElasticsearch(config Config) {
	log.Printf("connecting to elasticsearch at %s\n", config.ElasticsearchURL)
	err := runtime.Elasticsearch.Connect(config.ElasticsearchURL)
	if err != nil {
		log.Panicf("amplifer is unable to connect to elasticsearch on: %s\n%v", config.ElasticsearchURL, err)
	}
	log.Printf("connected to elasticsearch at %s\n", config.ElasticsearchURL)
}

func initKafka(config Config) {
	log.Printf("connecting to kafka at %s\n", config.KafkaURL)
	err := runtime.Kafka.Connect(config.KafkaURL)
	if err != nil {
		log.Panicf("amplifer is unable to connect to kafka on: %s\n%v", config.KafkaURL, err)
	}
	log.Printf("connected to kafka at %s\n", config.KafkaURL)
}

func initInfluxDB(config Config) {
	log.Printf("connecting to InfluxDB at %s\n", config.InfluxURL)
	runtime.Influx = influx.New(config.InfluxURL, "telegraf", "", "")
	if err := runtime.Influx.Connect(5 * time.Second); err != nil {
		log.Panicf("amplifer is unable to connect to influxDB on: %s\n%v", config.InfluxURL, err)
	}
	log.Printf("connected to influxDB at %s\n", config.InfluxURL)
}

func initDocker(config Config) {
	log.Printf("connecting to Docker API at %s version API: %s\n", config.DockerURL, config.DockerVersion)
	defaultHeaders := map[string]string{"User-Agent": "amplifier-1.0"}
	cli, err := client.NewClient(config.DockerURL, config.DockerVersion, nil, defaultHeaders)
	if err != nil {
		log.Panicf("amplifer is unable to connect to Docker on: %s\n%v", config.DockerURL, err)
	}
	runtime.Docker = cli
	log.Printf("connected to Docker at %s\n", config.DockerURL)
}
