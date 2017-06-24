package cluster

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"

	"github.com/appcelerator/amp/cli"
	"github.com/docker/docker/pkg/stringid"
)

// Supported plugin providers used by the factory function `NewPlugin`
const (
	Local = "local"
	AWS   = "aws"
)

// PluginConfig is used by the factory function `NewClusterPlugin` to create a new plugin instance.
type PluginConfig struct {
	// Provider is the name of the cluster provider, such as "local" or "aws"
	Provider string
	Options  map[string]string
	DockerOpts docker
}

// Plugin declares the methods that all plugin providers,
// such as local and aws, must implement
type Plugin interface {
	// Provider returns the name of the provider, such as "local" or "aws"
	Provider() string

	// Run executes the plugin with the specified arguments and environment variables
	Run(c cli.Interface, args []string, env map[string]string) error
}

// NewPlugin is a simple factory function to return a new instance of
// a specific cluster plugin based on the supplied config
// (config.Provider must be set to a valid provider or this
// function will return an error).
func NewPlugin(config PluginConfig) (Plugin, error) {
	var p Plugin

	switch config.Provider {
	case Local:
		p = &localPlugin{plugin{config: config}}
	case AWS:
		p = &awsPlugin{plugin{config: config}}
	default:
		return nil, errors.New(fmt.Sprintf("Not a valid plugin provider: %s", config.Provider))
	}

	return p, nil
}

// RunContainer starts a container using the specified image for the cluster plugin.
// Cluster plugin commands are `init`, `update`, and `destroy` (provided as the single
// `args` value). Additional arguments are supplied as environment variables in `env`, not `args`.
func RunContainer(c cli.Interface, img string, dockerOpts docker, args []string, env map[string]string) error {
	dockerArgs := []string{
		"run", "-t", "--rm", "--name", fmt.Sprintf("amp-cluster-plugin-%s", stringid.GenerateNonCryptoID()),
		"--network", "host",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-e", "GOPATH=/go",
	}

	for _, v := range dockerOpts.volumes {
		dockerArgs = append(dockerArgs, "-v", v)
	}

	// make environment variables available to container
	if env != nil {
		for k, v := range env {
			dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("%s=%s", k, v))
		}
	}

	// this completes the docker args
	dockerArgs = append(dockerArgs, img)

	cmd := "docker"
	args = append(dockerArgs, args...)

	proc := exec.Command(cmd, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	outscanner := bufio.NewScanner(stdout)
	go func() {
		for outscanner.Scan() {
			c.Console().Println(outscanner.Text())
		}
	}()

	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			c.Console().Println(errscanner.Text())
		}
	}()

	err = proc.Start()
	if err != nil {
		return err
	}

	err = proc.Wait()
	if err != nil {
		return err
	}

	return nil
}

// ========================================================
// base plugin implementation - should never be instantiated
// ========================================================

type plugin struct {
	config PluginConfig
}

func (p *plugin) Provider() string {
	return p.config.Provider
}

func (p *plugin) Run(c cli.Interface, args []string, env map[string]string) error {
	return errors.New("Run method invoked on base plugin type")
}

// ========================================================
// local plugin implementation
// ========================================================
type localPlugin struct {
	plugin
}

func (p *localPlugin) Run(c cli.Interface, args []string, env map[string]string) error {
	return queryCluster(c, args, env)
}

// ========================================================
// aws plugin implementation
// ========================================================
type awsPlugin struct {
	plugin
}

func (p *awsPlugin) Run(c cli.Interface, args []string, env map[string]string) error {
	img := "appcelerator/amp-aws"
	dockerOpts := p.plugin.config.DockerOpts
	return RunContainer(c, img, dockerOpts, args, env)
}

