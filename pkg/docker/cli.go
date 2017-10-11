package docker

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/flags"
)

func NewDockerCli(stdin io.ReadCloser, stdout, stderr io.Writer) *command.DockerCli {
	d := command.NewDockerCli(stdin, stdout, stderr)
	opts := flags.NewClientOptions()
	opts.Common = flags.NewCommonOptions()
	d.Initialize(opts)
	return d
}

func cliWrapper(cmd func(cli *command.DockerCli) error) (string, error) {
	r, w, _ := os.Pipe()
	cli := NewDockerCli(os.Stdin, w, w)
	if err := cmd(cli); err != nil {
		return "", err
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	return string(outs), nil
}
