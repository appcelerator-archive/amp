package stack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/compose/convert"
	"github.com/docker/cli/cli/compose/loader"
	composetypes "github.com/docker/cli/cli/compose/types"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	apiclient "github.com/docker/docker/client"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func deployCompose(ctx context.Context, dockerCli command.Cli, opts DeployOptions) error {
	configDetails, err := getConfigDetails(opts.Composefile)
	if err != nil {
		return err
	}

	config, err := loader.Load(configDetails)
	if err != nil {
		if fpe, ok := err.(*loader.ForbiddenPropertiesError); ok {
			return errors.Errorf("Compose file contains unsupported options:\n\n%s\n",
				propertyWarnings(fpe.Properties))
		}

		return err
	}

	unsupportedProperties := loader.GetUnsupportedProperties(configDetails)
	if len(unsupportedProperties) > 0 {
		fmt.Fprintf(dockerCli.Err(), "Ignoring unsupported options: %s\n\n",
			strings.Join(unsupportedProperties, ", "))
	}

	deprecatedProperties := loader.GetDeprecatedProperties(configDetails)
	if len(deprecatedProperties) > 0 {
		fmt.Fprintf(dockerCli.Err(), "Ignoring deprecated options:\n\n%s\n\n",
			propertyWarnings(deprecatedProperties))
	}

	if err := checkDaemonIsSwarmManager(ctx, dockerCli); err != nil {
		return err
	}

	namespace := convert.NewNamespace(opts.Namespace)

	if opts.Prune {
		services := map[string]struct{}{}
		for _, service := range config.Services {
			services[service.Name] = struct{}{}
		}
		pruneServices(ctx, dockerCli, namespace, services)
	}

	serviceNetworks := getServicesDeclaredNetworks(config.Services)
	networks, externalNetworks := convert.Networks(namespace, config.Networks, serviceNetworks)
	if err := validateExternalNetworks(ctx, dockerCli.Client(), externalNetworks); err != nil {
		return err
	}
	if err := createNetworks(ctx, dockerCli, namespace, networks); err != nil {
		return err
	}

	secrets, err := convert.Secrets(namespace, config.Secrets)
	if err != nil {
		return err
	}
	if err := createSecrets(ctx, dockerCli, secrets); err != nil {
		return err
	}

	configs, err := convert.Configs(namespace, config.Configs)
	if err != nil {
		return err
	}
	if err := createConfigs(ctx, dockerCli, configs); err != nil {
		return err
	}

	services, err := convert.Services(namespace, config, dockerCli.Client())
	if err != nil {
		return err
	}
	return deployServices(ctx, dockerCli, services, namespace, opts.SendRegistryAuth, opts.ResolveImage, opts.ExpectedState)
}

func getServicesDeclaredNetworks(serviceConfigs []composetypes.ServiceConfig) map[string]struct{} {
	serviceNetworks := map[string]struct{}{}
	for _, serviceConfig := range serviceConfigs {
		if len(serviceConfig.Networks) == 0 {
			serviceNetworks["default"] = struct{}{}
			continue
		}
		for network := range serviceConfig.Networks {
			serviceNetworks[network] = struct{}{}
		}
	}
	return serviceNetworks
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}

func getConfigDetails(composefile string) (composetypes.ConfigDetails, error) {
	var details composetypes.ConfigDetails

	absPath, err := filepath.Abs(composefile)
	if err != nil {
		return details, err
	}
	details.WorkingDir = filepath.Dir(absPath)

	configFile, err := getConfigFile(composefile)
	if err != nil {
		return details, err
	}
	// TODO: support multiple files
	details.ConfigFiles = []composetypes.ConfigFile{*configFile}
	details.Environment, err = buildEnvironment(os.Environ())
	return details, err
}

func buildEnvironment(env []string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, errors.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func getConfigFile(filename string) (*composetypes.ConfigFile, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := loader.ParseYAML(bytes)
	if err != nil {
		return nil, err
	}
	return &composetypes.ConfigFile{
		Filename: filename,
		Config:   config,
	}, nil
}

func validateExternalNetworks(
	ctx context.Context,
	client dockerclient.NetworkAPIClient,
	externalNetworks []string,
) error {
	for _, networkName := range externalNetworks {
		network, err := client.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{})
		switch {
		case dockerclient.IsErrNotFound(err):
			return errors.Errorf("network %q is declared as external, but could not be found. You need to create a swarm-scoped network before the stack is deployed", networkName)
		case err != nil:
			return err
		case container.NetworkMode(networkName).IsUserDefined() && network.Scope != "swarm":
			return errors.Errorf("network %q is declared as external, but it is not in the right scope: %q instead of \"swarm\"", networkName, network.Scope)
		}
	}
	return nil
}

