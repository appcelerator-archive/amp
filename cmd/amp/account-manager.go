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
)

func init() {
	AccountCmd.AddCommand(signUpCmd)
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
	client := account.NewAccountClient(amp.Conn)
	_, err = client.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process.")
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
