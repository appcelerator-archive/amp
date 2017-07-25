package secret

import (

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewRemoveCommand returns a new instance of the remove command for removing one or more secrets
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [OPTIONS]",
		Short:   "Remove one or more secrets",
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, cmd)
		},
	}

	return cmd
}

func remove(c cli.Interface, cmd *cobra.Command) error {
	// TODO call service to remove one or more secrets
	return nil
}

