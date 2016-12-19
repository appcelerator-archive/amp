package servercore

import (
	"fmt"
	"os"
	"strconv"
)

//ServerConfig configuration parameters
type ServerConfig struct {
	dockerEngine string
	apiPort      string
	grpcPort     string
}

var conf ServerConfig

//update conf instance with default value and environment variables
func (c *ServerConfig) init(version string, build string) {
	c.setDefault()
	c.loadConfigUsingEnvVariable()
	c.displayConfig(version, build)
}

//Set default value of configuration
func (c *ServerConfig) setDefault() {
	c.dockerEngine = "unix:///var/run/docker.sock"
	c.apiPort = "3000"
	c.grpcPort = "31315"
}

//Update config with env variables
func (c *ServerConfig) loadConfigUsingEnvVariable() {
	c.dockerEngine = c.getStringParameter("DOCKER", c.dockerEngine)
	c.apiPort = c.getStringParameter("API_PORT", c.apiPort)
	c.grpcPort = c.getStringParameter("SERVER_PORT", c.grpcPort)
}

//display amp-pilot configuration
func (c *ServerConfig) displayConfig(version string, build string) {
	fmt.Printf("adm-server version: %v build: %s\n", version, build)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Docker-engine: %s\n", c.dockerEngine)
	fmt.Printf("GRPC Port: %s\n", c.grpcPort)
	fmt.Println("----------------------------------------------------------------------------")
}

//return env variable value, if empty return default value
func (c *ServerConfig) getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value convert to int, if empty return default value
func (c *ServerConfig) getIntParameter(envVariableName string, def int) int {
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
