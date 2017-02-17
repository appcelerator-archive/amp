package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

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
		Short: "Get username for an account",
		Long:  `The forgot login command retrieves the account name, in case the user has forgotten it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(AMP)
		},
	}
)

func init() {
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(loginCmd)
	AccountCmd.AddCommand(forgotLoginCmd)
}

func signUp(amp *client.AMP) error {
	fmt.Println("This will sign you up for a new personal AMP account.")
	username := getUserName()
	email, err := getEmailAddress()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.SignUpRequest{
		Name:  username,
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process.")
	return nil
}

func verify(amp *client.AMP) error {
	fmt.Println("This will sign you up for a new personal AMP account.")
	token := getToken()
	request := &account.VerificationRequest{
		Token: token,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.Verify(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	fmt.Println("Your account has now been activated.")
	return nil
}

func login(amp *client.AMP) (err error) {
	fmt.Println("This will login to an existing AMP account.")
	username := getUserName()
	password, err := getPassword()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.Login(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err.Error())
	}
	fmt.Println("Welcome back, ", username, "!")
	return nil
}

func forgotLogin(amp *client.AMP) error {
	fmt.Println("This will send your username to your registered email address")
	email, err := getEmailAddress()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.ForgotLoginRequest{
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.ForgotLogin(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	fmt.Println("Your login name has been sent to the address: ", email)
	return nil
}

func getUserName() (username string) {
	fmt.Print("username: ")
	fmt.Scanln(&username)
	err := account.CheckUserName(username)
	if err != nil {
		fmt.Println("Username is mandatory. Try again!")
		fmt.Println("")
		return getUserName()
	}
	return
}

func getEmailAddress() (email string, err error) {
	fmt.Print("email: ")
	fmt.Scanln(&email)
	email, err = account.CheckEmailAddress(email)
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

func getPassword() (password string, err error) {
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
	} else {
		return
	}
	return
}
