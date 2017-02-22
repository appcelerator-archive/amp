package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

// Cobra definitions for account management related commands
var (
	signUpCmd = &cobra.Command{
		Use:   "signup",
		Short: "Signup for a new account",
		Long:  `The signup command creates a new account and sends a verification link to the registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP)
		},
	}

	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify account",
		Long:  `The verify command verifies an account by sending a verification code to their registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP)
		},
	}

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to account",
		Long:  `The login command logs the user into their existing account.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(AMP)
		},
	}

	forgotLoginCmd = &cobra.Command{
		Use:   "forgot-login",
		Short: "Retrieve account name",
		Long:  `The forgot login command retrieves the account name, in case the user has forgotten it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(AMP)
		},
	}

	pwdCmd = &cobra.Command{
		Use:   "password",
		Short: "Account password operations",
		Long:  "The password command allows users allows users to reset or update password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwd(AMP, cmd, args)
		},
	}

	change bool
	reset  bool
	set    bool

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

	pwdCmd.Flags().BoolVar(&change, "change", false, "Change Password")
	pwdCmd.Flags().BoolVar(&reset, "reset", false, "Reset Password")
	pwdCmd.Flags().BoolVar(&set, "set", false, "Set Password")
}

// signUp validates the input command line arguments and creates a new account
// by invoking the corresponding rpc/storage method
func signUp(amp *cli.AMP) (err error) {
	manager.printf(0, "This will sign you up for a new personal AMP account.\n")
	username := getUserName()
	email := getEmailAddress()
	password := getPassword()
	request := &account.SignUpRequest{
		Name:     username,
		Email:    email,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.SignUp(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Hi %s! Please check your email to complete the signup process.\n", username)
	return nil
}

// verify validates the input command line arguments and verifies an account
// by invoking the corresponding rpc/storage method
func verify(amp *cli.AMP) (err error) {
	manager.printf(0, "This will verify an existing AMP account.\n")
	token := getToken()
	request := &account.VerificationRequest{
		Token: token,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.Verify(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Your account has now been activated.\n")
	return nil
}

// login validates the input command line arguments and allows login to an existing account
// by invoking the corresponding rpc/storage method
func login(amp *cli.AMP) (err error) {
	manager.printf(0, "This will login to an existing AMP account.\n")
	username := getUserName()
	password := getPassword()
	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	header := metadata.MD{}
	_, err = accClient.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	if err := cli.SaveToken(header); err != nil {
		return err
	}
	manager.printf(4, "Welcome back, %s!\n", username)
	return nil
}

// forgotLogin validates the input command line arguments and retrieves account name
// by invoking the corresponding rpc/storage method
func forgotLogin(amp *cli.AMP) (err error) {
	manager.printf(0, "This will send your username to your registered email address.\n")
	email := getEmailAddress()
	request := &account.ForgotLoginRequest{
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.ForgotLogin(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Your login name has been sent to the address: %s\n", email)
	return nil
}

// pwd validates the input command line arguments and performs password-related operations
// by invoking the corresponding rpc/storage method
func pwd(amp *cli.AMP, cmd *cobra.Command, args []string) (err error) {
	if reset {
		return pwdReset(amp, cmd, args)
	}
	if change {
		return pwdChange(amp, cmd, args)
	}
	if set {
		return pwdSet(amp, cmd, args)
	}
	manager.printf(2, "Choose a command for password operation.\nUse amp account password -h for help.\n")
	return nil
}

// pwdReset validates the input command line arguments and resets password of an account
// by invoking the corresponding rpc/storage method
func pwdReset(amp *cli.AMP, cmd *cobra.Command, args []string) (err error) {
	manager.printf(0, "This will send a password reset email to your email address.\n")
	username := getUserName()
	request := &account.PasswordResetRequest{
		Name: username,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordReset(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Hi %s! Please check your email to complete the password reset process.\n", username)
	return nil
}

// pwdChange validates the input command line arguments and changes existing password of an account
// by invoking the corresponding rpc/storage method
func pwdChange(amp *cli.AMP, cmd *cobra.Command, args []string) (err error) {
	// Get inputs
	manager.printf(0, "This will allow you to update your existing password.\n")
	fmt.Println("Enter your current password.")
	existingPwd := getPassword()
	fmt.Println("Enter new password.")
	newPwd := getPassword()
	request := &account.PasswordChangeRequest{
		ExistingPassword: existingPwd,
		NewPassword:      newPwd,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordChange(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Your password change has been successful.\n")
	return nil
}

// pwdSet validates the input command line arguments and changes existing password of an account
// by invoking the corresponding rpc/storage method
func pwdSet(amp *cli.AMP, cmd *cobra.Command, args []string) (err error) {
	// Get inputs
	manager.printf(0, "This will allow you to set a new password.\n")
	token := getToken()
	fmt.Println("Enter new password.")
	password := getPassword()
	request := &account.PasswordSetRequest{
		Token:    token,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordSet(context.Background(), request)
	if err != nil {
		manager.fatalf("server error: %v\n", grpc.ErrorDesc(err))
		return
	}
	manager.printf(4, "Your password change has been successful.\n")
	return nil
}

func getUserName() (username string) {
	fmt.Print("username: ")
	fmt.Scanln(&username)
	username = strings.TrimSpace(username)
	err := schema.CheckName(username)
	if err != nil {
		manager.printf(2, "Username is invalid/incorrect. Try again!\n\n")
		return getUserName()
	}
	return
}

func getEmailAddress() (email string) {
	fmt.Print("email: ")
	fmt.Scanln(&email)
	email = strings.TrimSpace(email)
	_, err := schema.CheckEmailAddress(email)
	if err != nil {
		manager.printf(2, "Email in incorrect format. Try again!\n\n")
		return getEmailAddress()
	}
	return
}

func getToken() (token string) {
	fmt.Print("token: ")
	fmt.Scanln(&token)
	token = strings.TrimSpace(token)
	return
}

func getPassword() (password string) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if pw == nil || err != nil {
		manager.printf(2, "Password is mandatory. Try again!\n\n")
		return getPassword()
	}
	password = string(pw)
	password = strings.TrimSpace(password)
	err = schema.CheckPassword(password)
	if err != nil {
		manager.printf(2, "Password entered is too weak. Password must be at least 8 characters long. Try again!\n\n")
		return getPassword()
	}
	return
}
