package main

import (
	"log"

	"github.com/spf13/cobra"
)

func NewInstallCommand() *cobra.Command {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Set up amp services in swarm environment",
		Run:   install,
	}
	return installCmd
}

func install(cmd *cobra.Command, args []string) {
	log.Println("tbd")
}
