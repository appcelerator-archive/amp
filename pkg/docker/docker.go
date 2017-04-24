package docker

import (
	"fmt"

	"os"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"golang.org/x/net/context"
)

const (
	//DefaultURL docker default URL
	DefaultURL = "unix:///var/run/docker.sock"
	//DefaultVersion docker default version
	DefaultVersion = "1.27"
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

// ContainerCreate creates a container and pulls the image if needed
func (d *Docker) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, name string) (*container.ContainerCreateCreatedBody, error) {
	var (
		namedRef reference.Named
	)

	ref, err := reference.ParseAnyReference(config.Image)
	if err != nil {
		return nil, err
	}
	if named, ok := ref.(reference.Named); ok {
		namedRef = reference.TagNameOnly(named)
	}

	response, err := d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, name)
	if err == nil {
		return &response, nil
	}

	// if image not found try to pull it
	if client.IsErrImageNotFound(err) && namedRef != nil {
		fmt.Fprintf(os.Stderr, "Unable to find image '%s' locally\n", reference.FamiliarString(namedRef))
		if err = d.ImagePull(ctx, config.Image); err != nil {
			return nil, err
		}

		// Retry
		response, err := d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, name)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	return nil, err
}

// PullImage pulls a docker image
func (d *Docker) ImagePull(ctx context.Context, image string) error {
	//ref, err := reference.ParseNormalizedNamed(image)
	//if err != nil {
	//	return err
	//}
	//
	//// Resolve the Repository name from fqn to RepositoryInfo
	//repoInfo, err := registry.ParseRepositoryInfo(ref)
	//if err != nil {
	//	return err
	//}
	//
	//authConfig := command.ResolveAuthConfig(ctx, dockerCli, repoInfo.Index)
	//encodedAuth, err := command.EncodeAuthToBase64(authConfig)
	//if err != nil {
	//	return err
	//}
	//
	//options := types.ImageCreateOptions{
	//	RegistryAuth: encodedAuth,
	//}

	responseBody, err := d.client.ImageCreate(ctx, image, types.ImageCreateOptions{})
	if err != nil {
		return err
	}
	defer responseBody.Close()

	return jsonmessage.DisplayJSONMessagesStream(
		responseBody,
		os.Stdout,
		os.Stdout.Fd(),
		false,
		nil)
}

// ImageRemove remove a docker image from the local repository
func (d *Docker) ImageRemove(ctx context.Context, image string) error {
	_, err := d.client.ImageRemove(ctx, image, types.ImageRemoveOptions{Force: false, PruneChildren: true})
	if err != nil {
		return err
	}
	return nil
}
