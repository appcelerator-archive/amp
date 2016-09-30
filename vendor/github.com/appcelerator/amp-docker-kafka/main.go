package main

import (
	"github.com/appcelerator/amp-docker-kafka/pilot"
	"github.com/appcelerator/amp-docker-kafka/pilot/api/admin"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string
)

const (
	listenAddress = ":4242"
)

func main() {
	log.Printf("amp-docker-kafka (version: %s, build: %s)\n", Version, Build)

	pilot := pilot.New()
	if err := pilot.StartKafka(); err != nil {
		log.Fatalf("Unable to start Kafka: %s\n", err)
	}

	go watchKafka(pilot)

	s := grpc.NewServer()
	admin.RegisterAdminServer(s, &admin.Server{Pilot: pilot})

	// Start listening
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("amp-docker-kafka is unable to listen on: %s\n%v", listenAddress, err)
	}
	log.Printf("amp-docker-kafka is listening on port %s\n", listenAddress)
	if err := s.Serve(lis); err != nil {
		log.Panicf("Problem in api server: %s\n", err)
	}
}

func watchKafka(pilot *pilot.KafkaPilot) {
	if err := pilot.Wait(); err != nil {
		log.Panicf("Problem during Kafka execution: %s\n", err)
	}
	os.Exit(0)
}
