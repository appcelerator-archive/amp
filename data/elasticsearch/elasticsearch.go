package elasticsearch

import (
	"gopkg.in/olivere/elastic.v3"
)

var (
	// elasticsearch client
	client *elastic.Client
)

// Elasticsearch singleton
type Elasticsearch struct {
}

// Connect to the elastic search server
func (es *Elasticsearch) Connect(url string) error {
	// Create ES client
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	return err
}

// GetClient returns the native elastic search client
func (es *Elasticsearch) GetClient() *elastic.Client {
	return client
}
