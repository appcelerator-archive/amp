package elasticsearch

import (
	"gopkg.in/olivere/elastic.v3"
	"time"
)

var (
	// elasticsearch client
	client *elastic.Client
)

// Elasticsearch singleton
type Elasticsearch struct {
}

// Connect to the elastic search server
func (es *Elasticsearch) Connect(url string, timeout time.Duration) error {
	// Create ES client
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(true),
		elastic.SetHealthcheckTimeoutStartup(timeout),
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		//elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	return err
}

// GetClient returns the native elastic search client
func (es *Elasticsearch) GetClient() *elastic.Client {
	return client
}

// CreateIndexIfNotExists Creates an index if it doesn't already exists
func (es *Elasticsearch) CreateIndexIfNotExists(esIndex string, esType string, mapping string) error {
	// Use the IndexExists service to check if the index exists
	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(esIndex).Do()
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return err
		}

		response, err := client.PutMapping().Index(esIndex).Type(esType).BodyString(mapping).Do()
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
func (es *Elasticsearch) Index(esIndex string, esType string, body interface{}) error {
	// Add a document to the index
	_, err := client.Index().
		Index(esIndex).
		Type(esType).
		BodyJson(body).
		Refresh(true).
		Do()
	if err != nil {
		return err
	}
	return nil
}
