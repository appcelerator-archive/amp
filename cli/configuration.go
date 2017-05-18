package cli

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration is for all configurable client settings
type Configuration struct {
	Version string
	Build   string
	Server  string
	Verbose bool
	Theme   string
}

// ReadClientConfig reads the CLI configuration
func ReadClientConfig(config *Configuration) error {
	viper.AutomaticEnv() // Ensure env variables take precedence over config settings

	viper.SetConfigName("amp")               // Set CLI config file name
	viper.AddConfigPath("./.amp")            // Path to store config file in current working directory
	viper.AddConfigPath("$HOME/.config/amp") // Default path to store config file

	// Find and read the config file
	if err := viper.ReadInConfig(); err == nil {
		// Unmarshal the config file into CLI Configuration object
		if err := viper.Unmarshal(config); err != nil {
			return fmt.Errorf("error unmarshalling config file: %s", err)
		}
	}
	return nil
}
