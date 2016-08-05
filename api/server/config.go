package server

import (
	"fmt"
)

// Config is used for amplifier configuration settings
type Config struct {
	Port          string
	EtcdEndpoints []string
}

// String is used to display struct as a string
func (config Config) String() string {
	return fmt.Sprintf("{ Port: %v, EtcdEndpoints: %v }", config.Port, config.EtcdEndpoints)
}
