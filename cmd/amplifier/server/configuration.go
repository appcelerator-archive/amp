package server

import (
	"fmt"
	"log"

	"time"

	"strings"

	"github.com/appcelerator/amp/api/registration"
	"github.com/spf13/viper"
)

const (
	DefaultPort    = ":50101"
	DefaultTimeout = time.Minute
)

// Config is used for amplifier configuration settings
type Configuration struct {
	Version          string
	Build            string
	Port             string
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
	Registration     string
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

	// Read environment variable
	r, match := viper.GetString("registration"), false
	for _, valid := range []string{registration.None, registration.Email} {
		if strings.EqualFold(valid, r) {
			match = true
			break
		}
	}
	if match {
		config.Registration = strings.ToLower(r)
	} else {
		log.Printf("Invalid registration provided: %s, defaulting to: %s\n", r, registration.Default)
	}

	log.Println("Configuration file successfully loaded")
	return nil
}
