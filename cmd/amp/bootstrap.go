package main

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/appcelerator/amp/api/client"
	"github.com/spf13/cobra"
)

var (
	startArgs = [...]string{"-p", "docker"}
	stopArgs  = [...]string{"-c"}
)

// StartCmd will bootstrap a new cluster
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a local amp cluster",
	Long:  `The start command initializes a local amp cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start(AMP, cmd, args)
	},
}

// StopCmd will stop and cleanup a managed cluster
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a local amp cluster",
	Long:  `The stop command stops and cleans up a local amp cluster.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stop(AMP, cmd, args)
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(stopCmd)
}

func start(amp *client.AMP, cmd *cobra.Command, args []string) error {
	return updateCluster(append(startArgs[:], args[:]...))
}
func stop(amp *client.AMP, cmd *cobra.Command, args []string) error {
	return updateCluster(append(stopArgs[:], args[:]...))
}

// TODO: replace the bootstrap script with go code
func updateCluster(args []string) error {
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
