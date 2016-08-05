package server

import (
	"log"
	"net"

	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/etcd"
	"google.golang.org/grpc"
)

var (
	etc etcd.Etcd
)

// Start starts the server
func Start(config Config) {
	initEtcd(config)

	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	project.RegisterProjectServer(s, &projectService{})
	service.RegisterServiceServer(s, &serviceService{})
	s.Serve(lis)
}

// fail fast on initialization errors; there's no point in attempting
// to continue in a degraded state if there are problems at start up
func initEtcd(config Config) {
	etc = etcd.Etcd{}
	err := etc.Connect(config.EtcdEndpoints)
	if err != nil {
		panic(err)
	}
}
