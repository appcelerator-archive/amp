package core

import (
	"fmt"
	"os"
	"strconv"

	"github.com/appcelerator/amp/pkg/config"
)

//AgentConfig configuration parameters
type AgentConfig struct {
	dockerEngine string
	apiPort      string
	period       int
	natsURL      string
	clientID     string
	clusterID    string
}

var conf AgentConfig

//update conf instance with default value and environment variables
func (cfg *AgentConfig) init(version string) {
	cfg.setDefault()
	cfg.loadConfigUsingEnvVariable()
	cfg.displayConfig(version)
}

//Set default value of configuration
func (cfg *AgentConfig) setDefault() {
	cfg.dockerEngine = amp.DockerDefaultURL
	cfg.natsURL = amp.NatsDefaultURL
	cfg.apiPort = "3000"
	cfg.period = 3
	cfg.clientID = "amp-agent-" + os.Getenv("HOSTNAME")
	cfg.clusterID = amp.NatsClusterID
}

//Update config with env variables
func (cfg *AgentConfig) loadConfigUsingEnvVariable() {
	cfg.dockerEngine = getStringParameter("DOCKER", cfg.dockerEngine)
	cfg.natsURL = getStringParameter("NATS_URL", cfg.natsURL)
	cfg.apiPort = getStringParameter("API_PORT", cfg.apiPort)
	cfg.period = getIntParameter("PERIOD", cfg.period)
	cfg.clientID = getStringParameter("CLIENTID", cfg.clientID)
	cfg.clusterID = getStringParameter("CLIENTID", cfg.clusterID)
}

//display amp-pilot configuration
func (cfg *AgentConfig) displayConfig(version string) {
	fmt.Printf("amp-agent version: %v", version)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("Configuration:")
	fmt.Printf("Docker-engine: %s\n", conf.dockerEngine)
	fmt.Printf("Nats URL: %s\n", conf.natsURL)
	fmt.Printf("ClientId: %s\n", conf.clientID)
	fmt.Printf("ClusterId: %s\n", conf.clusterID)
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
