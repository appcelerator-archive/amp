package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//ControllerConfig Json format of conffile
type ControllerConfig struct {
	haproxyPort       int
	haProxyConffile   string
	dockerWatchPeriod int
	ampStackName      string
	debug             bool
}

var conf ControllerConfig

//Load Json conffile and instanciate new Config
func (config *ControllerConfig) load(version string, build string) {
	config.setDefault()
	config.loadConfigUsingEnvVariable()
	config.display(version, build)
}

//Set default value of configuration
func (config *ControllerConfig) setDefault() {
	config.haproxyPort = 80
	config.debug = false
	config.dockerWatchPeriod = 10
	config.ampStackName = "amp"
}

//Update config with env variables
func (config *ControllerConfig) loadConfigUsingEnvVariable() {
	config.haproxyPort = getIntParameter("PORT", config.haproxyPort)
	config.debug = getBoolParameter("DEBUG", config.debug)
	config.dockerWatchPeriod = getIntParameter("WATCH_PERIOD", config.dockerWatchPeriod)
	config.ampStackName = getStringParameter("AMP_STACK_NAME", config.ampStackName)
}

func (config *ControllerConfig) display(version string, build string) {
	fmt.Printf("HAProxy controller version %s  (build:%s)\n", version, build)
	fmt.Printf("Port: %d\n", config.haproxyPort)
	fmt.Printf("DockerWatchPeriod: %d s\n", config.dockerWatchPeriod)
	fmt.Printf("AMP stack name: %s\n", config.ampStackName)
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
func getBoolParameter(envVariableName string, def bool) bool {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	if strings.ToLower(value) == "true" {
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

//return env variable value, if empty return default value
func getStringArrayParameter(envVariableName string, def []string) []string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	list := strings.Split(strings.Replace(value, " ", "", -1), ",")
	return list
}
