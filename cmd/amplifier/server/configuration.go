package server

import (
	"fmt"

	"log"

	"github.com/spf13/viper"
)

// Config is used for amplifier configuration settings
type Configuration struct {
	Version          string
	Build            string
	Port             string
	PublicAddress    string
	EtcdEndpoints    []string
	ElasticsearchURL string
	NatsURL          string
	DockerURL        string
	DockerVersion    string
	EmailKey         string
	EmailSender      string
	SmsAccountID     string
	SmsKey           string
	SmsSender        string
}

// ReadConfig reads the configuration file
func ReadConfig(config *Configuration) error {
	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName("amplifier")
	viper.AddConfigPath("/etc/atomiq/")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error reading configuration file: %s", err)
	}

	// Unmarshal config into Configuration object
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("Fatal error unmarshalling configuration file: %s", err)
	}
	log.Println("Configuration file successfully loaded")
	return nil
}
