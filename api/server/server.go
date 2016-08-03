package server

import (
	//	"fmt"
	"log"
	"net"

	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/api/rpc/service"
	//	"github.com/appcelerator/amp/data/elasticsearch"

	"google.golang.org/grpc"
)

// const (
// 	esIndex   = "amp-project"
// 	esType    = "project"
// 	esMapping = `{
// 			"project":{
// 				"properties":{
// 					"name":{
// 						"type":"string"
// 					}
// 				}
// 			}
// 		}`
// )
//
// var (
// 	es elasticsearch.ElasticSearch
// )

// Start starts the server
func Start(port string) {
	// Initialize elastic search
	// es = elasticsearch.ElasticSearch{}
	// es.Connect()
	// es.CreateIndexIfNotExists(esIndex, esType, esMapping)

	// Start listening
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	project.RegisterProjectServer(s, &projectService{})
	service.RegisterServiceServer(s, &serviceService{})
	s.Serve(lis)
}
