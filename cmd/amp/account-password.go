package main

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
)

// PwdCmd is the main command for attaching password sub-commands.
var (
	pwdChangeCmd = &cobra.Command{
		Use:     "change",
		Short:   "Change password",
		Example: "amp password change --name=jdoe --password=p@s5wrd --new-password=v@larm0rghuli$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd)
		},
	}

	pwdResetCmd = &cobra.Command{
		Use:     "reset",
		Short:   "Reset password",
		Example: "amp password reset --name=jdoe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd)
		},
	}

	pwdSetCmd = &cobra.Command{
		Use:     "set",
		Short:   "Set password",
		Example: "amp password set --token=this-is-a-token-sample --password=v@lard0haeri$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdSet(AMP, cmd)
		},
	}

	newPwd string
)

func init() {
	PwdCmd.AddCommand(pwdChangeCmd)
	PwdCmd.AddCommand(pwdResetCmd)
	PwdCmd.AddCommand(pwdSetCmd)
	pwdSetCmd.Flags().StringVar(&token, "token", token, "Verification Token")
	pwdSetCmd.Flags().StringVar(&password, "password", password, "Password")

	pwdChangeCmd.Flags().StringVar(&username, "name", username, "Account Name")
	pwdChangeCmd.Flags().StringVar(&password, "password", password, "Current Password")
	pwdChangeCmd.Flags().StringVar(&newPwd, "new-password", newPwd, "New Password")

	pwdResetCmd.Flags().StringVar(&username, "name", username, "Account Name")
}

// pwdReset validates the input command line arguments and resets password of an account
// by invoking the corresponding rpc/storage method
func pwdReset(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = getName()
	}

	request := &account.PasswordResetRequest{
		Name: username,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordReset(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Hi %s! Please check your email to complete the password reset process.", username)
	return nil
}

// pwdChange validates the input command line arguments and changes existing password of an account
// by invoking the corresponding rpc/storage method
func pwdChange(amp *cli.AMP, cmd *cobra.Command) (err error) {
	fmt.Println("Enter your current password.")
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = getPassword()
	}
	fmt.Println("Enter new password.")
	if cmd.Flag("new-password").Changed {
		newPwd = cmd.Flag("new-password").Value.String()
	} else {
		newPwd = getPassword()
	}

	request := &account.PasswordChangeRequest{
		ExistingPassword: password,
		NewPassword:      newPwd,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordChange(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Your password change has been successful.")
	return nil
}

// pwdSet validates the input command line arguments and sets password of an account
// by invoking the corresponding rpc/storage method
func pwdSet(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("token").Changed {
		token = cmd.Flag("token").Value.String()
	} else {
		token = getToken()
	}
	fmt.Println("Enter new password.")
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = getPassword()
	}

	request := &account.PasswordSetRequest{
		Token:    token,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordSet(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Your password set has been successful.")
	return nil
}
