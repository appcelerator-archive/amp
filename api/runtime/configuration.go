package runtime

import (
	"github.com/appcelerator/amp/data/elasticsearch"
	"github.com/appcelerator/amp/data/influx"
	"github.com/appcelerator/amp/data/kafka"
	"github.com/appcelerator/amp/data/storage"
)

var (
	// Store is the interface used to access the key/value storage backend
	Store storage.Interface

	// Elasticsearch is the elasticsearch client
	Elasticsearch elasticsearch.Elasticsearch

	// Kafka is the kafka client
	Kafka kafka.Kafka

	// Influx is the influxDB client
	Influx influx.Influx
)
