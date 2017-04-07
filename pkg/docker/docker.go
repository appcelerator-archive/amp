package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	DefaultURL     = "unix:///var/run/docker.sock"
	DefaultVersion = "1.24"
)

// Docker wrapper
type Docker struct {
	url     string
	version string
	client  *client.Client
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
	if d.client, err = client.NewClient(d.url, d.version, nil, nil); err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", d.url, err)
	}
	return nil
}

// GetClient returns the native docker client
func (d *Docker) GetClient() *client.Client {
	return d.client
}

// DoesServiceExist returns whether the given service exists
func (d *Docker) DoesServiceExist(ctx context.Context, name string) bool {
	list, err := d.client.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil || len(list) == 0 {
		return false
	}
	for _, service := range list {
		if service.Spec.Annotations.Name == name {
			return true
		}
	}
	return false
}
