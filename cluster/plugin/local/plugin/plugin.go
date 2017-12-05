package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/filters"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/docker/docker/pkg/stdcopy"
)

const (
	InitTimeout            = 10
	ContainerName          = "ampagent"
	ImageName              = "appcelerator/ampagent"
	DockerSocket           = "/var/run/docker.sock"
	DockerSwarmSocket      = "/var/run/docker"
	CoreStackName          = "amp"
	MaxMapCountRequirement = 262144 // as specified here: https://www.elastic.co/guide/en/elasticsearch/reference/current/vm-max-map-count.html
)

var (
	ContainerLabels = map[string]string{"io.amp.role": "infrastructure"}
)

// RequestOptions stores parameters for the Docker API
type RequestOptions struct {
	InitRequest swarm.InitRequest
	// Node labels
	Labels map[string]string
	// Tag of the ampagent image
	Tag           string
	Registration  string
	Notifications bool
	ForceLeave    bool
	NoLogs        bool
	NoMetrics     bool
	NoProxy       bool
}

type FullSwarmInfo struct {
	Swarm swarm.Swarm
	Node  swarm.Node
}
type ShortSwarmInfo struct {
	SwarmStatus  string `json:"Swarm Status"`
	CoreServices int    `json:"Core Services"`
	UserServices int    `json:"User Services"`
}

// Check prerequisites
func CheckPrerequisites(opts *RequestOptions) error {
	switch runtime.GOOS {
	case "linux":
		if opts.NoLogs {
			return nil
		}

		// Check max_map_count settings for Elasticsearch
		path := path.Join("proc", "sys", "vm", "max_map_count")
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read system configuration at %s: %s", path, err.Error())
		}
		value := strings.TrimSpace(string(bytes))
		maxMapCount, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("unable to convert value %s: %s", value, err.Error())
		}
		if maxMapCount < MaxMapCountRequirement {
			return fmt.Errorf("vm.max_map_count value is too low. Got: %d, expected: >= %d, please check the documentation on local cluster deployment prerequisites", maxMapCount, MaxMapCountRequirement)
		}
		return nil

	default:
		return nil
	}
	return nil
}

// EnsureSwarmExists checks that the Swarm is initialized, and does it if it's not the case
func EnsureSwarmExists(ctx context.Context, c *docker.Client, opts *RequestOptions) error {
	timeout := make(chan bool, 1)
	done := make(chan bool, 1)
	var err error
	// the Init method may freeze if the Docker engine has issues (and it often has)
	go func() {
		time.Sleep(InitTimeout * time.Second)
		timeout <- true
	}()
	go func() {
		_, err = c.SwarmInit(ctx, opts.InitRequest)
		if err != nil {
			// if the swarm is already initialized, ignore the error
			if strings.Contains(fmt.Sprintf("%v", err), "This node is already part of a swarm") {
				err = nil
			} else {
				fmt.Printf("%v\n", err)
			}
		}
		done <- true
	}()
	select {
	case <-done:
		return err
	case <-timeout:
		return fmt.Errorf("Timed out")
	}
}

func LabelNode(ctx context.Context, c *docker.Client, opts *RequestOptions) error {
	node, err := InfoNode(ctx, c)
	if err != nil {
		return err
	}
	version := node.Meta.Version
	nodeSpec := node.Spec
	for label := range opts.Labels {
		node.Spec.Annotations.Labels[label] = opts.Labels[label]
	}
	return c.NodeUpdate(ctx, node.ID, version, nodeSpec)
}

func removeAgent(ctx context.Context, c *docker.Client, cid string, force bool) error {
	return c.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{Force: force})
}

func containerLogs(ctx context.Context, c *docker.Client, id string) {
	reader, err := c.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return
	}
	defer reader.Close()
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
}

