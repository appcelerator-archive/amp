package cluster

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/cli"
	aws "github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"github.com/appcelerator/amp/docker/docker/pkg/stringid"
	"github.com/mitchellh/go-homedir"
)

// Supported plugin providers used by the factory function `NewPlugin`
const (
	Local = "local"
	AWS   = "aws"
)

// PluginConfig is used by the factory function `NewClusterPlugin` to create a new plugin instance.
type PluginConfig struct {
	// Provider is the name of the cluster provider, such as "local" or "aws"
	Provider   string
	Options    map[string]string
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
	if config.Provider == "" {
		return nil, errors.New(fmt.Sprintf("Must specify a plugin provider: %s", config.Provider))
	}

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

func killContainer(c cli.Interface, name string, sig string) error {
	cmd := "docker"
	dockerArgs := []string{
		"kill", "--signal", sig, name,
	}
	proc := exec.Command(cmd, dockerArgs...)
	var out bytes.Buffer
	var e bytes.Buffer
	proc.Stdout = &out
	proc.Stderr = &e
	err := proc.Run()
	if err != nil {
		c.Console().Printf(out.String())
		c.Console().Printf(e.String())
		c.Console().Printf("failed to kill container %s, you may have to remove it manually\n", name)
	} else {
		c.Console().Printf("Plugin container %s has been successfully stopped\n", name)
	}
	return err
}

// RunContainer starts a container using the specified image for the cluster plugin.
// Cluster plugin commands are `init`, `update`, and `destroy` (provided as the single
// `args` value). Additional arguments are supplied as environment variables in `env`, not `args`.
// If f is not nil, then your func will be called (as a goroutine) with stdout from the container process;
// otherwise stdout from the container will be printed to the amp console stdout.
func RunContainer(c cli.Interface, img string, dockerOpts docker, args []string, env map[string]string, f func(r io.Reader, c chan bool)) error {
	containerName := fmt.Sprintf("amp-cluster-plugin-%s", stringid.GenerateNonCryptoID())
	dockerArgs := []string{
		"run", "-t", "--rm", "--name", containerName,
		"--network", "host",
		"--label", "io.amp.role=infrastructure",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-e", "GOPATH=/go",
	}

	// mount configured volumes
	for _, v := range dockerOpts.Volumes {
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

	interruption := make(chan os.Signal, 1)
	signalCaught := make(chan bool, 1)
	signal.Notify(interruption, os.Interrupt, os.Kill)
	go func() {
		sig := <-interruption
		signalCaught <- true
		c.Console().Printf("CLI received signal %s\n", sig.String())
		_ = killContainer(c, containerName, "INT")
		return
	}()

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}

	stdOutDone := make(chan bool)
	if f != nil {
		go f(stdout, stdOutDone)
	} else {
		outscanner := bufio.NewScanner(stdout)
		go func() {
			for outscanner.Scan() {
				c.Console().Println(outscanner.Text())
			}
			stdOutDone <- true
		}()
	}

	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}

	stdErrDone := make(chan bool)
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			c.Console().Println(errscanner.Text())
		}
		stdErrDone <- true
	}()

	err = proc.Start()
	if err != nil {
		return err
	}

	<-stdOutDone
	<-stdErrDone

	err = proc.Wait()
	if err != nil {
		// if it returns directly, we won't be able to process the interrupt signal
		select {
		case <-signalCaught:
			time.Sleep(5 * time.Second)
		default:
			return err
		}
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
	dockerOpts := p.config.DockerOpts
	if dockerOpts.Volumes == nil {
		p.config.DockerOpts.Volumes = []string{}
	}

	img := fmt.Sprintf("appcelerator/amp-local:%s", c.Version())
	return RunContainer(c, img, dockerOpts, args, env, nil)
}

// ========================================================
// aws plugin implementation
// ========================================================
type awsPlugin struct {
	plugin
}

func decodeAwsPluginOutput(c cli.Interface, d *json.Decoder) error {
	var po aws.PluginOutput
	if err := d.Decode(&po); err != nil {
		return err
	}
	// if there's an event, that's the relevant information
	if po.Event != nil {
		e := po.Event
		status := e.ResourceStatus
		if e.ResourceStatusReason != "" {
			status = fmt.Sprintf("%s (%s)", e.ResourceStatus, e.ResourceStatusReason)
		}
		c.Console().Printf("%-31s %-28s %s\n", e.Timestamp, e.LogicalResourceId, status)
		if e.ResourceType == "AWS::CloudFormation::Stack" && (e.ResourceStatus == "ROLLBACK_COMPLETE" || e.ResourceStatus == "DELETE_COMPLETE") {
			return fmt.Errorf("deployment failed")
		}
		// the caller should loop on the reader
		return nil
	}
	// next, look for an output
	if po.Output != nil {
		c.Console().Printf("-------------------------------------------------------------------------------\n")
		for _, o := range po.Output {
			c.Console().Printf("%-42s | %s\n", o.Description, o.OutputValue)
		}
		return nil
	}
	// last, look for an error
	if po.Error != "" {
		c.Console().Printf("Error: %s\n", po.Error)
		if strings.Contains(po.Error, "Throttling: Rate exceeded") == true {
			// we shouldn't stop
			return nil
		}
		return fmt.Errorf("deployment failed, see error details above")
	}
	return nil
}

func (p *awsPlugin) Run(c cli.Interface, args []string, env map[string]string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	dockerOpts := p.config.DockerOpts

	if dockerOpts.Volumes == nil {
		p.config.DockerOpts.Volumes = []string{}
	}

	mountAwsConfigFolder := true
	for _, a := range args {
		if a == "--access-key-id" {
			mountAwsConfigFolder = false
			break
		}
	}
	if mountAwsConfigFolder {
		// automatically attempt to mount aws credentials if present
		awshome := path.Join(home, ".aws")
		awscreds := path.Join(awshome, "credentials")
		if _, err := os.Stat(awscreds); err != nil {
			return fmt.Errorf("no $HOME/.aws/credentials found, please refer to http://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html to configure it")
		}
		dockerOpts.Volumes = append(dockerOpts.Volumes, fmt.Sprintf("%s:/root/.aws", awshome))
	}

	// function to print the aws plugin output
	f := func(r io.Reader, done chan bool) {
		d := json.NewDecoder(r)
		for {
			if err = decodeAwsPluginOutput(c, d); err == io.EOF {
				done <- true
				break
			} else if err != nil {
				done <- true
				// If there is an error here, it is because the plugin itself
				// failed to return json for plugin errors
				// so print the buffer for now
				outscanner := bufio.NewScanner(d.Buffered())
				for outscanner.Scan() {
					c.Console().Println(outscanner.Text())
				}
				c.Console().Error(err) // these errors need to be fixed
				return
			}
		}
	}

	img := fmt.Sprintf("appcelerator/amp-aws:%s", c.Version())
	return RunContainer(c, img, dockerOpts, args, env, f)
}
