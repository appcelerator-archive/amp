package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Cobra definitions for account management related commands
var (
	signUpCmd = &cobra.Command{
		Use:   "signup",
		Short: "Signup for a new account",
		Long:  `The signup command creates a new account and sends a verification link to the registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP, cmd)
		},
	}

	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify account",
		Long:  `The verify command verifies an account by checking the given verification code.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP, cmd)
		},
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to account",
		Long:  `The login command logs the user into their existing account.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(AMP, cmd)
		},
	}

	forgotLoginCmd = &cobra.Command{
		Use:   "forgot-login",
		Short: "Retrieve account name",
		Long:  `The forgot login command retrieves the account name, in case the user has forgotten it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(AMP, cmd)
		},
	}

	pwdCmd = &cobra.Command{
		Use:   "password",
		Short: "Account password operations",
		Long:  "The password command allows users to reset or update password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwd(AMP, cmd)
		},
	}

	switchCmd = &cobra.Command{
		Use:   "switch",
		Short: "Switch account",
		Long:  "The switch command allows users to switch between their personal and organization accounts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchAccount(AMP, cmd)
		},
	}

	change bool
	reset  bool
	set    bool

	username string
	email    string
	password string
	token    string
	newPwd   string

	//TODO: pass verbose as arg
	manager = NewCmdManager("")
)

//Adding account management commands to the account command
func init() {
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(loginCmd)
	AccountCmd.AddCommand(forgotLoginCmd)
	AccountCmd.AddCommand(pwdCmd)
	AccountCmd.AddCommand(switchCmd)

	signUpCmd.Flags().StringVar(&username, "name", username, "Account Name")
	signUpCmd.Flags().StringVar(&email, "email", email, "Email ID")
	signUpCmd.Flags().StringVar(&password, "password", password, "Password")

	verifyCmd.Flags().StringVar(&token, "token", token, "Verification Token")

	loginCmd.Flags().StringVar(&username, "name", username, "Account Name")
	loginCmd.Flags().StringVar(&password, "password", password, "Password")

	forgotLoginCmd.Flags().StringVar(&email, "email", email, "Email ID")

	pwdCmd.Flags().BoolVar(&change, "change", false, "Change Password")
	pwdCmd.Flags().BoolVar(&reset, "reset", false, "Reset Password")
	pwdCmd.Flags().BoolVar(&set, "set", false, "Set Password")
	pwdCmd.Flags().StringVar(&username, "name", username, "Account Name")
	pwdCmd.Flags().StringVar(&token, "token", token, "Verification Token")
	pwdCmd.Flags().StringVar(&password, "password", password, "Current Password")
	pwdCmd.Flags().StringVar(&newPwd, "new-password", newPwd, "New Password")

	switchCmd.Flags().StringVar(&username, "name", username, "Account Name")
}

// signUp validates the input command line arguments and creates a new account
// by invoking the corresponding rpc/storage method
func signUp(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = GetName()
	}
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = GetEmailAddress()
	}
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = GetPassword()
	}

	request := &account.SignUpRequest{
		Name:     username,
		Email:    email,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.SignUp(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Hi %s! Please check your email to complete the signup process.", username)
	return nil
}

// verify validates the input command line arguments and verifies an account
// by invoking the corresponding rpc/storage method
func verify(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("token").Changed {
		token = cmd.Flag("token").Value.String()
	} else {
		token = GetToken()
	}

	request := &account.VerificationRequest{
		Token: token,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.Verify(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Your account has now been activated.")
	return nil
}

// login validates the input command line arguments and allows login to an existing account
// by invoking the corresponding rpc/storage method
func login(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = GetName()
	}
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = GetPassword()
	}

	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	header := metadata.MD{}
	_, err = accClient.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	if err := cli.SaveToken(header); err != nil {
		return err
	}
	manager.printf(colSuccess, "Welcome back, %s!", username)
	return nil
}

// forgotLogin validates the input command line arguments and retrieves account name
// by invoking the corresponding rpc/storage method
func forgotLogin(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = GetEmailAddress()
	}

	request := &account.ForgotLoginRequest{
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.ForgotLogin(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Your login name has been sent to the address: %s", email)
	return nil
}

func switchAccount(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("account: ")
		username = GetName()
	}

	request := &account.SwitchRequest{
		Account: username,
	}
	accClient := account.NewAccountClient(amp.Conn)
	header := metadata.MD{}
	_, err = accClient.Switch(context.Background(), request, grpc.Header(&header))
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	if err := cli.SaveToken(header); err != nil {
		return err
	}
	manager.printf(colSuccess, "Your are now logged in as: %s", username)
	return nil
}

// pwd validates the input command line arguments and performs password-related operations
// by invoking the corresponding rpc/storage method
func pwd(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if reset {
		return pwdReset(amp, cmd)
	}
	if change {
		return pwdChange(amp, cmd)
	}
	if set {
		return pwdSet(amp, cmd)
	}
	manager.printf(colWarn, "Choose a command for password operation.\nUse amp account password -h for help.")
	return nil
}

// pwdReset validates the input command line arguments and resets password of an account
// by invoking the corresponding rpc/storage method
func pwdReset(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = GetName()
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
		password = GetPassword()
	}
	fmt.Println("Enter new password.")
	if cmd.Flag("new-password").Changed {
		newPwd = cmd.Flag("new-password").Value.String()
	} else {
		newPwd = GetPassword()
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

// pwdSet validates the input command line arguments and changes existing password of an account
// by invoking the corresponding rpc/storage method
func pwdSet(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("token").Changed {
		token = cmd.Flag("token").Value.String()
	} else {
		token = GetToken()
	}
	fmt.Println("Enter new password.")
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = GetPassword()
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
