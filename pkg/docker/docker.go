package docker

import (
	"fmt"
	"os"

	"docker.io/go-docker"
	log "github.com/sirupsen/logrus"
)

// Docker constants
const (
	DefaultURL     = "unix:///var/run/docker.sock"
	DefaultVersion = "1.32"
	MinVersion     = "1.32"
)

// Docker wrapper
type Docker struct {
	client  *docker.Client
	url     string
	version string
}

// NewEnvClient instantiates a new Docker wrapper
func NewEnvClient() *Docker {
	url := os.Getenv("DOCKER_HOST")
	if url == "" {
		url = DefaultURL
	}
	version := os.Getenv("DOCKER_API_VERSION")
	if version == "" {
		version = DefaultVersion
	}
	return &Docker{
		url:     url,
		version: version,
	}
}

// NewClient instantiates a new Docker wrapper
func NewClient(url string, version string) *Docker {
	if err := os.Setenv("DOCKER_HOST", url); err != nil {
		return nil
	}
	if err := os.Setenv("DOCKER_API_VERSION", version); err != nil {
		return nil
	}
	return &Docker{
		url:     url,
		version: version,
	}
}

// Connect to the docker API
func (d *Docker) Connect() (err error) {
	if d.client, err = docker.NewEnvClient(); err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", d.url, err)
	}
	log.Printf("Connected to Docker [%v]\n", d.client)
	return nil
}

// GetClient returns the native docker client
func (d *Docker) GetClient() *docker.Client {
	return d.client
}
