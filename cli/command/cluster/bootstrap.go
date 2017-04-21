package cluster

import (
	"bufio"
	"os/exec"

	"github.com/appcelerator/amp/cli"
)

var (
	dockerArgs []string
)

func init() {
	dockerArgs = []string{
		"run", "-t", "--rm", "--name", "amp-bootstrap",
		"--network", "host",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-e", "GOPATH=/go",
		"appcelerator/amp-bootstrap:1.0.2",
	}
}

func Run(c cli.Interface, args []string) error {
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
