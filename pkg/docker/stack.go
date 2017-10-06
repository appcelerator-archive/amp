package docker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/cli/cli/config/configfile"
	"github.com/appcelerator/amp/docker/cli/opts"
	"golang.org/x/net/context"
)

type StackStatus struct {
	RunningServices int32
	FailedServices  int32
	TotalServices   int32
	Status          string
}

// StackDeploy deploy a stack
func (d *Docker) StackDeploy(ctx context.Context, stackName string, composeFile []byte, configFile []byte) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		// Write the compose file to a temporary file
		compose, err := ioutil.TempFile("", stackName)
		if err != nil {
			return err
		}
		defer os.Remove(compose.Name()) // clean up
		if _, err := compose.Write(composeFile); err != nil {
			return err
		}

		if configFile != nil {
			log.Infoln("Using client configuration file")

			// Read client configuration file from reader
			cf := configfile.ConfigFile{
				AuthConfigs: make(map[string]types.AuthConfig),
			}
			if err := json.NewDecoder(bytes.NewReader(configFile)).Decode(&cf); err != nil {
				return err
			}

			// This method is specific to AMP. It updates the cli with the provided configuration.
			cli.SetConfigFile(&cf)
		}

		deployOpt := stack.NewDeployOptions("", compose.Name(), stackName, stack.ResolveImageAlways, true, false)
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
	return d.client.TaskList(ctx, options)
}

// NodeList list the nodes
func (d *Docker) NodeList(ctx context.Context, options types.NodeListOptions) ([]swarm.Node, error) {
	return d.client.NodeList(ctx, options)
}

// StackServices list the services of a stack
func (d *Docker) StackServices(ctx context.Context, stackName string, quietOption bool) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		servicesOpt := stack.NewServicesOptions(quietOption, "", opts.NewFilterOpt(), stackName)
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

// serviceIDs returns service ids of all services in a given stack
func (d *Docker) serviceIDs(ctx context.Context, stackName string) ([]string, error) {
	result, err := d.StackServices(ctx, stackName, true)
	if err != nil {
		return nil, err
	}
	serviceIDs := strings.Fields(result)
	return serviceIDs, nil
}

// StackStatus returns stack status
func (d *Docker) StackStatus(ctx context.Context, stackName string) (*StackStatus, error) {
	var readyServices, failedServices int32
	services, err := d.serviceIDs(ctx, stackName)
	if err != nil {
		return nil, err
	}
	totalServices := int32(len(services))
	for _, service := range services {
		status, err := d.ServiceStatus(ctx, service)
		if err != nil {
			return nil, err
		}
		if status.Status == StateNoMatchingNode || status.Status == StateError {
			failedServices++
		}
		if status.Status == StateRunning {
			readyServices++
		}
	}
	if failedServices != 0 && readyServices != totalServices {
		return &StackStatus{
			RunningServices: readyServices,
			TotalServices:   totalServices,
			FailedServices:  failedServices,
			Status:          StateError,
		}, nil
	}
	if readyServices == totalServices {
		return &StackStatus{
			RunningServices: readyServices,
			TotalServices:   totalServices,
			FailedServices:  failedServices,
			Status:          StateRunning,
		}, nil
	}
	return &StackStatus{
		RunningServices: readyServices,
		TotalServices:   totalServices,
		FailedServices:  failedServices,
		Status:          StateStarting,
	}, nil
}
