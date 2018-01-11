package cmd

import (
	ampdocker "github.com/appcelerator/amp/pkg/docker"
	"github.com/spf13/cobra"
)

var Docker *ampdocker.Docker

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "ampctl",
		PersistentPreRunE: Init,
		RunE:              Root,
		Short:             "Run commands in target amp cluster",
	}

	return rootCmd
}
func Init(cmd *cobra.Command, args []string) error {
	Docker = ampdocker.NewEnvClient()
	return Docker.Connect()
}

func Root(cmd *cobra.Command, args []string) error {
	// perform checks and install by default when no sub-command is specified
	if err := Checks(cmd, []string{}); err != nil {
		return err
	}
	return Install(cmd, args)
}
