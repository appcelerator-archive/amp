package docker

import (
	"fmt"

	"docker.io/go-docker"
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

// NewClient instantiates a new Docker wrapper
func NewClient(url string, version string) *Docker {
	return &Docker{
		url:     url,
		version: version,
	}
}

// Connect to the docker API
func (d *Docker) Connect() (err error) {
	if d.client, err = docker.NewClient(d.url, d.version, nil, nil); err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", d.url, err)
	}
	return nil
}

// GetClient returns the native docker client
func (d *Docker) GetClient() *docker.Client {
	return d.client
}
