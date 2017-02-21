package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var (
	dockerCmd  = "docker"
	dockerArgs []string
)

func init() {
	dockerArgs = []string{
		"run", "-t", "--rm", "--name", "ampcli",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-e", fmt.Sprintf("DOCKER_CMD=%s", dockerCmd),
		"-e", "GOPATH=/go",
		"--network", "hostnet",
		"appcelerator/amp:local",
	}
}

func main() {
	args := []string{}

	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}

	cmd := "docker"
	args = append(dockerArgs, args...)

	proc := exec.Command(cmd, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		panic(err)
	}
	outscanner := bufio.NewScanner(stdout)
	go func() {
		for outscanner.Scan() {
			fmt.Printf("%s\n", outscanner.Text())
		}
	}()

	stderr, err := proc.StderrPipe()
	if err != nil {
		panic(err)
	}
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			fmt.Fprintf(os.Stderr, "%s\n", errscanner.Text())
		}
	}()

	err = proc.Start()
	if err != nil {
		panic(err)
	}

	err = proc.Wait()
	if err != nil {
		// Just pass along the information that the process exited with a failure;
		// whatever error information it displayed is what the user will see.
		// TODO: return the process exit code
		os.Exit(1)

	}
}
