package cluster

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"os/exec"
)

var (
	err error
	cmdOut []byte
)
// NewStatusCommand returns a new instance of the status command for by providing groups and instances of local cluster.
func NewStatusCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Retrieve details about a local amp cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(c, args)
		},
	}
}

func status(c cli.Interface, args []string) error {
	cmdName := "docker"
	args = []string{"container", "ls"}
	if cmdOut, err = exec.Command(cmdName, args...).Output(); err != nil {
		c.Console().Fatalf("error while executing command: %v\n", err)
	}
	status := string(cmdOut)
	c.Console().Println(status)
	return err
}

