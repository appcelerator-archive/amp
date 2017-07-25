package secret

import (

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

type CreateOpts struct {
	Labels []string
}

var createOpts = &CreateOpts{
	Labels: []string{},
}

// NewCreateCommand returns a new instance of the create command for creating a secret.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create a secret from a file or STDIN as content",
		PreRunE: cli.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&createOpts.Labels, "labels", "l", []string{}, "Secret labels")

	return cmd
}

func create(c cli.Interface, cmd *cobra.Command) error {
	// TODO read from STDIN if no arg, otherwise arg[1] is the file to read
	// TODO call service to create secret
	return nil
}

