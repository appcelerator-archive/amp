package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mitchellh/go-homedir"
)

var (
	dockerCmd  = "docker"
	dockerArgs []string
)

func init() {
	homedir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	dockerArgs = []string{
		"run", "-i", "--rm", "--name", "ampcli",
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", fmt.Sprintf("%s/.config/amp:/root/.config/amp", homedir),
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

	// wire up stdin to the command's stdin
	stdin, err := proc.StdinPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		defer stdin.Close()
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			io.WriteString(stdin, input.Text()+"\n")
		}
	}()

	// display command's stdout
	stdout, err := proc.StdoutPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		output := make([]byte, 1)
		for {
			n, err := stdout.Read(output)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", string(output[0]))
		}
	}()

	// display command's stderr
	stderr, err := proc.StderrPipe()
	if err != nil {
		panic(err)
	}
	go func() {
		errscanner := bufio.NewScanner(stderr)
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
