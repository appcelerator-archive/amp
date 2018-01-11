package image

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewRegistryCommand returns a new instance of the registry command.
func NewImageCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "image",
		Short:   "Image management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}

	cmd.AddCommand(NewPushCommand(c))
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))

	return cmd
}
