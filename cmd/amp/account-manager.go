package main

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/spf13/cobra"
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
		Long:  `The verify command creates a new account and sends a verification link to the registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP)
		},
	}
)

func init() {
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
}

// signup signs up visitor for a new personal account.
// Sends a verification link to their email address.
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

// verify gets the unique code sent to the visitor in the email verification, registered username and new password,
// validates the command line inputs and activates their account.
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
		fmt.Errorf("Code is invalid. Try again!")
		fmt.Println("")
		return getToken()
	}
	return
}
