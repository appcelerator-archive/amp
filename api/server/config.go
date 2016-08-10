package server

import (
	"fmt"
)

// Config is used for amplifier configuration settings
type Config struct {
	Port          string
	EtcdEndpoints []string
	esURL         string
}

// String is used to display struct as a string
func (config Config) String() string {
	return fmt.Sprintf("{ Port: %s, EtcdEndpoints: %v, esURL: %s}", config.Port, config.EtcdEndpoints, config.esURL)
}
