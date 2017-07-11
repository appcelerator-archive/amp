package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

const (
	DefaultURL = "http://elasticsearch:9200"
)

// Elasticsearch wrapper
type Elasticsearch struct {
	client    *elastic.Client
	url       string
	timeout   time.Duration
	connected bool
}

// NewClient instantiates an Elasticsearch wrapper
func NewClient(url string, timeout time.Duration) *Elasticsearch {
	return &Elasticsearch{
		url:     url,
		timeout: timeout,
	}
}

// Connect connects to Elasticsearch
func (es *Elasticsearch) Connect() (err error) {
	if es.connected {
		return nil
	}
	es.client, err = elastic.NewClient(
		elastic.SetURL(es.url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckTimeoutStartup(es.timeout),
		elastic.SetMaxRetries(10),
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		//elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		es.connected = false
		return err
	}
	es.connected = true
	return nil
}

// GetClient returns the native elastic search client
func (es *Elasticsearch) GetClient() *elastic.Client {
	return es.client
}

// CreateIndexIfNotExists Creates an index if it doesn't already exists
func (es *Elasticsearch) CreateIndexIfNotExists(ctx context.Context, esIndex string, esType string, mapping string) error {
	// Use the IndexExists service to check if the index exists
	exists, err := es.client.IndexExists(esIndex).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := es.client.CreateIndex(esIndex).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}

		response, err := es.client.PutMapping().Index(esIndex).Type(esType).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}
		if response == nil {
			return err
		}
	}
	return nil
}

// Index store a document inside elastic search
func (es *Elasticsearch) Index(ctx context.Context, esIndex string, esType string, body interface{}) error {
	// Add a document to the index
	_, err := es.client.Index().
		Index(esIndex).
		Type(esType).
		BodyJson(body).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// FormatError formats an elastic.Error
func FormatError(err error) string {
	e, ok := err.(*elastic.Error)
	if !ok {
		return "Unable to cast to elastic.Error"
	}
	if len(e.Details.RootCause) == 0 {
		return e.Error()
	}
	details, err := json.MarshalIndent(e.Details.RootCause[0], "", "    ")
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s\n%s", e.Error(), string(details))
}
