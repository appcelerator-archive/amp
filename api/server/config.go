package server

import (
	"fmt"
)

// Config is used for amplifier configuration settings
type Config struct {
	Port             string
	EtcdEndpoints    []string
	ElasticsearchURL string
	ClientID         string
	ClientSecret     string
	KafkaURL         string
	InfluxURL        string
	DockerURL        string
	DockerVersion    string
}

// String is used to display struct as a string
func (config Config) String() string {
	return fmt.Sprintf("{ Port: %s, EtcdEndpoints: %v, ElasticsearchURL: %s, Kafka URL: %s, InfluxURL: %s, Docker: %s}", config.Port, config.EtcdEndpoints, config.ElasticsearchURL, config.KafkaURL, config.InfluxURL, config.DockerURL)
}
