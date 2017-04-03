package runtime

import (
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/mail"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/client"
)

// Runtime access to initialized clients for various services
var (
	// Docker is the Docker client
	Docker *client.Client

	// Store is the key/value storage client
	Store storage.Interface

	// Elasticsearch is the elasticsearch client
	Elasticsearch elasticsearch.Elasticsearch

	// NatsStreaming is the nats streaming client
	NatsStreaming ns.NatsStreaming

	// Mailer
	Mailer *mail.Mailer
)
