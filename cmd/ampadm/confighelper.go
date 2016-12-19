package main

//Duplicated from amp.cli because it can be used outside its package (main): to be relocated in api/client to make it sharable.
import (
	"fmt"
	ampClient "github.com/appcelerator/amp/api/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
	"path"
)

// InitConfig reads in a config file and ENV variables if set.
// Configuration variable lookup occurs in a specific order.
func InitConfig(cli *clusterClient, configFile string, config *ampClient.Configuration, verbose bool, serverAddr string) string {
	config.Verbose = verbose
	config.AdminServerAddress = serverAddr

	// Add matching environment variables - will be first in precedence.
	viper.AutomaticEnv()

	// Add config file specified using flag - will be next in precedence.
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Add default config file (without extension) - will be last in precedence.
	// First search .config/amp directory; if not found, then attempt to also search working
	// directory (will only succeed if process was started from application directory).
	viper.SetConfigName("amp")
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		viper.AddConfigPath("$XDG_CONFIG_HOME/amp")
	}
	viper.AddConfigPath("$HOME/.config/amp")
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	// Extra check for verbose because it might not have been set by
	// a flag, but might be set in the config
	if err := viper.ReadInConfig(); err == nil {
		if verbose || viper.GetBool("Verbose") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if verbose || viper.GetBool("Verbose") {
			fmt.Println("Warning: no valid configuration file (amp.yaml) found in ~/.config/amp/ or current directory")
		}
	}

	// Save viper into config
	err := viper.Unmarshal(config)
	if err != nil {
		fmt.Println(err)
		panic("Unable to process config")
	}

	// check for legacy configuration file for warning
	homedir, err := homedir.Dir()
	legacyConfig := path.Join(homedir, ".amp.yaml")
	if _, err := os.Stat(legacyConfig); err == nil {
		fmt.Printf("Warning: legacy configuration file found (%s)\nIt won't be read, consider moving it to $HOME/.config/amp/amp.yaml or removing it\n", legacyConfig)
	}
	return viper.ConfigFileUsed()
}
