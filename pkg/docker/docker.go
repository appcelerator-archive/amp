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
	env     bool
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
		env:     true,
	}
}

// NewClient instantiates a new Docker wrapper
func NewClient(url string, version string) *Docker {
	return &Docker{
		url:     url,
		version: version,
		env:     false,
	}
}

// Connect to the docker API
func (d *Docker) Connect() (err error) {
	if d.env {
		if d.client, err = docker.NewEnvClient(); err != nil {
			return fmt.Errorf("unable to connect to Docker based on the environment")
		}
	} else {
		if d.url == "" {
			d.url = DefaultURL
		}
		if d.version == "" {
			d.version = DefaultVersion
		}
		if d.client, err = docker.NewClient(d.url, d.version, nil, nil); err != nil {
			return fmt.Errorf("unable to connect to Docker at %s: %v", d.url, err)
		}
	}
	log.Printf("Connected to Docker [%v]\n", d.client)
	return nil
}

// GetClient returns the native docker client
func (d *Docker) GetClient() *docker.Client {
	return d.client
}
