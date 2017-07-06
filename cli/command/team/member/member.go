package member

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewTeamMemberCommand returns a new instance of the team member command.
func NewTeamMemberCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Short:   "Manage team members",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewAddTeamMemCommand(c))
	cmd.AddCommand(NewRemoveTeamMemCommand(c))
	cmd.AddCommand(NewListTeamMemCommand(c))
	return cmd
}
