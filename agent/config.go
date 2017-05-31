package core

import (
	"log"
	"os"
	"strconv"

	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/nats-streaming"
)

// AgentConfig configuration parameters
type AgentConfig struct {
	dockerEngine     string
	apiPort          string
	period           int
	natsURL          string
	clientID         string
	clusterID        string
	metricsPeriod    int
	logsBufferSize   int
	logsBufferPeriod int
}

var conf AgentConfig

// Update conf instance with default value and environment variables
func (cfg *AgentConfig) init(version, build string) {
	cfg.setDefault()
	cfg.loadConfigUsingEnvVariable()
	cfg.controlConfig()
	cfg.displayConfig(version, build)
}

// Set default value of configuration
func (cfg *AgentConfig) setDefault() {
	cfg.dockerEngine = docker.DefaultURL
	cfg.natsURL = ns.DefaultURL
	cfg.apiPort = "3000"
	cfg.period = 3
	cfg.clientID = "agent-" + os.Getenv("HOSTNAME")
	cfg.clusterID = ns.ClusterID
	cfg.metricsPeriod = 30
	cfg.logsBufferSize = 0
	cfg.logsBufferPeriod = 0
}

// Update config with env variables
func (cfg *AgentConfig) loadConfigUsingEnvVariable() {
	cfg.dockerEngine = getenv("DOCKER", cfg.dockerEngine)
	cfg.natsURL = getenv("NATS_URL", cfg.natsURL)
	cfg.apiPort = getenv("API_PORT", cfg.apiPort)
	cfg.period = getenvi("PERIOD", cfg.period)
	cfg.clientID = getenv("CLIENTID", cfg.clientID)
	cfg.clusterID = getenv("CLIENTID", cfg.clusterID)
	cfg.metricsPeriod = getenvi("METRICS_PERIOD", cfg.metricsPeriod)
	cfg.logsBufferSize = getenvi("LOGS_BUFFER_SIZE", cfg.logsBufferSize)
	cfg.logsBufferPeriod = getenvi("LOGS_BUFFER_PERIOD", cfg.logsBufferPeriod)
}

func (cfg *AgentConfig) controlConfig() {
	if cfg.metricsPeriod < 0 {
		cfg.metricsPeriod = 3 //Min 3 seconds
	}
	if cfg.logsBufferPeriod < 0 {
		cfg.logsBufferPeriod = 0
	}
	if cfg.logsBufferSize < 0 {
		cfg.logsBufferSize = 0
	}
}

// Display agent version and configuration information
func (cfg *AgentConfig) displayConfig(version, build string) {
	log.Printf("agent version: %s, build: %s", version, build)
	log.Println("----------------------------------------------------------------------------")
	log.Println("Configuration:")
	log.Printf("Docker-engine: %s\n", conf.dockerEngine)
	log.Printf("Nats URL: %s\n", conf.natsURL)
	log.Printf("ClientId: %s\n", conf.clientID)
	log.Printf("ClusterId: %s\n", conf.clusterID)
	log.Printf("MetricsPeriod: %d second(s)\n", conf.metricsPeriod)
	log.Printf("LogsBufferSize: %d\n", conf.logsBufferSize)
	log.Printf("LogsBufferPeriod: %d second(s)\n", conf.logsBufferPeriod)
	log.Println("----------------------------------------------------------------------------")
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, or the supplied default value if the variable is not present.
func getenv(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

// Getenv retrieves the value of the environment variable named by the key,
// converted to an int.
// It returns the value, or the supplied default value if the variable is not present.
func getenvi(envVariableName string, def int) int {
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
