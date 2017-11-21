package prometheus

import (
	"github.com/prometheus/client_golang/api/prometheus"
)

const DefaultURL = "http://prometheus:9090"

// Prometheus wrapper
type Prometheus struct {
	api prometheus.QueryAPI
	url string
}

// NewClient instantiates an Prometheus wrapper
func NewClient(url string) (*Prometheus, error) {
	// Define configuration parameters for a new client
	config := prometheus.Config{
		Address:   DefaultURL,
		Transport: prometheus.DefaultTransport,
	}

	// Return a new Client
	client, err := prometheus.New(config)
	if err != nil {
		return nil, err
	}

	// NewQueryAPI returns a new QueryAPI for the client
	api := prometheus.NewQueryAPI(client)
	return &Prometheus{
		url: url,
		api: api,
	}, nil
}

func (p *Prometheus) Api() prometheus.QueryAPI {
	return p.api
}
