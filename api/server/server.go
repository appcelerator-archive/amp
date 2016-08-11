package server

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"google.golang.org/grpc"
)

var (
	// Store is the interface used to access the key/value storage backend
	Store storage.Interface

	// es is the elasticsearch
	ES elasticsearch.Elasticsearch
)

// Start starts the server
func Start(config Config) {
	initEtcd(config)
	initElasticsearch(config)

	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("amplifer is unable to listen on: %s\n%v", config.Port[1:], err)
	}
	log.Printf("amplifier is listening on port %s\n", config.Port[1:])
	s := grpc.NewServer()
	// project.RegisterProjectServer(s, &project.Service{})
	logs.RegisterLogsServer(s, &logs.Logs{})
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
		log.Fatalf("amplifer is unable to connect to elasticsearch on: %s\n%v", config.ElasticsearchURL, err)
	}
	log.Printf("connected to elasticsearch at %s\n", config.ElasticsearchURL)
}
