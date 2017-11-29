package object_store

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewObjectCommand returns a new instance of the object-store command.
func NewObjectStoreCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "object-store",
		Short:   "Object storage operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewListCommand(c))
	cmd.AddCommand(NewCreateCommand(c))
	cmd.AddCommand(NewForgetCommand(c))
	cmd.AddCommand(NewRemoveCommand(c))
	return cmd
}
