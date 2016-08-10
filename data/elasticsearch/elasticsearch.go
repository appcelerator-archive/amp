package elasticsearch

import (
	"gopkg.in/olivere/elastic.v3"
)

var (
	// elasticSearch client
	client *elastic.Client
)

// ElasticSearch singleton
type ElasticSearch struct {
}

// Connect to the elastic search server
func (es *ElasticSearch) Connect(url string) error {
	// Create ES client
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	return err
}

// CreateIndexIfNotExists Creates an index if it doesn't already exists
func (es *ElasticSearch) CreateIndexIfNotExists(esIndex string, esType string, mapping string) {
	// Use the IndexExists service to check if the index exists
	exists, err := client.IndexExists(esIndex).Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(esIndex).Do()
		if err != nil {
			// TODO: Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// TODO: Handle not acknowledged
			panic("not acked")
		}

		response, err := client.PutMapping().Index(esIndex).Type(esType).BodyString(mapping).Do()
		if err != nil {
			// TODO: Handle error
			panic(err)
		}
		if response == nil {
			// TODO: Handle error
			panic(err)
		}
	}
}

// Index store a document inside elastic search
func (es *ElasticSearch) Index(esIndex string, esType string, body interface{}) {
	// Add a document to the index
	_, err := client.Index().
		Index(esIndex).
		Type("project").
		BodyJson(body).
		Refresh(true).
		Do()
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
}

// All Returns all the documents for a given index
func (es *ElasticSearch) All(esIndex string) []*elastic.SearchHit {
	// Search with a term query
	searchResult, err := client.Search().
		Index(esIndex).
		Do()
	if err != nil {
		// TODO: Handle error
		panic(err)
	}
	return searchResult.Hits.Hits
}

// GetNative returns the native elastic search client
func (es *ElasticSearch) GetClient() *elastic.Client {
	return client
}
