package main

import (
	"os"

	"github.com/appcelerator/amp/cmd/amplifier/server"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/elasticsearch"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/appcelerator/amp/pkg/sms"
	log "github.com/sirupsen/logrus"
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info level or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("amplifier (server version: %s, build: %s)\n", Version, Build)

	// Default Configuration
	cfg = &configuration.Configuration{
		Version:          Version,
		Build:            Build,
		Port:             configuration.DefaultPort,
		H1Port:           configuration.DefaultH1Port,
		EmailSender:      mail.DefaultSender,
		SmsSender:        sms.DefaultSender,
		EtcdEndpoints:    []string{etcd.DefaultEndpoint},
		ElasticsearchURL: elasticsearch.DefaultURL,
		NatsURL:          ns.DefaultURL,
		DockerURL:        docker.DefaultURL,
		DockerVersion:    docker.DefaultVersion,
		Registration:     configuration.RegistrationDefault,
		Notifications:    configuration.NotificationsDefault,
	}

	// Override with configuration file
	if err := configuration.ReadConfig(cfg); err != nil {
		log.Fatalln(err)
	}

	log.Infoln(cfg)
	amplifier, err := server.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(amplifier.Start())
}
