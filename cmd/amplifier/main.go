package main

import (
	"fmt"
	"log"

	"github.com/appcelerator/amp/cmd/amplifier/server"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/elasticsearch"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/appcelerator/amp/pkg/sms"
)

// build vars
var (
	// Version is set with a linker flag (see Makefile)
	Version string

	// Build is set with a linker flag (see Makefile)
	Build string

	// Amplifier configuration
	config *server.Configuration
)

func main() {
	fmt.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)

	// Default Configuration
	config = &server.Configuration{
		Version:          Version,
		Build:            Build,
		Port:             server.DefaultPort,
		EmailSender:      mail.DefaultSender,
		SmsSender:        sms.DefaultSender,
		EtcdEndpoints:    []string{etcd.DefaultEndpoint},
		ElasticsearchURL: elasticsearch.DefaultURL,
		NatsURL:          ns.DefaultURL,
		DockerURL:        docker.DefaultURL,
		DockerVersion:    docker.DefaultVersion,
	}

	// Override with configuration file
	if err := server.ReadConfig(config); err != nil {
		log.Fatalln(err)
	}
	amplifier := server.New(config)
	amplifier.Start()
}
