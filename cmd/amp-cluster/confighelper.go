package main

//Duplicated from amp.cli because it can be used outside its package (main): to be relocated in api/client to make it sharable.
import (
	ampClient "github.com/appcelerator/amp/api/client"
	"github.com/spf13/viper"
)

// InitConfig reads in a config file and ENV variables if set.
// Configuration variable lookup occurs in a specific order.
func InitConfig(cli *clusterClient, configFile string, config *ampClient.Configuration, verbose bool, serverAddr string) string {
	config.Verbose = verbose
	config.ClusterServerAddress = serverAddr

	// Add matching envirionment variables - will be first in precedence.
	viper.AutomaticEnv()

	// Add config file specified using flag - will be next in precedence.
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Add default config file (without extension) - will be last in precedence.
	// First search home directory; if not found, then attempt to also search working
	// directory (will only succeed if process was started from application directory).
	viper.SetConfigName(".amp")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	// Extra check for verbose because it might not have been set by
	// a flag, but might be set in the config
	if err := viper.ReadInConfig(); err == nil {
		if verbose || viper.GetBool("Verbose") {
			cli.printfc(colInfo, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	} else {
		if verbose || viper.GetBool("Verbose") {
			cli.printfc(colWarn, "Warning: no valid configuration file (.amp.yaml) found in home or current directory")
		}
	}

	// Save viper into config
	err := viper.Unmarshal(config)
	if err != nil {
		cli.printfc(colError, "Config format error: %v\n", err)
		//panic("Unable to process config")
	}
	if verbose {
		cli.printfc(colInfo, "Configuration: %+v\n", config)
	}
	return viper.ConfigFileUsed()
}
