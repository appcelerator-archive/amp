package amp

import (
	"fmt"
	"github.com/spf13/viper"
)

//AmplifierConfig amplifier config
var AmplifierConfig Config

// Config is used for amplifier configuration settings
type Config struct {
	Version            string
	Port               string
	ServerAddress      string
	EtcdEndpoints      []string
	ElasticsearchURL   string
	ClientID           string
	ClientSecret       string
	InfluxURL          string
	DockerURL          string
	DockerVersion      string
	NatsURL            string
	EmailServerAddress string
	EmailServerPort    string
	EmailSender        string
	EmailPwd           string
	EmailKey           string
	SmsAccountID       string
	SmsSender          string
	SmsKey             string
}

// String is used to display struct as a string
func (config Config) String() string {
	return fmt.Sprintf("{ Port: %s, EtcdEndpoints: %v, ElasticsearchURL: %s, NatsURL: %s, InfluxURL: %s, Docker: %s}", config.Port, config.EtcdEndpoints, config.ElasticsearchURL, config.NatsURL, config.InfluxURL, config.DockerURL)
}

//GetConfig get amplifier config
func GetConfig() *Config {
	return &AmplifierConfig
}

// InitConfig reads secret variables in conffile
func InitConfig(config *Config) {
	// set the parameter(s) which have a default, but we don't want to set using amplifier arguments
	config.EmailSender = DefaultEmailSender

	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName("amplifier")
	viper.AddConfigPath("/.config/amp")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: unable to load /.config/amp/amplifier.yaml\n")
		return
	}

	// Save viper into config
	err := viper.Unmarshal(config)
	if err != nil {
		fmt.Println("Unmarshal amplifier conffile error: %v\n", err)
	}
	fmt.Printf("Amplifier conffile /.config/amp/amplifier.yaml loaded\n")
}
