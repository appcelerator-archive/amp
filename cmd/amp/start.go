package main

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/appcelerator/amp/api/client"
	"github.com/spf13/cobra"
)

// StartCmd is the main command for attaching local swarm commands.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a local amp cluster",
	Long:  `The start command initializes a local amp cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start(AMP, cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
}

func start(amp *client.AMP, cmd *cobra.Command, args []string) error {
	return startCluster()
}

// TODO: replace the stacks/local-bootstrap script with go code
func startCluster() error {
	// TODO: use AMPHOME environment variable for path
	cmd := "local-bootstrap"
	proc := exec.Command(cmd)
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
			fmt.Printf("%s\n", outscanner.Text())
		}
	}()
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for errscanner.Scan() {
			fmt.Printf("%s\n", errscanner.Text())
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
