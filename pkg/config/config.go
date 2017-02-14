package amp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	//"github.com/appcelerator/amp/api/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose            bool
	GitHub             string
	Target             string
	AmpAddress         string
	ServerPort         string
	AdminServerPort    string
	WebMailServerPort  string
	CmdTheme           string
	EmailServerAddress string
	EmailServerPort    string
	EmailSender        string
	EmailPwd           string
}

const (
	//DefaultServerAddress amplifier address
	DefaultAmpAddress = "local.appcelerator.io"

	//DefaultServerAddress amplifier address
	DefaultServerPort = "8080"

	//DefaultAdminServerAddress adm-server address
	DefaultAdminServerPort = "8081"

	//DefaultWebServerAddress web address
	DefaultWebMailServerPort = "8082"

	//DefaultEmailServerAddress email server address
	DefaultEmailServerAddress = "smtp.gmail.com"

	//DefaultEmailServerPort email server port
	DefaultEmailServerPort = "587"

	//DefaultEmailSender email sender
	DefaultEmailSender = "amp.appcelerator@gmail.com"
)

// GetRegularConfig return configuration struct load from its regular file (~/.config/amp/amp.yaml)
func GetRegularConfig(verbose bool) *Configuration {
	configFilePath := viper.ConfigFileUsed()
	config := &Configuration{}
	InitConfig(configFilePath, config, verbose, "")
	return config
}

// InitConfig reads in a config file and ENV variables if set.
// Configuration variable lookup occurs in a specific order.
func InitConfig(configFile string, config *Configuration, verbose bool, ampAddr string) {
	config.Verbose = verbose
	config.AmpAddress = ampAddr
	if config.AmpAddress == "" {
		config.AmpAddress = DefaultAmpAddress
	}
	if config.ServerPort == "" {
		config.ServerPort = DefaultServerPort
	}
	if config.AdminServerPort == "" {
		config.AdminServerPort = DefaultAdminServerPort
	}
	if config.WebMailServerPort == "" {
		config.WebMailServerPort = DefaultWebMailServerPort
	}
	if config.EmailServerAddress == "" {
		config.EmailServerAddress = DefaultEmailServerAddress
	}
	if config.EmailServerPort == "" {
		config.EmailServerPort = DefaultEmailServerPort
	}
	if config.EmailSender == "" {
		config.EmailSender = DefaultEmailSender
	}
	// Add matching environment variables - will take precedence over config files.
	viper.AutomaticEnv()

	// Add default config file search paths in order of decreasing precedence.
	viper.SetConfigName("amp")
	viper.AddConfigPath(".")
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		viper.AddConfigPath("$XDG_CONFIG_HOME/amp")
	} else {
		homedir, err := homedir.Dir()
		if err != nil {
			return
		}
		viper.AddConfigPath(path.Join(homedir, ".config/amp"))
	}
	// last place to look: system dir on *nix
	viper.AddConfigPath("/etc/amp/")

	// this must be last: config file specified using --use-config option will take precedence over other paths.
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// If a config file is found, read it in.
	// Extra check for verbose because it might not have been set by
	// a flag, but might be set in the config
	if err := viper.ReadInConfig(); err == nil {
		if verbose || viper.GetBool("Verbose") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if verbose || viper.GetBool("Verbose") {
			if configFile != "" {
				fmt.Printf("Warning: unable to load %s, using default configuration\n", configFile)
			} else {
				fmt.Println("Warning: no valid configuration file (amp.yaml) found in ~/.config/amp/ or current directory")
			}

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
}

// SaveConfiguration saves the configuration to ~/.config/amp/amp.yaml
func SaveConfiguration(c interface{}) (err error) {
	configFilePath := viper.ConfigFileUsed()

	if configFilePath == "" {
		var configdir string
		xdgdir := os.Getenv("XDG_CONFIG_HOME")
		if xdgdir != "" {
			configdir = path.Join(xdgdir, "amp")
		} else {
			homedir, err := homedir.Dir()
			if err != nil {
				return err
			}
			configdir = path.Join(homedir, ".config/amp")
		}
		err = os.MkdirAll(configdir, 0755)
		if err != nil {
			return
		}
		configFilePath = path.Join(configdir, "amp.yaml")
	}

	contents, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(configFilePath, contents, os.ModePerm)
	if err != nil {
		return
	}
	return
}
