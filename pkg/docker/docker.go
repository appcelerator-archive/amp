package docker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/appcelerator/amp/api/rpc/cluster/constants"
	"github.com/appcelerator/amp/pkg/docker/docker/stack"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/compose/loader"
	types2 "github.com/docker/docker/cli/compose/types"
	"github.com/docker/docker/cli/flags"
	"github.com/docker/docker/client"
	dopts "github.com/docker/docker/opts"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/manager/constraint"
	"golang.org/x/net/context"
)

// Docker constants
const (
	DefaultURL               = "unix:///var/run/docker.sock"
	DefaultVersion           = "1.27"
	StackStateStarting       = "STARTING"
	StackStateRunning        = "RUNNING"
	StackStateNoMatchingNode = "NO MATCHING NODE"
	NoMatchingNodes          = -1
)

// Docker wrapper
type Docker struct {
	client    *client.Client
	connected bool
	url       string
	version   string
}

type StackStatus struct {
	RunningServices int32
	TotalServices   int32
	Status          string
}

type ServiceStatus struct {
	RunningTasks int32
	TotalTasks   int32
	Status       string
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
func (d *Docker) StackServices(ctx context.Context, stackName string, quietOption bool) (output string, err error) {
	cmd := func(cli *command.DockerCli) error {
		servicesOpt := stack.NewServicesOptions(quietOption, "", dopts.NewFilterOpt(), stackName)
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
	var readyServices int32 = 0
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
		if status.Status == StackStateRunning {
			readyServices++
		}
	}
	if readyServices == totalServices {
		return &StackStatus{
			RunningServices: readyServices,
			TotalServices:   totalServices,
			Status:          StackStateRunning,
		}, nil
	}
	return &StackStatus{
		RunningServices: readyServices,
		TotalServices:   totalServices,
		Status:          StackStateStarting,
	}, nil
}

// ServiceInspect inspects a service
func (d *Docker) ServiceInspect(ctx context.Context, service string) (swarm.Service, error) {
	serviceEntity, _, err := d.client.ServiceInspectWithRaw(ctx, service, types.ServiceInspectOptions{InsertDefaults: true})
	if err != nil {
		return swarm.Service{}, err
	}
	return serviceEntity, nil
}

// ServiceScale scales a service
func (d *Docker) ServiceScale(ctx context.Context, service string, scale uint64) error {
	serviceEntity, _, err := d.client.ServiceInspectWithRaw(ctx, service, types.ServiceInspectOptions{InsertDefaults: true})
	if err != nil {
		return err
	}
	serviceMode := &serviceEntity.Spec.Mode
	if serviceMode.Replicated == nil {
		return fmt.Errorf("scale can only be used with replicated mode")
	}
	serviceMode.Replicated.Replicas = &scale

	response, err := d.client.ServiceUpdate(ctx, serviceEntity.ID, serviceEntity.Version, serviceEntity.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		return err
	}
	for _, warning := range response.Warnings {
		log.Println(warning)
	}
	log.Printf("service %s scaled to %d\n", service, scale)
	return nil
}

// NodeInspect inspects a node
func (d *Docker) NodeInspect(ctx context.Context, nodeID string) (swarm.Node, error) {
	nodeEntity, _, err := d.client.NodeInspectWithRaw(ctx, nodeID)
	if err != nil {
		return swarm.Node{}, err
	}
	return nodeEntity, nil
}

// ExpectedNumberOfTasks returns expected number of tasks of a service
func (d *Docker) ExpectedNumberOfTasks(ctx context.Context, serviceID string) (int, error) {
	var expectedTasks int
	serviceInfo, err := d.ServiceInspect(ctx, serviceID)
	if err != nil {
		return 0, err
	}
	matchingNodeCount, err := d.numberOfMatchingNodes(ctx, serviceInfo)
	if err != nil {
		return 0, err
	}
	if serviceInfo.Spec.Mode.Global != nil {
		expectedTasks = matchingNodeCount
	} else {
		expectedTasks = int(*serviceInfo.Spec.Mode.Replicated.Replicas)
	}
	return expectedTasks, nil
}

// numberOfMatchingNodes returns number of nodes matching placement constraints
func (d *Docker) numberOfMatchingNodes(ctx context.Context, serviceInfo swarm.Service) (int, error) {
	var matchingNodes int
	// placement constraints
	constraints, _ := constraint.Parse(serviceInfo.Spec.TaskTemplate.Placement.Constraints)
	// list all nodes in the swarm
	nodes, err := d.client.NodeList(ctx, types.NodeListOptions{})
	if err != nil {
		return 0, err
	}
	// inspect every node on the swarm to check for satisfying constraints
	for _, node := range nodes {
		apiNode := d.nodeToNode(ctx, node)
		if !constraint.NodeMatches(constraints, apiNode) {
			return NoMatchingNodes, nil
		}
		matchingNodes++
	}
	return matchingNodes, nil
}

// nodeToNode converts a swarm node type into an api node type
func (d *Docker) nodeToNode(ctx context.Context, swarmNode swarm.Node) *api.Node {
	apiNode := &api.Node{
		ID: swarmNode.ID,
		Status: api.NodeStatus{
			Addr:  swarmNode.Status.Addr,
			State: api.NodeStatus_State(api.NodeStatus_State_value[strings.ToUpper(string(swarmNode.Status.State))]),
		},
		Spec: api.NodeSpec{
			Availability: api.NodeSpec_Availability(api.NodeSpec_Availability_value[strings.ToUpper(string(swarmNode.Spec.Availability))]),
			Annotations: api.Annotations{
				Labels: swarmNode.Spec.Labels,
			},
		},
		Description: &api.NodeDescription{
			Hostname: swarmNode.Description.Hostname,
			Platform: &api.Platform{
				OS:           swarmNode.Description.Platform.OS,
				Architecture: swarmNode.Description.Platform.Architecture,
			},
			Engine: &api.EngineDescription{
				EngineVersion: swarmNode.Description.Engine.EngineVersion,
				Labels:        swarmNode.Description.Engine.Labels,
			},
		},
		Role: api.NodeRole(api.NodeRole_value[strings.ToUpper(string(swarmNode.Spec.Role))]),
	}
	if swarmNode.ManagerStatus != nil {
		apiNode.ManagerStatus = &api.ManagerStatus{
			Leader:       swarmNode.ManagerStatus.Leader,
			Addr:         swarmNode.ManagerStatus.Addr,
			Reachability: api.RaftMemberStatus_Reachability(api.RaftMemberStatus_Reachability_value[strings.ToUpper(string(swarmNode.ManagerStatus.Reachability))]),
		}
	}

	for _, plugin := range swarmNode.Description.Engine.Plugins {
		apiNode.Description.Engine.Plugins = append(apiNode.Description.Engine.Plugins, api.PluginDescription{Type: plugin.Type, Name: plugin.Name})
	}
	return apiNode
}

// ServiceList list the services
func (d *Docker) ServicesList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	if !d.connected {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}
	return d.client.ServiceList(ctx, options)
}

