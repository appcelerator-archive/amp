package main

import (
	"log"

	"github.com/appcelerator/amp/cmd/amplifier/server"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
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
	cfg *configuration.Configuration
)

func main() {
	log.Printf("amplifier (server version: %s, build: %s)\n", Version, Build)

	// Default Configuration
	cfg = &configuration.Configuration{
		Version:          Version,
		Build:            Build,
		Port:             configuration.DefaultPort,
		EmailSender:      mail.DefaultSender,
		SmsSender:        sms.DefaultSender,
		EtcdEndpoints:    []string{etcd.DefaultEndpoint},
		ElasticsearchURL: elasticsearch.DefaultURL,
		NatsURL:          ns.DefaultURL,
		DockerURL:        docker.DefaultURL,
		DockerVersion:    docker.DefaultVersion,
		Registration:     configuration.RegistrationDefault,
		Notifications:    true,
	}

	// Override with configuration file
	if err := configuration.ReadConfig(cfg); err != nil {
		log.Fatalln(err)
	}

	log.Println(cfg)
	amplifier, err := server.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	amplifier.Start()
}
