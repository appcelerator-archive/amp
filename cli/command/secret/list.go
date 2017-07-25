package secret

import (

	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

type ListOpts struct {
	Filters []string
	Format string
	Quiet bool
}

var listOpts = &ListOpts{
	Filters: []string{},
}

// NewListCommand returns a new instance of the list command for listing secrets
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [OPTIONS]",
		Short:   "List secrets",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c, cmd)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.Filters, "filter", "f", []string{}, "Filter output based on conditions provided")
	flags.StringVar(&listOpts.Format, "format", "", "Pretty-print secrets using a Go template")
	flags.BoolVarP(&listOpts.Quiet, "quiet", "q", false, "Only display IDs")

	return cmd
}

func list(c cli.Interface, cmd *cobra.Command) error {
	// TODO call service to get secrets list
	return nil
}

