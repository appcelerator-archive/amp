package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/data/account/schema"
	"github.com/appcelerator/amp/pkg/auth"
	"github.com/howeyc/gopass"
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

	reset  bool
	change bool
)

//Adding account management commands to the account command
func init() {
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(loginCmd)
	AccountCmd.AddCommand(forgotLoginCmd)
	AccountCmd.AddCommand(pwdCmd)

	pwdCmd.Flags().BoolVar(&reset, "reset", false, "Reset Password")
	pwdCmd.Flags().BoolVar(&change, "change", false, "Change Password")

}

// signUp validates the input command line arguments and creates a new account
// by invoking the corresponding rpc/storage method
func signUp(amp *client.AMP) (err error) {
	fmt.Println("This will sign you up for a new personal AMP account.")
	username := getUserName()
	email := getEmailAddress()
	password := getPassword()
	request := &account.SignUpRequest{
		Name:     username,
		Email:    email,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process.")
	fmt.Println("token", reply.Token)
	return nil
}

// verify validates the input command line arguments and verifies an account
// by invoking the corresponding rpc/storage method
func verify(amp *client.AMP) (err error) {
	fmt.Println("This will sign you up for a new personal AMP account.")
	token := getToken()
	request := &account.VerificationRequest{
		Token: token,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.Verify(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	fmt.Println("Your account has now been activated.")
	return nil
}

// login validates the input command line arguments and allows login to an existing account
// by invoking the corresponding rpc/storage method
func login(amp *client.AMP) (err error) {
	fmt.Println("This will login to an existing AMP account.")
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
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	if err := SaveToken(header); err != nil {
		return err
	}
	fmt.Println("Welcome back, ", username, "!")
	return nil
}

// forgotLogin validates the input command line arguments and retrieves account name
// by invoking the corresponding rpc/storage method
func forgotLogin(amp *client.AMP) (err error) {
	fmt.Println("This will send your username to your registered email address")
	email := getEmailAddress()
	request := &account.ForgotLoginRequest{
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.ForgotLogin(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	fmt.Println("Your login name has been sent to the address: ", email)

	return nil
}

// pwd validates the input command line arguments and performs password-related operations
// by invoking the corresponding rpc/storage method
func pwd(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	if reset {
		return pwdReset(amp, cmd, args)
	}
	if change {
		return pwdChange(amp, cmd, args)
	}
	fmt.Println("Choose a command for password operation")
	fmt.Println("Use amp account password -h for help")
	return nil
}

// pwdReset validates the input command line arguments and resets password of an account
// by invoking the corresponding rpc/storage method
func pwdReset(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will send a password reset email to your email address.")
	username := getUserName()
	request := &account.PasswordResetRequest{
		Name: username,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordReset(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	fmt.Println("Hi", username, "! Please check your email to complete the password reset process.")
	return nil
}

// pwdChange validates the input command line arguments and changes existing password of an account
// by invoking the corresponding rpc/storage method
func pwdChange(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	token, err := ReadToken()
	if err != nil {
		return err
	}
	// Get inputs
	fmt.Println("This will allow you to update your existing password.")
	username := getUserName()
	fmt.Println("Enter your current password.")
	existingPwd := getPassword()
	fmt.Println("Enter new password.")
	newPwd := getPassword()

	// Call the backend
	// Set the authN token on the request header
	md := metadata.Pairs(auth.TokenKey, token)
	ctx := metadata.NewContext(context.Background(), md)
	request := &account.PasswordChangeRequest{
		Name:             username,
		ExistingPassword: existingPwd,
		NewPassword:      newPwd,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.PasswordChange(ctx, request)
	if err != nil {
		return fmt.Errorf("server error: %v", grpc.ErrorDesc(err))
	}
	fmt.Println("Hi ", username, "! Your recent password change has been successful.")
	return nil
}

func getUserName() (username string) {
	fmt.Print("username: ")
	fmt.Scanln(&username)
	err := schema.CheckName(username)
	if err != nil {
		fmt.Println("Username is mandatory. Try again!")
		fmt.Println("")
		return getUserName()
	}
	return
}

func getEmailAddress() (email string) {
	fmt.Print("email: ")
	fmt.Scanln(&email)
	_, err := schema.CheckEmailAddress(email)
	if err != nil {
		fmt.Println("Email in incorrect format. Try again!")
		fmt.Println("")
		return getEmailAddress()
	}
	return
}

func getToken() (token string) {
	fmt.Print("token: ")
	fmt.Scanln(&token)
	err := account.CheckVerificationCode(token)
	if err != nil {
		fmt.Println("Code is invalid. Try again!")
		fmt.Println("")
		return getToken()
	}
	return
}

func getPassword() (password string) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if pw == nil || err != nil {
		fmt.Println("Password is mandatory. Try again!")
		fmt.Println("")
		return getPassword()
	}
	password = string(pw)
	err = account.CheckPassword(password)
	if err != nil {
		fmt.Println("Password entered is too weak. Password must be at least 8 characters long. Try again!")
		fmt.Println("")
		return getPassword()
	}
	return
}
