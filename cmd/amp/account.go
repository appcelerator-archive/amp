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
)

func init() {
	RootCmd.AddCommand(AccountCmd)
}
