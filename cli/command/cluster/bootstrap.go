package cluster

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/appcelerator/amp/cli"
)

const (
	bootstrapImg = "appcelerator/amp-bootstrap:%s"
	bootstrapTag = "0.13.1"
)

var (
	dockerArgs []string
)

func init() {
	dockerArgs = []string{
		"run", "-t", "--rm", "--name", "amp-bootstrap",
		"--network", "host",
		"--label", "io.amp.role=infrastructure",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-e", "GOPATH=/go",
	}
}

func Run(c cli.Interface, args []string, env map[string]string) error {
	// make environment variables available to container
	if env != nil {
		for k, v := range env {
			dockerArgs = append(dockerArgs, "-e", fmt.Sprintf("%s=%s", k, v))
		}
	}

	// update the bootstrapImg template string to use either the default tag
	// or the tag specified by the TAG environment variable
	img := bootstrapImg
	tag := bootstrapTag
	if env != nil && env["TAG"] != "" {
		tag = env["TAG"]
	}
	img = fmt.Sprintf(img, tag)

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
