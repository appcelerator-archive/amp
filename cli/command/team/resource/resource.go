package resource

import (
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
)

// NewTeamResourceCommand returns a new instance of the team resource command.
func NewTeamResourceCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resource",
		Short:   "Manage team resources",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewAddTeamResCommand(c))
	cmd.AddCommand(NewRemoveTeamResCommand(c))
	cmd.AddCommand(NewListTeamResCommand(c))
	cmd.AddCommand(NewChangeTeamResPermissionLevelCommand(c))
	return cmd
}
