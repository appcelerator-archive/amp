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
	InfluxURL        string
	NatsURL          string
}

// String is used to display struct as a string
func (config Config) String() string {
	return fmt.Sprintf("{ Port: %s, EtcdEndpoints: %v, ElasticsearchURL: %s, InfluxURL: %s, NatsURL: %s}", config.Port, config.EtcdEndpoints, config.ElasticsearchURL, config.InfluxURL, config.NatsURL)
}