func createSecrets(
	ctx context.Context,
	dockerCli command.Cli,
	secrets []swarm.SecretSpec,
) error {
	client := dockerCli.Client()

	for _, secretSpec := range secrets {
		secret, _, err := client.SecretInspectWithRaw(ctx, secretSpec.Name)
		switch {
		case err == nil:
			// secret already exists, then we update that
			if err := client.SecretUpdate(ctx, secret.ID, secret.Meta.Version, secretSpec); err != nil {
				return errors.Wrapf(err, "failed to update secret %s", secretSpec.Name)
			}
		case apiclient.IsErrSecretNotFound(err):
			// secret does not exist, then we create a new one.
			if _, err := client.SecretCreate(ctx, secretSpec); err != nil {
				return errors.Wrapf(err, "failed to create secret %s", secretSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func createConfigs(
	ctx context.Context,
	dockerCli command.Cli,
	configs []swarm.ConfigSpec,
) error {
	client := dockerCli.Client()

	for _, configSpec := range configs {
		config, _, err := client.ConfigInspectWithRaw(ctx, configSpec.Name)
		switch {
		case err == nil:
			// config already exists, then we update that
			if err := client.ConfigUpdate(ctx, config.ID, config.Meta.Version, configSpec); err != nil {
				errors.Wrapf(err, "failed to update config %s", configSpec.Name)
			}
		case apiclient.IsErrConfigNotFound(err):
			// config does not exist, then we create a new one.
			if _, err := client.ConfigCreate(ctx, configSpec); err != nil {
				errors.Wrapf(err, "failed to create config %s", configSpec.Name)
			}
		default:
			return err
		}
	}
	return nil
}

func createNetworks(
	ctx context.Context,
	dockerCli command.Cli,
	namespace convert.Namespace,
	networks map[string]types.NetworkCreate,
) error {
	client := dockerCli.Client()

	existingNetworks, err := getStackNetworks(ctx, client, namespace.Name())
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	for internalName, createOpts := range networks {
		name := namespace.Scope(internalName)
		if _, exists := existingNetworkMap[name]; exists {
			continue
		}

		if createOpts.Driver == "" {
			createOpts.Driver = DefaultNetworkDriver
		}

		fmt.Fprintf(dockerCli.Out(), "Creating network %s\n", name)
		if _, err := client.NetworkCreate(ctx, name, createOpts); err != nil {
			return errors.Wrapf(err, "failed to create network %s", internalName)
		}
	}
	return nil
}

func deployServices(
	ctx context.Context,
	dockerCli command.Cli,
	services map[string]swarm.ServiceSpec,
	namespace convert.Namespace,
	sendAuth bool,
	resolveImage string,
	expectedState swarm.TaskState,
) error {
	apiClient := dockerCli.Client()
	out := dockerCli.Out()

	existingServices, err := getServices(ctx, apiClient, namespace.Name())
	if err != nil {
		return err
	}

	existingServiceMap := make(map[string]swarm.Service)
	for _, service := range existingServices {
		existingServiceMap[service.Spec.Name] = service
	}

	for internalName, serviceSpec := range services {
		name := namespace.Scope(internalName)

		encodedAuth := ""
		image := serviceSpec.TaskTemplate.ContainerSpec.Image
		if sendAuth {
			// Retrieve encoded auth token from the image reference
			encodedAuth, err = command.RetrieveAuthTokenFromImage(ctx, dockerCli, image)
			if err != nil {
				return err
			}
		}

		// service stabilization defaults
		stabilizeDelay := time.Duration(5) * time.Second
		stabilizeTimeout := time.Duration(1) * time.Minute

		// override service stabilization default settings based on spec labels
		labels := serviceSpec.TaskTemplate.ContainerSpec.Labels
		if labels["amp.service.stabilize.delay"] != "" {
			stabilizeDelay, err = time.ParseDuration(labels["amp.service.stabilize.delay"])
			if err != nil {
				return err
			}
		}
		if labels["amp.service.stabilize.timeout"] != "" {
			stabilizeTimeout, err = time.ParseDuration(labels["amp.service.stabilize.timeout"])
			if err != nil {
				return err
			}
		}

		// apply service stabilization timeout setting - the service must be stable before the timeout
		ctx, _ := context.WithTimeout(ctx, stabilizeTimeout)
		var imageName string
		var serviceID string

		if service, exists := existingServiceMap[name]; exists {
			fmt.Fprintf(out, "Updating service       %s (id: %s)\n", name, service.ID)
			fmt.Fprintf(out, "service:               %+v\n", service)
			imageName = service.Spec.TaskTemplate.ContainerSpec.Image
			serviceID = service.ID

			updateOpts := types.ServiceUpdateOptions{EncodedRegistryAuth: encodedAuth}

			if resolveImage == ResolveImageAlways || (resolveImage == ResolveImageChanged && image != service.Spec.Labels[convert.LabelImage]) {
				updateOpts.QueryRegistry = true
			}

			response, err := apiClient.ServiceUpdate(
				ctx,
				service.ID,
				service.Version,
				serviceSpec,
				updateOpts,
			)
			if err != nil {
				return errors.Wrapf(err, "failed to update service %s", name)
			}

			for _, warning := range response.Warnings {
				fmt.Fprintln(dockerCli.Err(), warning)
			}
		} else {
			fmt.Fprintf(out, "Creating service %s\n", name)
			createOpts := types.ServiceCreateOptions{EncodedRegistryAuth: encodedAuth}
			// query registry if flag disabling was not set
			if resolveImage == ResolveImageAlways || resolveImage == ResolveImageChanged {
				createOpts.QueryRegistry = true
			}

			var resp types.ServiceCreateResponse
			if resp, err = apiClient.ServiceCreate(ctx, serviceSpec, createOpts); err != nil {
				return errors.Wrapf(err, "failed to create service %s", name)
			}
			fmt.Fprintf(out, "service:               %+v\n", resp)
			serviceID = resp.ID
			imageName = serviceSpec.TaskTemplate.ContainerSpec.Image
		}

		fmt.Fprintf(out, "image:                 %s\n", imageName)
		fmt.Fprintf(out, "Stabilization delay:   %s\n", stabilizeDelay)
		fmt.Fprintf(out, "Stabilization timeout: %s\n", stabilizeTimeout)
		done := make(chan error)

		// create a watcher for service/container events based on the service image
		options := NewEventsWatcherOptions(events.ServiceEventType, events.ContainerEventType)
		options.AddImageFilter(imageName)

		w := NewEventsWatcherWithCancel(ctx, apiClient, options)
		w.On("*", func(m events.Message) {
			//fmt.Fprintf(out, "EVENT: %s\n", MessageString(m))
		})
		w.OnError(func(err error) {
			//fmt.Fprintf(out, "OnError: %s\n", err)
			w.Cancel()
			done <- err
		})
		w.Watch()

		NotifyState(ctx, apiClient, serviceID, expectedState, stabilizeDelay, func(err error) {
			done <- err
		})

		err = <-done
		// unlike what docker does with stack deployment,
		// we consider that a failing service should fail the stack deployment
		if err != nil {
			w.Cancel()
			return err
		}
	}
	return nil
}

// MessageString returns a formatted event message
func MessageString(m events.Message) string {
	a := ""
	for k, v := range m.Actor.Attributes {
		a += fmt.Sprintf("    %s: %s\n", k, v)
	}
	return fmt.Sprintf("ID: %s\n  Status: %s\n  From: %s\n  Type: %s\n  Action: %s\n  Actor ID: %s\n  Actor Attributes: \n%s\n  Scope: %s\n  Time: %d\n  TimeNano: %d\n\n",
		m.ID, m.Status, m.From, m.Type, m.Action, m.Actor.ID, a, m.Scope, m.Time, m.TimeNano)
}

// NotifyState calls the provided callback when the desired service state is achieved for all tasks or when the deadline is exceeded
func NotifyState(ctx context.Context, apiClient apiclient.APIClient, serviceID string, desiredState swarm.TaskState, stabilizeDelay time.Duration, callback func(error)) {
	deadline, isSet := ctx.Deadline()
	if !isSet {
		deadline = time.Now().Add(1 * time.Minute)
	}

	go func() {
		taskOpts := types.TaskListOptions{}
		taskOpts.Filters = filters.NewArgs()
		taskOpts.Filters.Add("service", serviceID)
		counter := 0

		for {
			// all tasks need to match the desired state and be stable within the deadline
			if time.Now().After(deadline) {
				callback(errors.New("failed to achieve desired state before deadline"))
				return
			}

			// get tasks
			tasks, err := ListTasks(ctx, apiClient, taskOpts)
			if err != nil {
				callback(err)
				return
			}

			// if *any* task does not match the desired state then wait for another loop iteration to check again
			failure := false
			for _, t := range tasks {
				if t.Status.State != desiredState {
					failure = true
					break
				}
			}

			// all tasks matched the desired state - now wait for things to stabilize, or if already stabilized,
			// then callback with success
			if !failure {
				if counter < 1 {
					// make sure we have enough time to wait for things to stabilize within the deadline
					if time.Now().Add(stabilizeDelay).After(deadline) {
						callback(errors.New("failed to achieve desired state with stabilization delay before deadline"))
						return
					}
					time.Sleep(stabilizeDelay)
					counter++
				} else {
					// success!
					callback(nil)
					return
				}
			}

			// task polling interval
			time.Sleep(1 * time.Second)
		}
	}()
}
