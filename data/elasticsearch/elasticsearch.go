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