// readyTasks returns the running tasks of a service
func (d *Docker) readyTasks(ctx context.Context, service string) (int, error) {
	args := filters.NewArgs()
	args.Add("service", service)
	var readyTasks int = 0
	serviceTasks, err := d.TaskList(ctx, types.TaskListOptions{Filters: args})
	if err != nil {
		return 0, err
	}
	for _, serviceTask := range serviceTasks {
		if serviceTask.Status.State == swarm.TaskStateRunning {
			readyTasks++
		}
	}
	return readyTasks, nil
}

// ServiceStatus returns service status
func (d *Docker) ServiceStatus(ctx context.Context, service string) (*ServiceStatus, error) {
	readyTasks, err := d.readyTasks(ctx, service)
	if err != nil {
		return &ServiceStatus{}, err
	}
	totalTasks, err := d.ExpectedNumberOfTasks(ctx, service)
	if err != nil {
		return &ServiceStatus{}, err
	}
	if readyTasks == NoMatchingNodes {
		return &ServiceStatus{
			RunningTasks: 0,
			TotalTasks:   0,
			Status:       StackStateNoMatchingNode,
		}, nil
	}
	if readyTasks == totalTasks {
		return &ServiceStatus{
			RunningTasks: int32(readyTasks),
			TotalTasks:   int32(totalTasks),
			Status:       StackStateRunning,
		}, nil
	}
	return &ServiceStatus{
		RunningTasks: int32(readyTasks),
		TotalTasks:   int32(totalTasks),
		Status:       StackStateStarting,
	}, nil
}
