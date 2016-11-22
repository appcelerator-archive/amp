package agentcore

import (
	"fmt"
	"os"
	"strconv"
)

//AgentConfig configuration parameters
type AgentConfig struct {
	dockerEngine string
	apiPort      string
	agentID      string
	grpcPort     string
	serverPort   string
	serverAddr   string
}

var conf AgentConfig

//update conf instance with default value and environment variables
func (c *AgentConfig) init(version string, build string) {
	c.setDefault()
	c.loadConfigUsingEnvVariable()
	c.displayConfig(version, build)
}

//Set default value of configuration
func (c *AgentConfig) setDefault() {
	c.dockerEngine = "unix:///var/run/docker.sock"
	c.apiPort = "3000"
	c.agentID = os.Getenv("HOSTNAME")
	c.serverAddr = "adm-server"
	c.serverPort = "31315"
	c.grpcPort = "31316"
}

//Update config with env variables
func (c *AgentConfig) loadConfigUsingEnvVariable() {
	c.dockerEngine = c.getStringParameter("DOCKER", c.dockerEngine)
	c.apiPort = c.getStringParameter("API_PORT", c.apiPort)
	c.grpcPort = c.getStringParameter("AGENT_PORT", c.grpcPort)
	c.serverPort = c.getStringParameter("SERVER_PORT", c.serverPort)
	c.serverAddr = c.getStringParameter("SERVER_ADDR", c.serverAddr)
}

//display amp-pilot configuration
func (c *AgentConfig) displayConfig(version string, build string) {
	fmt.Printf("adl-agent version: %v build: %s\n", version, build)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Docker-engine: %s\n", c.dockerEngine)
	fmt.Printf("AgentId: %s\n", c.agentID)
	fmt.Printf("Agent port: %s\n", c.grpcPort)
	fmt.Printf("adm-server: %s:%s\n", c.serverAddr, c.serverPort)
	fmt.Println("----------------------------------------------------------------------------")
}

//return env variable value, if empty return default value
func (c *AgentConfig) getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value convert to int, if empty return default value
func (c *AgentConfig) getIntParameter(envVariableName string, def int) int {
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
