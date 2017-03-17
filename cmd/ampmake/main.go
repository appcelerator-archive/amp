package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mitchellh/go-homedir"
)

var (
	ug         = "0:0"
	version    = "0.0.0"
	build      = "-"
	owner      = "appcelerator"
	repo       = "github.com/appcelerator/amp"
	dockerCmd  = "sudo docker"
	dockerArgs []string
)

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dockerArgs = []string{
		"run", "-t", "--rm", "--name", "amptools",
		"-u", ug, //fmt.Sprintf("%s:%s", strconv.Itoa(os.Getuid()), strconv.Itoa(os.Getgid())),
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", fmt.Sprintf("%s/.ssh:/root/.ssh:ro", homedir),
		"-v", fmt.Sprintf("%s/.config:/root/.config:ro", homedir),
		"-v", fmt.Sprintf("%s:/go/src/%s", wd, repo),
		"-w", fmt.Sprintf("/go/src/%s", repo),
		"-e", fmt.Sprintf("VERSION=%s", version),
		"-e", fmt.Sprintf("BUILD=%s", build),
		"-e", fmt.Sprintf("OWNER=%s", owner),
		"-e", fmt.Sprintf("REPO=%s", repo),
		"-e", fmt.Sprintf("DOCKER_CMD=%s", dockerCmd),
		"-e", "GOPATH=/go",
		"appcelerator/amptools:1.1.0",
	}
}

func main() {
	args := []string{
		"make",
		"-f",
		"Makefile.refactor.make",
	}

	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
		fmt.Println(strings.Join(args, " "))
	}

	cmd := "docker"
	args = append(dockerArgs, args...)
	//fmt.Printf("%s %s\n", cmd, strings.Join(args, " "))

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
