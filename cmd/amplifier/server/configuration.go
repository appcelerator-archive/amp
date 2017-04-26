package server

import (
	"fmt"
	"log"

	"time"

	"github.com/spf13/viper"
)

const (
	DefaultPort       = ":50101"
	DefaultPublicHost = "127.0.0.1"
	DefaultTimeout    = time.Minute
)

// Config is used for amplifier configuration settings
type Configuration struct {
	Version          string
	Build            string
	Port             string
	PublicHost       string
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

func (c *Configuration) String() string {
	s := fmt.Sprintf("Version: %s\n", c.Version)
	s += fmt.Sprintf("Build: %s\n", c.Build)
	s += fmt.Sprintf("Port: %s\n", c.Port)
	s += fmt.Sprintf("PublicHost: %s\n", c.PublicHost)
	s += fmt.Sprintf("EtcdEndpoints: %s\n", c.EtcdEndpoints)
	s += fmt.Sprintf("ElasticsearchURL: %s\n", c.ElasticsearchURL)
	s += fmt.Sprintf("NatsURL: %s\n", c.NatsURL)
	s += fmt.Sprintf("DockerURL: %s\n", c.DockerURL)
	s += fmt.Sprintf("DockerVersion: %s\n", c.DockerVersion)
	return s
}

// ReadConfig reads the configuration file
func ReadConfig(config *Configuration) error {
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

	// Override with environment
	viper.SetEnvPrefix("amp")
	viper.AutomaticEnv()
	config.PublicHost = viper.GetString("PublicHost")

	log.Println("Configuration file successfully loaded")
	return nil
}
