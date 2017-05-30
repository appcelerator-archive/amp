package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	"github.com/appcelerator/amp/pkg/docker/docker/stack"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/compose/loader"
	types2 "github.com/docker/docker/cli/compose/types"
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

// ComposeParse parses a compose file
func (d *Docker) ComposeParse(ctx context.Context, composeFile []byte) (*types2.Config, error) {
	var details types2.ConfigDetails
	var err error

	// WorkingDir
	details.WorkingDir, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Parsing compose file
	config, err := loader.ParseYAML(composeFile)
	if err != nil {
		return nil, err
	}
	details.ConfigFiles = []types2.ConfigFile{{
		Filename: "filename",
		Config:   config,
	}}

	// Environment
	env := os.Environ()
	details.Environment = make(map[string]string, len(env))
	for _, s := range env {
		if !strings.Contains(s, "=") {
			return nil, fmt.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		details.Environment[kv[0]] = kv[1]
	}

	return loader.Load(details)
}

// ComposeIsAuthorized checks if the given compose file is authorized
func (d *Docker) ComposeIsAuthorized(compose *types2.Config) bool {
	for _, reservedSecret := range constants.Secrets {
		if _, exists := compose.Secrets[reservedSecret]; exists {
			return false
		}
	}

	for _, service := range compose.Services {
		for _, reservedSecret := range constants.Secrets {
			for _, secret := range service.Secrets {
				if strings.EqualFold(secret.Source, reservedSecret) {
					return false
				}
			}
		}
		for _, reservedLabel := range constants.Labels {
			if _, exists := service.Labels[reservedLabel]; exists {
				return false
			}
			if _, exists := service.Deploy.Labels[reservedLabel]; exists {
				return false
			}
		}
		for _, secret := range service.Secrets {
			if strings.EqualFold(secret.Source, constants.SecretCertificate) {
				return false
			}
			if strings.EqualFold(secret.Source, constants.SecretCertificate) {
				return false
			}
		}
		for _, secret := range service.Secrets {
			if strings.EqualFold(secret.Source, constants.SecretCertificate) {
				return false
			}
			if strings.EqualFold(secret.Source, constants.SecretCertificate) {
				return false
			}
		}
	}
	return true
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

// StackRemove remove a stack
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

// TaskList list the tasks
func (d *Docker) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}
	return d.client.TaskList(ctx, options)
}

// NodeList list the nodes
func (d *Docker) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}
	return d.client.NodeList(ctx, options)
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
