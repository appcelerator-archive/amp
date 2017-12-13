package configuration

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DefaultPort          = ":50101"
	DefaultH1Port        = ":5100"
	DefaultTimeout       = time.Minute
	RegistrationNone     = "none"
	RegistrationEmail    = "email"
	RegistrationDefault  = RegistrationEmail
	NotificationsDefault = true
	SecretsDir           = "/run/secrets/"
	ConfigsDir           = "/run/configs/"
	ConfigName           = "amplifier"
	CertificateSecret    = "cert0.pem"
)

// Configuration is used for amplifier configuration settings
type Configuration struct {
	Version          string
	Build            string
	Port             string
	H1Port           string
	EtcdEndpoints    []string
	ElasticsearchURL string
	NatsURL          string
	JWTSecretKey     string `yaml:"JWTSecretKey"`
	SUPassword       string `yaml:"SUPassword"`
	EmailKey         string `yaml:"EmailKey"`
	EmailSender      string
	SmsAccountID     string
	SmsKey           string
	SmsSender        string
	Registration     string
	Notifications    bool
}

func (c *Configuration) String() string {
	return fmt.Sprintf("Version: %s\nBuild: %s\nPort: %s\nH1Port: %s\nEtcdEndpoints: %v\nElasticsearchURL: %s\nNatsURL: %s\nDockerURL: %s\nDockerTLSVerify: %s\nRegistration: %s\nNotifications: %v\n", c.Version, c.Build, c.Port, c.H1Port, c.EtcdEndpoints, c.ElasticsearchURL, c.NatsURL, os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_TLS_VERIFY"), c.Registration, c.Notifications)
}

// ReadConfig reads the configuration file
func ReadConfig(config *Configuration) error {
	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName(ConfigName)
	viper.AddConfigPath(SecretsDir)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error reading configuration file: %s", err)
	}

	// Unmarshal config into Configuration object
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("Fatal error unmarshalling configuration file: %s", err)
	}

	// validate registration environment variable
	registration, match := viper.GetString("registration"), false
	for _, v := range []string{RegistrationNone, RegistrationEmail} {
		if strings.EqualFold(v, registration) {
			match = true
			break
		}
	}
	if match {
		config.Registration = strings.ToLower(registration)
	} else {
		log.Warnf("Invalid registration policy specified: %s, defaulting to: %s\n", registration, RegistrationDefault)
	}
	config.Notifications = viper.GetBool("notifications")

	log.Infoln("Configuration file successfully loaded")
	return nil
}
