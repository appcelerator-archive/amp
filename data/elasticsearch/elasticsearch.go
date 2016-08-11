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
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	return err
}

// Index stores a document inside elastic search
func (es *ElasticSearch) Index(esIndex string, esType string, body interface{}) error {
	// Add a document to the index
	_, err := client.Index().
		Index(esIndex).
		Type("project").
		BodyJson(body).
		Refresh(true).
		Do()
	return err
}

// All returns all the documents for a given index
func (es *ElasticSearch) All(esIndex string) ([]*elastic.SearchHit, error) {
	// Search with a term query
	searchResult, err := client.Search().
		Index(esIndex).
		Do()
	if err != nil {
		return nil, err
	}
	return searchResult.Hits.Hits, nil
}
