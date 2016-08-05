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
	etc etcd.Etcd
)

func init() {
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
