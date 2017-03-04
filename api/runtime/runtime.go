package runtime

import (
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/client"
)

var (
	// Store is the interface used to access the key/value storage backend
	Store storage.Interface

	// Elasticsearch is the elasticsearch client
	Elasticsearch elasticsearch.Elasticsearch

	// Influx is the influxDB client
	Influx influx.Influx

	// Docker is the Docker client
	Docker *client.Client

	//NatsStreaming is the nats streaming client
	NatsStreaming ns.NatsStreaming
)
