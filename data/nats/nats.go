package messaging

import (
	"github.com/nats-io/go-nats-streaming"
	"log"
)

// Nats singleton
type Nats struct {
	clusterID string
	clientID  string
	async     bool
	URL       string
}

var (
	sc stan.Conn
)

// Connect to Nats
func (Nats *Nats) Connect() error {
	var err error
	Nats.clusterID = "test-cluster"
	Nats.clientID = "amp-agent"
	Nats.URL = stan.DefaultNatsURL
	sc, err = stan.Connect(Nats.clusterID, Nats.clientID, stan.NatsURL(Nats.URL))
	return err
}

func ackHandler(ackedNuid string, err error) {
	if err != nil {
		log.Printf("Warning: error publishing msg id %s: %v\n", ackedNuid, err.Error())
	}
}

// Publish a message asynchronously
func (Nats *Nats) Publish(subject string, data []byte) (string, error) {
	return sc.PublishAsync(subject, data, ackHandler)
}

// Close nats
func (Nats *Nats) Close() error {
	return sc.Close()
}
