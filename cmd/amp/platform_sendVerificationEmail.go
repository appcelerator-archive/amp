package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/pkg/ampmail"
	"github.com/spf13/cobra"
	"os"
)

// PlatformSendVerificationEmail to send a verification email to a user
var PlatformSendVerificationEmail = &cobra.Command{
	Use:   "sendVerificationEmail [email] [accountName] [token]",
	Short: "send a verification email",
	Long:  `Send a verification email to control user email at signUp`,
	Run: func(cmd *cobra.Command, args []string) {
		sendVerificationEmail(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformSendVerificationEmail)
}

func sendVerificationEmail(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		fmt.Printf("usage: amp sendVerificationEmail [email] [accountName] [token]\n")
		os.Exit(1)
	}
	if err := ampmail.SendAccountVerificationEmail(args[0], args[1], args[2]); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
