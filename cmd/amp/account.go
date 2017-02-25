package main

import (
	"github.com/spf13/cobra"
)

// AccountCmd is the main command for attaching account sub-commands.
var (
	AccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Account operations",
		Long:  `The account command manages all account-related operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	OrgCmd = &cobra.Command{
		Use:   "org",
		Short: "Organization operations",
		Long:  `The organization command manages all organization-related operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	TeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team operations",
		Long:  `The team command manages all team-related operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}
)

func init() {
	RootCmd.AddCommand(AccountCmd)
	RootCmd.AddCommand(OrgCmd)
	RootCmd.AddCommand(TeamCmd)
}
