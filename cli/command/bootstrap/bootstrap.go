package bootstrap

import (
	"bufio"
	"os/exec"

	"github.com/appcelerator/amp/cli"
)

// TODO: replace the bootstrap script with go code
func updateCluster(c cli.Interface, args []string) error {
	// TODO: use AMPHOME environment variable for path
	cmd := "bootstrap"
	proc := exec.Command(cmd, args...)
	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}

	outscanner := bufio.NewScanner(stdout)
	go func() {
		for outscanner.Scan() {
			c.Console().Printf("%s\n", outscanner.Text())
		}
	}()
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			c.Console().Printf("%s\n", errscanner.Text())
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
		return err

	}

	return nil
}
