package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var (
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "log in to amp",
		Long:  `The login command logs the user into existing account`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AMP.Connect()
			if err != nil {
				return err
			}
			return login(AMP)
		},
	}

	AccountCmd = &cobra.Command{
		Use:   "account",
		Short: "Account operations",
		Long:  `The account command manages all account-related operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	signUpCmd = &cobra.Command{
		Use:   "signup",
		Short: "Create a new account and login",
		Long:  `The signup command creates a new account and sends a verification link to their registered email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP)
		},
	}

	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify email using code",
		Long: `The verify command is used to verify the users account using the code sent to them via email.
		This is used if the user cannot access the verification link sent.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP)
		},
	}

	forgotLoginCmd = &cobra.Command{
		Use:   "forgot-login",
		Short: "Get username for an account",
		Long: `The forgot login command is used when a user has forgotten their username.
		An email is sent to the email entered with the username registered to it.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(AMP)
		},
	}

	pwdResetCmd = &cobra.Command{
		Use:   "password-reset USERNAME EMAIL",
		Short: "Reset Password",
		Long:  "The password reset command allows users to reset password. A link to reset password will be sent to their registered email address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd, args)
		},
	}

	pwdChangeCmd = &cobra.Command{
		Use:   "password-change USERNAME EXISTING-PASSWORD NEW-PASSWORD CONFIRM-NEW-PASSWORD",
		Short: "Change Password",
		Long:  "The password change command allows users to reset existing password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd, args)
		},
	}

	switchRoleCmd = &cobra.Command{
		Use:   "switch [ORGANIZATION]",
		Short: "Switch primary organization",
	}

	createOrganizationCmd = &cobra.Command{
		Use:   "create organization [NAME] [EMAIL]",
		Short: "Create an organization",
	}

	listUsersCmd = &cobra.Command{
		Use:   "list users [ORGANIZATION] [TEAM]",
		Short: "list users, optionally filter by organization and team",
	}

	listOrganizationsCmd = &cobra.Command{
		Use:   "list organizations",
		Short: "list organizations",
	}

	listTeamsCmd = &cobra.Command{
		Use:   "list teams [ORGANIZATION]",
		Short: "list teams by organization",
	}

	listPermissionsCmd = &cobra.Command{
		Use:   "list permissions [ORGANIZATION] [TEAM]",
		Short: "list permissions by team",
	}

	infoCmd = &cobra.Command{
		Use:   "info [name]",
		Short: "list account information",
	}

	editCmd = &cobra.Command{
		Use:   "edit [name]",
		Short: "edits account information",
	}

	deleteCmd = &cobra.Command{
		Use:   "delete [name]",
		Short: "deletes account",
	}

	addOrganizationMembersCmd = &cobra.Command{
		Use:   "add organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Add users to an organization",
	}

	addTeamMembersCmd = &cobra.Command{
		Use:   "add team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Add users to a team",
	}

	removeOrganizationMembersCmd = &cobra.Command{
		Use:   "remove organization members [ORGANIZATION] [MEMBERS...]",
		Short: "Remove users from an organization",
	}

	removeTeamMembersCmd = &cobra.Command{
		Use:   "remove team members [ORGANIZATION] [TEAM] [MEMBERS...]",
		Short: "Remove users from a team",
	}

	grantPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM]",
		Short: "Grant permission to a team for a resource",
	}

	editPermissionCmd = &cobra.Command{
		Use:   "grant [RESOURCE_ID] [LEVEL] [ORGANIZATION] [TEAM] ",
		Short: "Edit a permission for a resource",
	}

	revokePermissionCmd = &cobra.Command{
		Use:   "revoke [RESOURCE_ID] [ORGANIZATION] [TEAM]",
		Short: "Revoke permission from a team for a resource",
	}

	transferOwnershipCmd = &cobra.Command{
		Use:   "transfer [RESOURCE_ID] [ORGANIZATION]",
		Short: "Transfer ownership of a resource to a different organization",
	}
)

func init() {
	RootCmd.AddCommand(LoginCmd)
	RootCmd.AddCommand(AccountCmd)
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(forgotLoginCmd)
	AccountCmd.AddCommand(pwdResetCmd)
	AccountCmd.AddCommand(pwdChangeCmd)
}

