package runtime

import (
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/pkg/mq"
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

	// MQ is the interface used to access the message queuer backend
	MQ mq.Interface
)
