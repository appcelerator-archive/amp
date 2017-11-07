package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/cli/cli/compose/convert"
	"github.com/appcelerator/amp/docker/cli/cli/config/configfile"
	"github.com/appcelerator/amp/docker/cli/opts"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type StackStatus struct {
	RunningServices int32
	FailedServices  int32
	TotalServices   int32
	Status          string
}

// StackDeploy deploy a stack
func (d *Docker) StackDeploy(ctx context.Context, stackName string, composeFile []byte, configFile []byte, environment []string) (output string, err error) {
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
			log.Infoln("Using client registry credentials")

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

		opts := stack.DeployOptions{
			Composefile:      compose.Name(),
			Namespace:        stackName,
			ResolveImage:     stack.ResolveImageAlways,
			SendRegistryAuth: true,
			Prune:            false,
			Environment:      environment,
		}
		if err := stack.RunDeploy(cli, opts); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return output, nil
}

// StackList list the stacks
func (d *Docker) StackList(ctx context.Context) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		opts := stack.ListOptions{}
		if err := stack.RunList(cli, opts); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return output, nil
}

// StackRemove remove a stack
func (d *Docker) StackRemove(ctx context.Context, stackName string) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		opts := stack.RemoveOptions{Namespaces: []string{stackName}}
		if err := stack.RunRemove(cli, opts); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return output, nil
}

// StackServices list the services of a stack
func (d *Docker) StackServices(ctx context.Context, stackName string, quiet bool) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		opts := stack.ServicesOptions{
			Quiet:     quiet,
			Filter:    opts.NewFilterOpt(),
			Namespace: stackName,
		}
		if err := stack.RunServices(cli, opts); err != nil {
			return err
		}
		return nil
	}
	if output, err = cliWrapper(cmd); err != nil {
		return "", err
	}
	return output, nil
}

// StackServiceIDs returns service ids of all services in a given stack
func (d *Docker) StackServiceIDs(ctx context.Context, stackName string) ([]string, error) {
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
	services, err := d.StackServiceIDs(ctx, stackName)
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

// waitOnService waits for the service to converge. It outputs a progress bar,
// if appropriate based on the CLI flags.
func (d *Docker) WaitOnStack(ctx context.Context, namespace string, progressWriter io.Writer) []error {
	// List stack services
	options := types.ServiceListOptions{Filters: filters.NewArgs()}
	options.Filters.Add("label", convert.LabelNamespace+"="+namespace)
	services, err := d.client.ServiceList(context.Background(), options)
	if err != nil {
		return []error{err}
	}

	errors := make([]error, len(services))
	var wg sync.WaitGroup
	wg.Add(len(services))

	for i, service := range services {
		go func(serviceID string, err *error) {
			defer wg.Done()
			*err = d.WaitOnService(ctx, serviceID, true)
		}(service.ID, &errors[i])
	}

	// Wait for all services to converge.
	wg.Wait()

	return errors
}