// login gets the username and password, validates the command line inputs
// and logs the user into their account
func login(amp *client.AMP) (err error) {
	fmt.Println("This will login an existing personal AMP account.")
	username := getUserName()
	password, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.LogInRequest{
		Name:     username,
		Password: password,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Login(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err.Error())
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Welcome back,", username)
	color.Unset()
	return nil
}

// signup signs up visitor for a new personal account.
// Sends a verification link to their email address.
func signUp(amp *client.AMP) error {
	fmt.Println("This will sign you up for a new personal AMP account.")
	username := getUserName()
	email := getEmailAddress()
	request := &account.SignUpRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err := client.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("Server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process")
	color.Unset()
	return nil
}

// verify gets the unique code sent to the visitor in the email verification, registered username and new password,
// validates the command line inputs and activates their account.
func verify(amp *client.AMP) (err error) {
	fmt.Println("This will verify your account and confirm your password")
	code := getCode()
	username := getUserName()
	password, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	request := &account.VerificationRequest{
		Name:     username,
		Password: password,
		Code:     code,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Verify(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "! Your account has now be activated")
	color.Unset()
	return nil
}

// forgotLogin validates the input command line arguments and sends a users username
// to their registered email address by invoking the corresponsind rpc/storage method
func forgotLogin(amp *client.AMP) error {
	fmt.Println("This will send your username to your registered email address")
	email := getEmailAddress()
	request := &account.ForgotLoginRequest{
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err := client.ForgotLogin(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Your login name has been sent to the address:", email)
	color.Unset()
	return nil
}

// pwdReset validates the input command line arguments and resets the current password
// by invoking the corresponding rpc/storage method
func pwdReset(amp *client.AMP, cmd *cobra.Command, args []string) error {
	fmt.Println("This will send a password reset email to your email address")
	username := getUserName()
	email := getEmailAddress()
	request := &account.PasswordResetRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err := client.PasswordReset(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "! Please check your email to complete the password reset process.")
	color.Unset()
	return nil
}

// pwdChange validates the input command line arguments and changes the current password
// by invoking the corresponding rpc/storage method
func pwdChange(amp *client.AMP, cmd *cobra.Command, args []string) error {
	fmt.Println("This will allow you to update your existing password.")
	username := getUserName()
	fmt.Print("existing ")
	existingPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	fmt.Print("new ")
	newPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	getConfirmPwd(newPwd)
	request := &account.PasswordChangeRequest{
		Name:             username,
		ExistingPassword: existingPwd,
		NewPassword:      newPwd,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.PasswordChange(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error : %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi ", username, "! Your recent password change has been successful.")
	color.Unset()
	return nil
}

func getUserName() (username string) {
	fmt.Print("username: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&username)
	color.Unset()
	err := account.CheckUserName(username)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("username is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getUserName()
	}
	return
}

func getPwd() (password string, err error) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		if err == gopass.ErrInterrupted {
			err = fmt.Errorf(err.Error())
			return
		} else {
			return
		}
	}
	password = string(pw)
	err = account.CheckPassword(password)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("password is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getPwd()
	}
	err = account.CheckPasswordStrength(password)
	if err != nil {
		if strings.Contains(err.Error(), "password too weak") {
			color.Set(color.FgRed, color.Bold)
			fmt.Println("password entered is too weak. password must be at least 8 characters long. Try again!")
			color.Unset()
			fmt.Println("")
			return getPwd()
		} else {
			return
		}
	}
	return
}

func getConfirmPwd(newPwd string) error {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Enter Password again for confirmation.")
	color.Unset()
	fmt.Print("confirm ")
	confirmNewPwd, err := getPwd()
	if err != nil {
		return err
	}
	if confirmNewPwd != newPwd {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Password mismatch. Try again!")
		color.Unset()
		fmt.Println("")
		getConfirmPwd(newPwd)
	}
	return nil
}

func getCode() (code string) {
	fmt.Print("Code: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&code)
	color.Unset()
	err := account.CheckVerificationCodeFormat(code)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("code is invalid. Code must be 8 characters long. Try again!")
		color.Unset()
		fmt.Println("")
		return getCode()
	}
	return
}

func getEmailAddress() (email string) {
	fmt.Print("email: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&email)
	color.Unset()
	email, err := account.CheckEmailAddress(email)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("email in incorrect format. Try again!")
		color.Unset()
		fmt.Println("")
		return getEmailAddress()
	}
	return
}
