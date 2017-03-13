package main

import (
	"github.com/spf13/cobra"
)

// AccountCmd is the main command for attaching account sub-commands.
var (
	OrgCmd = &cobra.Command{
		Use:   "org",
		Short: "Organization operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	TeamCmd = &cobra.Command{
		Use:   "team",
		Short: "Team operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	UserCmd = &cobra.Command{
		Use:   "user",
		Short: "User operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	PwdCmd = &cobra.Command{
		Use:   "password",
		Short: "Password operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}
)

func init() {
	RootCmd.AddCommand(OrgCmd)
	RootCmd.AddCommand(TeamCmd)
	RootCmd.AddCommand(UserCmd)
	RootCmd.AddCommand(PwdCmd)
}
