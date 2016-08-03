package server

import (
	//	"fmt"
	"log"
	"net"

	"fmt"
	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/elasticsearch"
	"google.golang.org/grpc"
	"gopkg.in/olivere/elastic.v3"
	"os"
)

var (
	es elasticsearch.ElasticSearch
)

func init() {
	// Get elasticsearch url from environment
	ES_URL := os.Getenv("ES_URL")
	if ES_URL == "" {
		ES_URL = elastic.DefaultURL
	}
	fmt.Printf("ES_URL: %v\n", ES_URL)

	// Initialize elastic search
	es = elasticsearch.ElasticSearch{}
	es.Connect(ES_URL)
	es.CreateIndexIfNotExists(esIndex, esType, esMapping)
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
