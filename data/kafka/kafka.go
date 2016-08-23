package kafka

import (
	"github.com/Shopify/sarama"
)

var (
	client sarama.Client
)

// Kafka singleton
type Kafka struct {
}

// Connect to kafka
func (kafka *Kafka) Connect(host string) error {
	// Create Kafka client
	var err error
	client, err = sarama.NewClient([]string{host}, nil)
	return err
}

// NewConsumer creates a new consumer
func (kafka *Kafka) NewConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumerFromClient(client)
}
