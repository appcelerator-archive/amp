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
	elasticSearchURL := os.Getenv("ES_URL")
	if elasticSearchURL == "" {
		elasticSearchURL = elastic.DefaultURL
	}
	fmt.Printf("ES_URL: %v\n", elasticSearchURL)

	// Initialize elastic search
	es = elasticsearch.ElasticSearch{}
	es.Connect(elasticSearchURL)
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
