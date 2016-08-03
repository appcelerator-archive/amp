package server

import (
	//	"fmt"
	"log"
	"net"

	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/api/rpc/service"
	"google.golang.org/grpc"
)

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
