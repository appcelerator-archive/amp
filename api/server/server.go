package server

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/oauth"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stat"
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/kafka"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"google.golang.org/grpc"
)

var (
	// Store is the interface used to access the key/value storage backend
	Store storage.Interface

	// ES is the elasticsearch client
	ES elasticsearch.Elasticsearch

	// Kafka is the kafka client
	Kafka kafka.Kafka

	//Influx is the influxDB client
	Influx influx.Influx
)

// Start starts the server
func Start(config Config) {
	initEtcd(config)
	initElasticsearch(config)
	initKafka(config)
	initInfluxDB(config)

	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("amplifer is unable to listen on: %s\n%v", config.Port[1:], err)
	}
	log.Printf("amplifier is listening on port %s\n", config.Port[1:])
	s := grpc.NewServer()
	// project.RegisterProjectServer(s, &project.Service{})
	logs.RegisterLogsServer(s, &logs.Logs{ES, Store, Kafka})
	stat.RegisterStatServer(s, &stat.Stat{
		Influx: Influx,
	})
	oauth.RegisterGithubServer(s, &oauth.Oauth{
		Store:        Store,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	})
	service.RegisterServiceServer(s, &service.Service{})
	s.Serve(lis)
}

// fail fast on initialization errors; there's no point in attempting
// to continue in a degraded state if there are problems at start up
func initEtcd(config Config) {
	log.Printf("connecting to etcd at %v", strings.Join(config.EtcdEndpoints, ","))
	Store = etcd.New(config.EtcdEndpoints, "amp")
	if err := Store.Connect(5 * time.Second); err != nil {
		panic(err)
	}
	log.Printf("connected to etcd at %v", strings.Join(Store.Endpoints(), ","))
}

// fail fast on initialization errors; there's no point in attempting
// to continue in a degraded state if there are problems at start up
func initElasticsearch(config Config) {
	log.Printf("connecting to elasticsearch at %s\n", config.ElasticsearchURL)
	err := ES.Connect(config.ElasticsearchURL)
	if err != nil {
		log.Panicf("amplifer is unable to connect to elasticsearch on: %s\n%v", config.ElasticsearchURL, err)
	}
	log.Printf("connected to elasticsearch at %s\n", config.ElasticsearchURL)
}

// fail fast on initialization errors; there's no point in attempting
// to continue in a degraded state if there are problems at start up
func initKafka(config Config) {
	log.Printf("connecting to kafka at %s\n", config.KafkaURL)
	err := Kafka.Connect(config.KafkaURL)
	if err != nil {
		log.Panicf("amplifer is unable to connect to kafka on: %s\n%v", config.KafkaURL, err)
	}
	log.Printf("connected to kafka at %s\n", config.KafkaURL)
}

// fail fast on initialization errors; there's no point in attempting
// to continue in a degraded state if there are problems at start up
func initInfluxDB(config Config) {
	log.Printf("connecting to InfluxDB at %s\n", config.InfluxURL)
	Influx = influx.New(config.InfluxURL, "telegraf", "", "")
	if err := Influx.Connect(5 * time.Second); err != nil {
		log.Panicf("amplifer is unable to connect to influxDB on: %s\n%v", config.InfluxURL, err)
	}
	log.Printf("connected to influxDB at %s\n", config.InfluxURL)
}
