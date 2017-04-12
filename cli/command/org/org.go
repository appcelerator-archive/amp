package org

import (
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/cli/command/org/member"
	"github.com/spf13/cobra"
)

// NewOrgCommand returns a new instance of the org command.
func NewOrgCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "org",
		Short:   "Organization management operations",
		PreRunE: cli.NoArgs,
		RunE:    c.ShowHelp,
	}
	cmd.AddCommand(NewOrgListCommand(c))
	cmd.AddCommand(NewOrgCreateCommand(c))
	cmd.AddCommand(NewOrgRemoveCommand(c))
	cmd.AddCommand(NewOrgGetCommand(c))
	cmd.AddCommand(member.NewOrgMemCommand(c))
	cmd.AddCommand(NewSwitchCommand(c))
	return cmd
}
