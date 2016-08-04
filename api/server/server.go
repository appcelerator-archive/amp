package server

import (
	//	"fmt"
	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/etcd"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	//es elasticsearch.ElasticSearch
	etc etcd.Etcd
)

func init() {
	//// Get elasticsearch url from environment
	//elasticSearchURL := os.Getenv("ES_URL")
	//if elasticSearchURL == "" {
	//	elasticSearchURL = elastic.DefaultURL
	//}
	//fmt.Printf("ES_URL: %v\n", elasticSearchURL)
	//
	//// Initialize elastic search
	//es = elasticsearch.ElasticSearch{}
	//es.Connect(elasticSearchURL)
	//es.CreateIndexIfNotExists(esIndex, esType, esMapping)

	// Initialize etcd
	etc = etcd.Etcd{}
	etc.Connect()
}

// Start starts the server
func Start(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	project.RegisterProjectServer(s, &projectService{})
	service.RegisterServiceServer(s, &serviceService{})
	s.Serve(lis)
}