// RunAgent runs the ampagent image to init (action ="install") or destroy (action="uninstall")
func RunAgent(ctx context.Context, c *docker.Client, action string, opts *RequestOptions) error {
	// first remove any exited or dead container with same name
	// this can happen when Docker crashes while the plugin and the agent are running
	if err := removeExitedAgent(c); err != nil {
		return err
	}
	image := fmt.Sprintf("%s:%s", ImageName, opts.Tag)
	config := container.Config{
		Image: image,
		Env: []string{
			fmt.Sprintf("TAG=%s", opts.Tag),
			fmt.Sprintf("REGISTRATION=%s", opts.Registration),
			fmt.Sprintf("NOTIFICATIONS=%t", opts.Notifications),
		},
		Labels: ContainerLabels,
		Tty:    false,
	}
	var actionArgs []string
	if opts.NoLogs {
		actionArgs = append(actionArgs, "--no-logs")
	}
	if opts.NoMetrics {
		actionArgs = append(actionArgs, "--no-metrics")
	}
	if opts.NoProxy {
		actionArgs = append(actionArgs, "--no-proxy")
	}
	switch action {
	case "install":
		action = ""
		config.Cmd = actionArgs
	case "uninstall":
		config.Cmd = append([]string{action}, actionArgs...)
	default:
		return fmt.Errorf("action %s is not implemented", action)
	}
	mounts := []mount.Mount{
		{
			Type:   "bind",
			Source: DockerSocket,
			Target: DockerSocket,
		},
		{
			Type:   "bind",
			Source: DockerSwarmSocket,
			Target: DockerSwarmSocket,
		},
	}
	hostConfig := container.HostConfig{
		AutoRemove: false,
		Mounts:     mounts,
	}
	reader, err := c.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		// don't exit, if not in the registry we may still want to run the container with a local image
		fmt.Println("ampagent image pull failed, which is expected on a development version")
	} else {
		// wait for the image to be pulled
		data := make([]byte, 1000, 1000)
		for {
			_, err := reader.Read(data)
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				return err
			}
		}
	}
	r, err := c.ContainerCreate(ctx, &config, &hostConfig, nil, ContainerName)
	if err != nil {
		return err
	}

	done := make(chan bool, 1)
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Interrupt, os.Kill)
	go func() {
		sig := <-interruption
		fmt.Printf("Received signal %s\n", sig.String())
		err = c.ContainerKill(ctx, r.ID, "SIGINT")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		for i := 0; i < 10; i++ {
			if err := removeAgent(ctx, c, r.ID, false); err == nil {
				break
			}
			time.Sleep(time.Second)
		}
		done <- true
		return
	}()

	if err = c.ContainerStart(ctx, r.ID, types.ContainerStartOptions{}); err != nil {
		_ = removeAgent(ctx, c, r.ID, true)
		return err
	}

	go containerLogs(ctx, c, r.ID)

	go func() {
		for {
			filter := filters.NewArgs()
			filter.Add("id", r.ID)
			l, _ := c.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
			if len(l) == 0 {
				done <- true
				return
			}
			time.Sleep(time.Second)
		}
	}()

	<-done
	_ = removeAgent(ctx, c, r.ID, true)
	return nil
}

// InfoCluster returns the Swarm info
func InfoCluster(ctx context.Context, c *docker.Client) (swarm.Swarm, error) {
	return c.SwarmInspect(ctx)
}

// InfoNode returns the Node info
func InfoNode(ctx context.Context, c *docker.Client) (swarm.Node, error) {
	nodes, err := c.NodeList(ctx, types.NodeListOptions{})
	if len(nodes) != 1 {
		return swarm.Node{}, fmt.Errorf("expected 1 node, got %d", len(nodes))
	}
	node, _, err := c.NodeInspectWithRaw(ctx, nodes[0].ID)
	return node, err
}

// InfoAMPCore returns the number of AMP core services
func InfoAMPCore(ctx context.Context, c *docker.Client) (int, error) {
	var count int
	services, err := c.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return 0, err
	}
	for _, service := range services {
		if strings.HasPrefix(service.Spec.Name, fmt.Sprintf("%s_", CoreStackName)) {
			count++
		}
	}
	return count, nil
}

// InfoUser returns the number of user services
func InfoUser(ctx context.Context, c *docker.Client) (int, error) {
	var count int
	services, err := c.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return 0, err
	}
	for _, service := range services {
		if !strings.HasPrefix(service.Spec.Name, CoreStackName) {
			count++
		}
	}
	return count, err
}

func InfoToJSON(status string, csCount int, usCount int) (string, error) {
	// filter the swarm content
	si := ShortSwarmInfo{SwarmStatus: status, CoreServices: csCount, UserServices: usCount}
	j, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// DeleteSwarm starts the delete operation
func DeleteCluster(ctx context.Context, c *docker.Client, opts *RequestOptions) error {
	if opts.ForceLeave {
		return c.SwarmLeave(ctx, true)
	}
	return nil
}

// SwarmNodeStatus returns the swarm status for this node
func SwarmNodeStatus(c *docker.Client) (swarm.LocalNodeState, error) {
	info, err := c.Info(context.Background())
	if err != nil {
		return "", err
	}
	return info.Swarm.LocalNodeState, nil
}

func removeExitedAgent(c *docker.Client) error {
	filter := filters.NewArgs()
	filter.Add("name", ContainerName)
	l, err := c.ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return err
	}
	if len(l) == 1 {
		switch l[0].State {
		case "exited", "dead":
			fmt.Printf("found an exited ampagent container [%s], removing it\n", l[0].ID)
			return removeAgent(context.Background(), c, l[0].ID, false)
		case "created", "running":
			return fmt.Errorf("ampagent is already %s, this means you probably already have an amp CLI running", l[0].State)
		default:
			return fmt.Errorf("ampagent is in an unexpected state [%s]", l[0].State)
		}
	}
	return nil
}
