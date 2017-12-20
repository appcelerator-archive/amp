package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/command/stack"
	"github.com/appcelerator/amp/docker/docker/pkg/term"
	ampdocker "github.com/appcelerator/amp/pkg/docker"
	"github.com/spf13/cobra"
)

const (
	TARGET_SINGLE  = "single"
	TARGET_CLUSTER = "cluster"
)

type InstallOptions struct {
	NoLogs    bool
	NoMetrics bool
	NoProxy   bool
}

var InstallOpts = &InstallOptions{}

func NewInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up amp services in swarm environment",
		RunE:  Install,
	}

	return installCmd
}

func Install(cmd *cobra.Command, args []string) error {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := ampdocker.NewDockerCli(stdin, stdout, stderr)

	// Create initial secrets
	createInitialSecrets()

	// Create initial configs
	createInitialConfigs()

	// Create initial networks
	createInitialNetworks()

	namespace := "amp"
	if len(args) > 0 && args[0] != "" {
		namespace = args[0]
	}

	// Handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("\nReceived an interrupt signal - removing AMP services")
			stack.RunRemove(dockerCli, stack.RemoveOptions{Namespaces: []string{namespace}})
			os.Exit(1)
		}
	}()

	deploymentMode, err := serviceDeploymentMode(dockerCli.Client(), "amp.type.kv", "true")
	if err != nil {
		return err
	}
	log.Printf("Deploying in %s mode\n", deploymentMode)

	stackFiles, err := getStackFiles("./stacks", deploymentMode)
	if err != nil {
		return err
	}

	for _, stackFile := range stackFiles {
		if strings.Contains(stackFile, "logs") && InstallOpts.NoLogs ||
			strings.Contains(stackFile, "metrics") && InstallOpts.NoMetrics ||
			strings.Contains(stackFile, "proxy") && InstallOpts.NoProxy {
			continue
		}
		log.Println("Deploying stack", stackFile)
		if err := deploy(dockerCli, stackFile, namespace); err != nil {
			stack.RunRemove(dockerCli, stack.RemoveOptions{Namespaces: []string{namespace}})
			return err
		}
	}
	return nil
}

// returns the deployment mode
// based on the number of nodes with the label passed as argument
// if number of nodes > 2, mode = cluster, else mode = single
func serviceDeploymentMode(c docker.APIClient, labelKey string, labelValue string) (string, error) {
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
func getStackFiles(path string, deploymentMode string) ([]string, error) {
	if path == "" {
		path = "./stacks"
	}
	path += "/" + deploymentMode

	// a bit more work but we can't just use filepath.Glob
	// since we need to match both *.yml and *.yaml
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	stackfiles := []string{}
	for _, f := range files {
		name := f.Name()
		if matched, _ := regexp.MatchString("\\.ya?ml$", name); matched {
			stackfiles = append(stackfiles, filepath.Join(path, name))
		}
	}
	return stackfiles, nil
}

func deploy(d *command.DockerCli, stackfile string, namespace string) error {
	if namespace == "" {
		// use the stackfile basename as the default stack namespace
		namespace = filepath.Base(stackfile)
		namespace = strings.TrimSuffix(namespace, filepath.Ext(namespace))
	}

	options := stack.DeployOptions{
		Namespace:        namespace,
		Composefile:      stackfile,
		ResolveImage:     stack.ResolveImageAlways,
		SendRegistryAuth: false,
		Prune:            false,
	}

	if err := stack.RunDeploy(d, options); err != nil {
		return err
	}

	for _, err := range Docker.WaitOnStack(context.Background(), namespace, os.Stdout) {
		if err != nil {
			return err
		}
	}
	return nil
}

// AMP configs map: Config name paired to config file in ./defaults
var ampConfigs = map[string]string{
	"prometheus_alerts_rules": "prometheus_alerts.rules",
}

// This is the default configs path
const defaultConfigsPath = "defaults"

func createInitialConfigs() error {
	// Computing config path
	configPath := path.Join("/", defaultConfigsPath)
	pe, err := pathExists(configPath)
	if err != nil {
		return err
	}
	if !pe {
		configPath = defaultConfigsPath
	}
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return err
	}
	log.Println("Using the following path for configs:", configPath)

	// Creating configs
	for config, filename := range ampConfigs {
		// Check if config already exists
		exists, err := Docker.ConfigExists(config)
		if err != nil {
			return err
		}
		if exists {
			log.Println("Skipping already existing config:", config)
			continue
		}

		// Load config data
		data, err := ioutil.ReadFile(path.Join(configPath, filename))
		if err != nil {
			return err
		}

		// Create config
		if _, err := Docker.CreateConfig(config, data); err != nil {
			return err
		}
		log.Println("Successfully created config:", config)
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
		exists, err := Docker.SecretExists(secret)
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
		if _, err := Docker.CreateSecret(secret, data); err != nil {
			return err
		}
		log.Println("Successfully created secret:", secret)
	}
	return nil
}

var ampnetworks = []string{"public", "monit", "core"}

func createInitialNetworks() error {
	for _, network := range ampnetworks {
		// Check if network already exists
		exists, err := Docker.NetworkExists(network)
		if err != nil {
			return err
		}
		if exists {
			log.Println("Skipping already existing network:", network)
			continue
		}
		if _, err := Docker.CreateNetwork(network, true, true); err != nil {
			return err
		}
		log.Println("Successfully created network:", network)
	}
	return nil
}
