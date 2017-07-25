package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/appcelerator/amp/cluster/agent/admin"
	"github.com/appcelerator/amp/cluster/agent/pkg/docker"
	"github.com/appcelerator/amp/cluster/agent/pkg/docker/stack"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/compose/convert"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"
)

const (
	TARGET_SINGLE  = "single"
	TARGET_CLUSTER = "cluster"
)

func NewInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up amp services in swarm environment",
		RunE:  install,
	}
	return installCmd
}

func install(cmd *cobra.Command, args []string) error {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := docker.NewDockerCli(stdin, stdout, stderr)

	// Create initial secrets
	createInitialSecrets()

	// Create initial networks
	createInitialNetworks()

	namespace := "amp"
	if len(args) > 0 {
		namespace = args[0]
	}

	etcdClusterMode, err := serviceDeploymentMode(dockerCli.Client(), "amp.type.kv", "true")
	if err != nil {
		return err
	}
	elasticsearchClusterMode, err := serviceDeploymentMode(dockerCli.Client(), "amp.type.search", "true")
	if err != nil {
		return err
	}
	clusterMode := map[string]string{"elasticsearch": elasticsearchClusterMode, "etcd": etcdClusterMode}
	files, err := getStackFiles("./stacks", clusterMode)
	if err != nil {
		return err
	}

	for _, f := range files {
		log.Println(f)
		if strings.Contains(f, "test") {
			err := deployTest(dockerCli, f, "test", 60 /* timeout in seconds */)
			stack.Remove(dockerCli, stack.RemoveOptions{Namespaces: []string{"test"}})
			if err != nil {
				return err
			}
		} else {
			err := deploy(dockerCli, f, namespace)
			if err != nil {
				return err
			}
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}

// returns the deployment mode
// based on the number of nodes with the label passed as argument
// if number of nodes > 2, mode = cluster, else mode = single
func serviceDeploymentMode(c client.APIClient, labelKey string, labelValue string) (string, error) {
	// unfortunately filtering labels on NodeList won't work as expected, Cf. https://github.com/moby/moby/issues/27231
	nodes, err := c.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return "", err
	}
	matchingNodes := 0
	for _, node := range nodes {
		// node is a swarm.Node
		for k, v := range node.Spec.Labels {
			if k == labelKey {
				if labelValue == "" || labelValue == v {
					matchingNodes++
				}
			}
		}
	}
	switch matchingNodes {
	case 0:
		return "", fmt.Errorf("can't find a node with label %s", labelKey)
	case 1:
		fallthrough
	case 2:
		return TARGET_SINGLE, nil
	default:
		return TARGET_CLUSTER, nil
	}
}

// returns sorted list of yaml file pathnames
func getStackFiles(path string, clusterMode map[string]string) ([]string, error) {
	if path == "" {
		path = "./stacks"
	}

	// a bit more work but we can't just use filepath.Glob
	// since we need to match both *.yml and *.yaml
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	stackfiles := []string{}
	for _, f := range files {
		name := f.Name()
		// not compiling regex since only expecting less than a dozen stackfiles
		matched, err := regexp.MatchString("\\.ya?ml$", name)
		if err != nil {
			log.Println(err)
		} else if matched {
			// looking for the service name, in case there's an indication for the cluster mode (single vs cluster)
			// expecting a file with a name NN-SERVICENAME-mode.*
			split := strings.Split(name, "-")
			if len(split) == 3 {
				serviceName := split[1]
				if strings.Contains(name, TARGET_SINGLE) && clusterMode[serviceName] != TARGET_SINGLE {
					continue
				}
				if strings.Contains(name, TARGET_CLUSTER) && clusterMode[serviceName] != TARGET_CLUSTER {
					continue
				}
			}
			stackfiles = append(stackfiles, filepath.Join(path, name))
		}
	}
	return stackfiles, nil
}

func deploy(d *command.DockerCli, stackfile string, namespace string) error {
	return deployExpectingState(d, stackfile, namespace, swarm.TaskStateRunning)
}

func deployExpectingState(d *command.DockerCli, stackfile string, namespace string, expectedState swarm.TaskState) error {
	if namespace == "" {
		// use the stackfile basename as the default stack namespace
		namespace = filepath.Base(stackfile)
		namespace = strings.TrimSuffix(namespace, filepath.Ext(namespace))
	}

	opts := stack.DeployOptions{
		Namespace:        namespace,
		Composefile:      stackfile,
		ResolveImage:     stack.ResolveImageNever,
		SendRegistryAuth: false,
		Prune:            false,
		ExpectedState:    expectedState,
	}

	return stack.Deploy(context.Background(), d, opts)
}

func deployTest(d *command.DockerCli, stackfile string, namespace string, timeout int) error {
	// Deploy the test stack
	if err := deployExpectingState(d, stackfile, namespace, swarm.TaskStateComplete); err != nil {
		return err
	}

	// Create a docker client
	c, err := client.NewClient(admin.DefaultURL, admin.DefaultVersion, nil, nil)
	if err != nil {
		return err
	}

	// List stack tasks
	options := types.TaskListOptions{Filters: filters.NewArgs()}
	options.Filters.Add("label", convert.LabelNamespace+"="+namespace)
	tasks, err := stack.ListTasks(context.Background(), c, options)
	if err != nil {
		return err
	}

	// Assert we have at least one task
	if len(tasks) == 0 {
		return fmt.Errorf("no task for test")
	}

	// Assert we have only one task
	if len(tasks) != 1 {
		return fmt.Errorf("too many tasks for test: %d", len(tasks))
	}

	// If the task has an error, the test has failed
	task := tasks[0]
	if task.Status.Err != "" {
		return fmt.Errorf("test failed with status: %s", task.Status.Err)
	}

	log.Println("Test successful")
	return nil
}

func ListSecrets() ([]swarm.Secret, error) {
	c, err := client.NewClient(admin.DefaultURL, admin.DefaultVersion, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.SecretList(context.Background(), types.SecretListOptions{})
}

func SecretExists(name string) (bool, error) {
	secrets, err := ListSecrets()
	if err != nil {
		return false, err
	}
	for _, secret := range secrets {
		if secret.Spec.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func CreateSecret(name string, data []byte) error {
	c, err := client.NewClient(admin.DefaultURL, admin.DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: name,
		},
		Data: data,
	}
	_, err = c.SecretCreate(context.Background(), spec)
	if err != nil {
		return err
	}
	return nil
}

// AMP secrets map: Secret name paired to secret file in ./defaults
var ampSecrets = map[string]string{
	"alertmanager_yml": "alertmanager.yml",
	"amplifier_yml":    "amplifier.yml",
	"certificate_amp":  "certificate.amp",
}

// This is the default secrets path
const defaultSecretsPath = "defaults"

// exists returns whether the given file or directory exists or not
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func createInitialSecrets() error {
	// Computing secret path
	secretPath := path.Join("/", defaultSecretsPath)
	pe, err := pathExists(secretPath)
	if err != nil {
		return err
	}
	if !pe {
		secretPath = defaultSecretsPath
	}
	secretPath, err = filepath.Abs(secretPath)
	if err != nil {
		return err
	}
	log.Println("Using the following path for secrets:", secretPath)

	// Creating secrets
	for secret, filename := range ampSecrets {
		// Check if secret already exists
		exists, err := SecretExists(secret)
		if err != nil {
			return err
		}
		if exists {
			log.Println("Skipping already existing secret:", secret)
			continue
		}

		// Load secret data
		data, err := ioutil.ReadFile(path.Join(secretPath, filename))
		if err != nil {
			return err
		}

		// Create secret
		if err := CreateSecret(secret, data); err != nil {
			return err
		}
		log.Println("Successfully created secret:", secret)
	}
	return nil
}

func ListNetworks() ([]types.NetworkResource, error) {
	c, err := client.NewClient(admin.DefaultURL, admin.DefaultVersion, nil, nil)
	if err != nil {
		return nil, err
	}
	return c.NetworkList(context.Background(), types.NetworkListOptions{})
}

func NetworkExists(name string) (bool, error) {
	networks, err := ListNetworks()
	if err != nil {
		return false, err
	}
	for _, network := range networks {
		if network.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func CreateNetwork(name string, overlay bool, attachable bool) error {
	c, err := client.NewClient(admin.DefaultURL, admin.DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	spec := types.NetworkCreate{
		CheckDuplicate: true,
		Attachable:     attachable,
	}
	if overlay {
		spec.Driver = "overlay"
	}
	_, err = c.NetworkCreate(context.Background(), name, spec)
	if err != nil {
		return err
	}
	return nil
}

const ampnet = "ampnet"

func createInitialNetworks() error {
	// Check if network already exists
	exists, err := NetworkExists(ampnet)
	if err != nil {
		return err
	}
	if exists {
		log.Println("Skipping already existing network:", ampnet)
		return nil
	}

	// Create network
	if err := CreateNetwork(ampnet, true, true); err != nil {
		return err
	}
	log.Println("Successfully created network:", ampnet)
	return nil
}
