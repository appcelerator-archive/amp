package team

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/team/member"
	"github.com/spf13/cobra"
)

// NewTeamCommand returns a new instance of the team command.
func NewTeamCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "team",
		Short:   "Team management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewTeamCreateCommand(c))
	cmd.AddCommand(NewTeamListCommand(c))
	cmd.AddCommand(NewTeamRemoveCommand(c))
	cmd.AddCommand(NewTeamGetCommand(c))
	cmd.AddCommand(member.NewTeamMemberCommand(c))
	return cmd
}
