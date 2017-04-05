package member

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewOrgMemCommand returns a new instance of the organization member command.
func NewOrgMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Short:   "Manage organization members",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewOrgAddMemCommand(c))
	cmd.AddCommand(NewOrgRemoveMemCommand(c))
	cmd.AddCommand(NewOrgListMemCommand(c))
	cmd.AddCommand(NewOrgChangeMemRoleCommand(c))
	return cmd
}
