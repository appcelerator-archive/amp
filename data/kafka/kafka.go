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
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_0_0

	var err error
	client, err = sarama.NewClient([]string{host}, config)
	return err
}

// NewConsumer creates a new consumer
func (kafka *Kafka) NewConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumerFromClient(client)
}
