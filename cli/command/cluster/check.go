package cluster

import (
	"errors"
	"os/exec"
)

// check if a cluster command is not already in progress, and provides hints if it's the case
func check(provider string) error {
	cmd := "docker"
	switch provider {
	case "local", "docker":
		// check that the amp-boostrap container does not exist
		args := []string{
			"container", "inspect", "amp-bootstrap",
		}
		proc := exec.Command(cmd, args...)
		err := proc.Run()
		if err == nil {
			return errors.New("A cluster operation is still in progress. Please wait for it to finish unless you want to terminate it now (docker rm -f amp-bootstrap) and try again.")
		}
	default:
		// no check
	}
	return nil
}
