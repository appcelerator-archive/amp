package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/appcelerator/amp/pkg/docker/docker/stack"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/flags"
	"github.com/docker/docker/client"
	dopts "github.com/docker/docker/opts"
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
	client    *client.Client
	connected bool
	url       string
	version   string
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

// ContainerCreate creates a container and pulls the image if needed
func (d *Docker) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, name string) (*container.ContainerCreateCreatedBody, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

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

// ImagePull pulls a docker image
func (d *Docker) ImagePull(ctx context.Context, image string) error {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return err
		}
	}
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

// StackDeploy deploy a stack
func (d *Docker) StackDeploy(ctx context.Context, stackName string, composeFile []byte) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		// Write the compose file to a temporary file
		tmp, err := ioutil.TempFile("", stackName)
		if err != nil {
			return err
		}
		defer os.Remove(tmp.Name()) // clean up
		if _, err := tmp.Write(composeFile); err != nil {
			return err
		}
		deployOpt := stack.NewDeployOptions(stackName, tmp.Name(), true)
		if err := stack.RunDeploy(cli, deployOpt); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return string(output), nil
}

// StackList list the stacks
func (d *Docker) StackList(ctx context.Context) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		listOpt := stack.NewListOptions()
		if err := stack.RunList(cli, listOpt); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return string(output), nil
}

// StackRemove remoce a stack
func (d *Docker) StackRemove(ctx context.Context, stackName string) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		rmOpt := stack.NewRemoveOptions([]string{stackName})
		if err := stack.RunRemove(cli, rmOpt); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return string(output), nil
}

func (d *Docker) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}
	return d.client.TaskList(ctx, options)
}

// StackServices list the services of a stack
func (d *Docker) StackServices(ctx context.Context, stackName string) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		servicesOpt := stack.NewServicesOptions(false, "", dopts.NewFilterOpt(), stackName)
		if err := stack.RunServices(cli, servicesOpt); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return string(output), nil
}

/*
func (d *Docker)  NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	cmd := func(cli *command.DockerCli) error {
		//servicesOpt := stack.NewServicesOptions(false, "", dopts.NewFilterOpt(), stackName)
		if err := stack.RunServices(cli, servicesOpt); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return string(output), nil
}
*/

func cliWrapper(cmd func(cli *command.DockerCli) error) (string, error) {
	r, w, _ := os.Pipe()
	cli := command.NewDockerCli(os.Stdin, w, w)
	if err := cli.Initialize(flags.NewClientOptions()); err != nil {
		return "", err
	}
	if err := cmd(cli); err != nil {
		return "", err
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	return string(outs), nil
}
