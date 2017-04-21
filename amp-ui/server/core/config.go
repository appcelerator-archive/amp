package core

import (
	"fmt"
	"os"
	"strconv"
)

//ServerConfig configuration parameters
type ServerConfig struct {
	port string
}

var conf ServerConfig

//update conf instance with default value and environment variables
func (c *ServerConfig) init(version string, build string) {
	conf.setDefault()
	conf.loadConfigUsingEnvVariable()
	conf.displayConfig(version, build)
}

//Set default value of configuration
func (c *ServerConfig) setDefault() {
	c.port = "4200" //warning 8080 css
}

//Update config with env variables
func (c *ServerConfig) loadConfigUsingEnvVariable() {
	c.port = getStringParameter("SERVER_PORT", c.port)
}

//display amp-pilot configuration
func (c *ServerConfig) displayConfig(version string, build string) {
	fmt.Printf("amp-ui server version: %v build: %s\n", version, build)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Port: %s\n", c.port)
	fmt.Println("----------------------------------------------------------------------------")
}

//return env variable value, if empty return default value
func getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value convert to int, if empty return default value
func getIntParameter(envVariableName string, def int) int {
	value := os.Getenv(envVariableName)
	if value != "" {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return ivalue
	}
	return def
}
