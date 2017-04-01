package exec

import (
	"bufio"
	"os/exec"

	"github.com/appcelerator/amp/cli"
)

// Run is a helper to exec an os command using the streams configured for the cli.
func Run(c cli.Interface, name string, args []string) error {
	proc := exec.Command(name, args...)
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
			c.Console().Println(outscanner.Text())
		}
	}()
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
		// Just pass along the information that the process exited with a failure;
		// whatever error information it displayed is what the user will see.
		return err

	}

	return nil
}
