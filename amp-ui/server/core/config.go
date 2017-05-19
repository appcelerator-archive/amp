package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//ServerConfig configuration parameters
type ServerConfig struct {
	port          string
	amplifierAddr string
}

//update conf instance with default value and environment variables
func (c *ServerConfig) init(version string, build string) {
	c.setDefault()
	c.loadConfigUsingEnvVariable()
	c.displayConfig(version, build)
}

//Set default value of configuration
func (c *ServerConfig) setDefault() {
	c.port = "3333" //warning 8080 css
	c.amplifierAddr = "amp_amplifier:50101"
}

//Update config with env variables
func (c *ServerConfig) loadConfigUsingEnvVariable() {
	c.port = getStringParameter("SERVER_PORT", c.port)
	c.amplifierAddr = getStringParameter("AMPLIFIER_ADDR", c.amplifierAddr)
}

//display amp-pilot configuration
func (c *ServerConfig) displayConfig(version string, build string) {
	fmt.Printf("amp-ui server version: %v build: %s\n", version, build)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Port: %s\n", c.port)
	fmt.Printf("Amplifier addr: %s\n", c.amplifierAddr)
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

//return env variable value, if empty return default value
func getArrayParameter(envVariableName string, def []string) []string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	list := strings.Split(value, " ")
	for i := range list {
		list[i] = strings.TrimSpace(list[i])
	}
	return list
}

//return env variable value, if empty return default value
func getBoolParameter(envVariableName string, def bool) bool {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	if value == "true" || value == "1" {
		return true
	}
	return false
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
