package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"docker.io/go-docker"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/flags"
)

// Docker constants
const (
	DefaultURL          = "unix:///var/run/docker.sock"
	DefaultVersion      = "1.32"
	StateStarting       = "STARTING"
	StateRunning        = "RUNNING"
	StateError          = "ERROR"
	StateNoMatchingNode = "NO MATCHING NODE"
)

// Docker wrapper
type Docker struct {
	client  *docker.Client
	url     string
	version string
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
	if d.client, err = docker.NewClient(d.url, d.version, nil, nil); err != nil {
		return fmt.Errorf("unable to connect to Docker at %s: %v", d.url, err)
	}
	return nil
}

// GetClient returns the native docker client
func (d *Docker) GetClient() *docker.Client {
	return d.client
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
