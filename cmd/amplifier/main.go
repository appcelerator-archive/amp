package main

import (
	"fmt"
	"github.com/appcelerator/amp/cmd/amplifier/server"
	"github.com/appcelerator/amp/pkg/config"
	"github.com/prometheus/common/log"
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
		Port:             amp.AmplifierDefaultPort,
		PublicAddress:    amp.AmplifierDefaultPublicAddress,
		EmailSender:      amp.EmailDefaultSender,
		SmsSender:        amp.SmsDefaultSender,
		EtcdEndpoints:    []string{amp.EtcdDefaultEndpoint},
		ElasticsearchURL: amp.ElasticsearchDefaultURL,
		NatsURL:          amp.NatsDefaultURL,
		DockerURL:        amp.DockerDefaultURL,
		DockerVersion:    amp.DockerDefaultVersion,
	}

	// Override with configuration file
	if err := server.ReadConfig(config); err != nil {
		log.Fatalln(err)
	}
	server.Start(config)
}
