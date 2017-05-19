package configuration

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultPort         = ":50101"
	DefaultTimeout      = time.Minute
	RegistrationNone    = "none"
	RegistrationEmail   = "email"
	RegistrationDefault = RegistrationEmail
)

// Configuration is used for amplifier configuration settings
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
	Notifications    bool
}

func (c *Configuration) String() string {
	return fmt.Sprintf("Version: %s\nBuild: %s\nPort: %s\nEtcdEndpoints: %v\nElasticsearchURL: %s\nNatsURL: %s\nDockerURL: %s\nDockerVersion: %s\nRegistration: %s\nNotifications: %v\n", c.Version, c.Build, c.Port, c.EtcdEndpoints, c.ElasticsearchURL, c.NatsURL, c.DockerURL, c.DockerVersion, c.Registration, c.Notifications)
}

// ReadConfig reads the configuration file
func ReadConfig(config *Configuration) error {
	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName("amplifier")
	viper.AddConfigPath("/run/secrets/")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error reading configuration file: %s", err)
	}

	// Unmarshal config into Configuration object
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("Fatal error unmarshalling configuration file: %s", err)
	}

	// Read environment variable
	registration, match := viper.GetString("registration"), false
	for _, valid := range []string{RegistrationNone, RegistrationEmail} {
		if strings.EqualFold(valid, registration) {
			match = true
			break
		}
	}
	if match {
		config.Registration = strings.ToLower(registration)
	} else {
		log.Printf("Invalid registration policy specified: %s, defaulting to: %s\n", registration, RegistrationDefault)
	}
	config.Notifications = viper.GetBool("notifications")

	log.Println("Configuration file successfully loaded")
	return nil
}
