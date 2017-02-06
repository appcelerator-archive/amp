package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	LoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Log in to AMP",
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
		Use:   "signup ACCOUNT-NAME",
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

	pwdResetCmd = &cobra.Command{
		Use:   "password-reset ACCOUNT-NAME EMAIL",
		Short: "Reset password",
		Long:  "The password reset command allows users to reset password. A link to reset password will be sent to their registered email address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd, args)
		},
	}

	pwdChangeCmd = &cobra.Command{
		Use:   "password-change ACCOUNT-NAME EXISTING-PASSWORD NEW-PASSWORD CONFIRM-NEW-PASSWORD",
		Short: "Change password",
		Long:  "The password change command allows users to reset existing password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd, args)
		},
	}

	switchRoleCmd = &cobra.Command{
		Use:   "switch ORGANIZATION-NAME",
		Short: "Switch primary organization",
		Long:  `The switch command changes the current login from a user account to the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchRole(AMP, cmd, args)
		},
	}

	createOrganizationCmd = &cobra.Command{
		Use:   "create-organization ORGANIZATION-NAME EMAIL",
		Short: "Create an organization",
		Long:  `The create organization command creates an organization with a name and email address.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrganization(AMP, cmd, args)
		},
	}

	editOrganizationCmd = &cobra.Command{
		Use:   "edit-organization ORGANIZATION-NAME",
		Short: "Edit an organization",
		Long:  `The edit organization command updates an existing organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editOrganization(AMP, cmd, args)
		},
	}

	deleteOrganizationCmd = &cobra.Command{
		Use:   "delete-organization ORGANIZATION-NAME",
		Short: "Delete an organization",
		Long:  `The delete organization command deletes an existing organization and all related information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteOrganization(AMP, cmd, args)
		},
	}

	listOrganizationsCmd = &cobra.Command{
		Use:   "list-organization",
		Short: "List organizations",
		Long:  `The list organization command displays all the available organizations.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrganization(AMP, cmd, args)
		},
	}

	addOrganizationMembersCmd = &cobra.Command{
		Use:   "add-organization-member ORGANIZATION-NAME MEMBERS...",
		Short: "Add users to an organization",
		Long:  `The add-organization command allows an owner team member to add new members in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMember(AMP, cmd, args)
		},
	}

	removeOrganizationMembersCmd = &cobra.Command{
		Use:   "remove-organization-member ORGANIZATION-NAME MEMBERS...",
		Short: "Remove users from an organization",
		Long:  `The remove-organization command allows an owner team member to remove existing members from an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrgMember(AMP, cmd, args)
		},
	}

	transferOwnershipCmd = &cobra.Command{
		Use:   "transfer-ownership RESOURCE-ID ORGANIZATION-NAME",
		Short: "Transfer ownership of a resource to a different organization",
		Long:  `The transfer command allows a resource owner to transfer a particular resource to a different organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return transferOwnership(AMP, cmd, args)
		},
	}

	createTeamCmd = &cobra.Command{
		Use:   "create-team ORGANIZATION-NAME TEAM-NAME          ",
		Short: "Create a team within an organization",
		Long:  `The create team command creates a team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTeam(AMP, cmd, args)
		},
	}

	editTeamCmd = &cobra.Command{
		Use:   "edit-team ORGANIZATION-NAME TEAM-NAME          ",
		Short: "Edit team information",
		Long:  `The edit team command updates information of team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editTeam(AMP, cmd, args)
		},
	}

	deleteTeamCmd = &cobra.Command{
		Use:   "delete-team ORGANIZATION-NAME TEAM-NAME          ",
		Short: "Delete a team within an organization",
		Long:  `The delete team command deletes a team within the specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTeam(AMP, cmd, args)
		},
	}

	listTeamsCmd = &cobra.Command{
		Use:   "list-team ORGANIZATION-NAME",
		Short: "List teams by organization",
		Long:  `The list team command displays the available teams in a specified organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTeam(AMP, cmd, args)
		},
	}

	addTeamMembersCmd = &cobra.Command{
		Use:   "add-team-member ORGANIZATION-NAME TEAM-NAME MEMBERS...",
		Short: "Add users to a team",
		Long:  `The add-team command allows an owner team member to add new members to a team in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTeamMember(AMP, cmd, args)
		},
	}

	removeTeamMembersCmd = &cobra.Command{
		Use:   "remove-team-member ORGANIZATION-NAME TEAM-NAME MEMBERS...",
		Short: "Remove users from a team",
		Long:  `The remove-team command allows an owner team member to remove existing members from a team in the organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTeamMember(AMP, cmd, args)
		},
	}

	infoCmd = &cobra.Command{
		Use:   "info ACCOUNT-NAME",
		Short: "Display account information",
		Long: `The info command displays information about the specified account name.
	If the input account name belongs to the user who is currently logged-in, the following information is displayed :
	Account Name, Email, Organization Name, Team Name, Billing Information, Other Settings

	If the input account name belongs to a different user, only the following information can be viewed :
	Account Name, Email, Organization Name, Team Name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAccount(AMP, cmd, args)
		},
	}

	editCmd = &cobra.Command{
		Use:   "edit ACCOUNT-NAME",
		Short: "Edit account information",
		Long:  `The update command allows an account owner to modify the specified account information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editAccount(AMP, cmd, args)
		},
	}

	deleteCmd = &cobra.Command{
		Use:   "delete ACCOUNT-NAME",
		Short: "Delete an account",
		Long:  `The delete command allows an account owner to delete the specified account.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteAccount(AMP, cmd, args)
		},
	}

	//TODO
	listUsersCmd = &cobra.Command{
		Use:   "list-user [ORGANIZATION-NAME] [TEAM-NAME]",
		Short: "List users, optionally filter by organization and team",
		Long:  `The list user command displays information about all users currently on the system, which can be filtered by team or organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(AMP, cmd, args)
		},
	}

	grantPermissionCmd = &cobra.Command{
		Use:   "grant-permission RESOURCE-ID PERMISSION-LEVEL ORGANIZATION-NAME TEAM-NAME",
		Short: "Grant permission to a team for a resource",
		Long: `The grant command permits an account owner to grant permissions to a team having access to specified resource.
The permissions levels can be Read, Write, Read/Write and Delete.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return grantPermission(AMP, cmd, args)
		},
	}

	editPermissionCmd = &cobra.Command{
		Use:   "edit-permission RESOURCE-ID PERMISSION-LEVEL ORGANIZATION-NAME TEAM-NAME",
		Short: "Edit permission of a team for a resource",
		Long: `The edit command permits an account owner to edit permissions of a team having access to specified resource.
The permissions levels can be Read, Write, Read/Write and Delete.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return editPermission(AMP, cmd, args)
		},
	}

	revokePermissionCmd = &cobra.Command{
		Use:   "revoke-permission RESOURCE-ID ORGANIZATION-NAME TEAM-NAME",
		Short: "Revoke permission from a team for a resource",
		Long:  `The revoke command allows an account owner to revoke permissions of a team from the specified resource.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return revokePermission(AMP, cmd, args)
		},
	}

	listPermissionsCmd = &cobra.Command{
		Use:   "list-permission ORGANIZATION-NAME TEAM-NAME",
		Short: "List permissions by team",
		Long:  `The list permission command displays the permissions, filtered by teams in an organization.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPermission(AMP, cmd, args)
		},
	}

	change bool
	reset  bool
)

func init() {
	RootCmd.AddCommand(LoginCmd)
	RootCmd.AddCommand(AccountCmd)
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(pwdResetCmd)
	AccountCmd.AddCommand(pwdChangeCmd)
	AccountCmd.AddCommand(switchRoleCmd)
	AccountCmd.AddCommand(listUsersCmd)
	AccountCmd.AddCommand(infoCmd)
	AccountCmd.AddCommand(editCmd)
	AccountCmd.AddCommand(deleteCmd)
	AccountCmd.AddCommand(createOrganizationCmd)
	AccountCmd.AddCommand(editOrganizationCmd)
	AccountCmd.AddCommand(deleteOrganizationCmd)
	AccountCmd.AddCommand(listOrganizationsCmd)
	AccountCmd.AddCommand(addOrganizationMembersCmd)
	AccountCmd.AddCommand(removeOrganizationMembersCmd)
	AccountCmd.AddCommand(transferOwnershipCmd)
	AccountCmd.AddCommand(createTeamCmd)
	AccountCmd.AddCommand(editTeamCmd)
	AccountCmd.AddCommand(deleteTeamCmd)
	AccountCmd.AddCommand(listTeamsCmd)
	AccountCmd.AddCommand(addTeamMembersCmd)
	AccountCmd.AddCommand(removeTeamMembersCmd)
	AccountCmd.AddCommand(grantPermissionCmd)
	AccountCmd.AddCommand(editPermissionCmd)
	AccountCmd.AddCommand(revokePermissionCmd)
	AccountCmd.AddCommand(listPermissionsCmd)
}

// login validates the input command line arguments and logs a user into their account
// by invoking the corresponding rpc/storage method
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
	fmt.Println("Welcome back, ", username, "!")
	color.Unset()
	return nil
}

// signup validates the input command line arguments and allows a user to create an account
// by invoking the corresponding rpc/storage method
func signUp(amp *client.AMP) (err error) {
	fmt.Println("This will sign you up for a new personal AMP account.")
	username := getUserName()
	email, er := getEmailAddress()
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.SignUpRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.SignUp(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Hi", username, "!, Please check your email to complete the signup process.")
	color.Unset()
	return nil
}

// verify validates the input command line arguments and verifies an account
// by invoking the corresponding rpc/storage method
func verify(amp *client.AMP) (err error) {
	fmt.Println("This will verify your account and confirm your password.")
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
	fmt.Println("Hi", username, "! Your account has now been activated.")
	color.Unset()
	return nil
}

// pwdReset validates the input command line arguments and resets the current password
// by invoking the corresponding rpc/storage method
func pwdReset(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will send a password reset link to your email address.")
	username := getUserName()
	email, er := getEmailAddress()
	if er != nil {
		return fmt.Errorf("user error: %v", er)
	}
	request := &account.PasswordResetRequest{
		Name:  username,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.PasswordReset(context.Background(), request)
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
func pwdChange(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will allow you to update your existing password.")
	username := getUserName()
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Enter your existing password.")
	color.Unset()
	existingPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Enter new password.")
	color.Unset()
	newPwd, err := getPwd()
	if err != nil {
		return fmt.Errorf("user error: %v", err)
	}
	getConfirmNewPwd(newPwd)
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

// switchRole validates the input command line arguments and switches to an organization
// to the input value by invoking the corresponding rpc/storage method
func switchRole(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will switch between the primary account and an organization account that a user is part of.")
	var org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.TeamRequest{
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.Switch(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Switch organization successful - ", org)
	color.Unset()
	return nil
}

// listUser validates the input command line arguments and lists the users available on the system
// by invoking the corresponding rpc/storage method
func listUser(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all users available on the system.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.AccountsRequest{
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListAccounts(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Accounts) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\tORGANIZATION\tTEAM\t")
	fmt.Fprintln(w, "----\t-----\t------------------\t------------\t------------\t----\t")
	for _, info := range response.Accounts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", info.Name, info.Email, info.EmailVerified, info.AccountType, org, team)
	}
	w.Flush()
	return nil
}

// getAccountDetails validates the input command line arguments and displays the account
// information by invoking the corresponding rpc/storage method
func getAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will display the information of a specific account.")
	var name string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify account name")
		color.Unset()
		name = getUserName()
	case 1:
		name = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.AccountRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.GetAccountDetails(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ACCOUNT NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\t")
	fmt.Fprintln(w, "------------\t-----\t------------------\t------------\t")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", response.Account.Name, response.Account.Email, response.Account.EmailVerified, response.Account.AccountType)
	w.Flush()
	return nil
}

// editAccount validates the input command line arguments and modifies the account
// information by invoking the corresponding rpc/storage method
func editAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will modify the information of a specific account.")
	var name string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify account name")
		color.Unset()
		name = getUserName()
	case 1:
		name = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.UpdateRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.UpdateAccount(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ACCOUNT NAME\tEMAIL\tIS EMAIL VERIFIED?\tACCOUNT TYPE\t")
	fmt.Fprintln(w, "------------\t-----\t------------------\t------------\t")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", response.Account.Name, response.Account.Email, response.Account.EmailVerified, response.Account.AccountType)
	w.Flush()
	return nil
}

// deleteAccountDetails validates the input command line arguments and deletes the account
// by invoking the corresponding rpc/storage method
func deleteAccount(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will delete the specified account.")
	var name string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify account name")
		color.Unset()
		name = getUserName()
	case 1:
		name = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.AccountRequest{
		Name: name,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, er := client.DeleteAccount(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Account '", name, "' deleted successfully.")
	color.Unset()
	return nil
}

// createOrganization validates the input command line arguments and creates an organization
// with name and email by invoking the corresponding rpc/storage method
func createOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will create an organization account with the specified account name and email address.")
	var org, email string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 2:
		org = args[0]
		email = args[1]
		mail, er := account.CheckEmailAddress(email)
		if er != nil {
			return fmt.Errorf("user error: %v", er)
		} else {
			email = mail
		}
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.OrganizationRequest{
		Name:  org,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.CreateOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' created successfully.")
	color.Unset()
	return nil
}

// editOrganization validates the input command line arguments and creates an organization
// with name and email by invoking the corresponding rpc/storage method
func editOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit an organization account.")
	var org, email string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify email")
		color.Unset()
		email, err = getEmailAddress()
	case 2:
		org = args[0]
		email = args[1]
		mail, er := account.CheckEmailAddress(email)
		if er != nil {
			return fmt.Errorf("user error: %v", er)
		} else {
			email = mail
		}
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.OrganizationRequest{
		Name:  org,
		Email: email,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' edited successfully.")
	color.Unset()
	return nil
}

// deleteOrganization validates the input command line arguments and deletes an organization
// by invoking the corresponding rpc/storage method
func deleteOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will delete an organization account.")
	var org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.OrganizationRequest{
		Name: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Organization '", org, "' deleted successfully.")
	color.Unset()
	return nil
}

// listOrganization validates the input command line arguments and lists the available
// organizations by invoking the corresponding rpc/storage method
func listOrganization(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all the organizations available.")
	if len(args) > 0 {
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.AccountsRequest{}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListAccounts(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Accounts) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tEMAIL\t")
	fmt.Fprintln(w, "----\t----\t")
	for _, info := range response.Accounts {
		fmt.Fprintf(w, "%s\t%s\t\n", info.Name, info.Email)
	}
	w.Flush()
	return nil
}

// addOrgMember validates the input command line arguments and adds new members to
// an organization by invoking the corresponding rpc/storage method
func addOrgMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will add new members to the specified organization.")
	var org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	members := getMembers()
	request := &account.OrganizationMembershipsRequest{
		Name:    org,
		Members: members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, er := client.AddOrganizationMemberships(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully added to organization '", org, "'.")
	color.Unset()
	return nil
}

// removeOrgMember validates the input command line arguments and removes members from
// an organization by invoking the corresponding rpc/storage method
func removeOrgMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will remove members from the specified organization.")
	var org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	members := getMembers()
	request := &account.OrganizationMembershipsRequest{
		Name:    org,
		Members: members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteOrganizationMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully removed from organization '", org, "'.")
	color.Unset()
	return nil
}

// transferOwnership validates the input command line arguments and transfer ownership of a resource
// by invoking the corresponding rpc/storage method
func transferOwnership(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will transfer ownership of a resource to the specified organization.")
	var resId, org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 2:
		resId = args[0]
		org = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.TransferOwnership(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Ownership Transfer successful for organization '", org, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// createTeam validates the input command line arguments and creates a team in an organization
// by invoking the corresponding rpc/storage method
func createTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will create a team within the specified organization.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.CreateTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' created successfully.")
	color.Unset()
	return nil
}

// editTeam validates the input command line arguments and creates a team in an organization
// by invoking the corresponding rpc/storage method
func editTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit team information.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' updated successfully.")
	color.Unset()
	return nil
}

// deleteTeam validates the input command line arguments and deletes a team in an organization
// by invoking the corresponding rpc/storage method
func deleteTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will deletes a team within the specified organization.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.TeamRequest{
		Organization: org,
		Name:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteTeam(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Team '", team, "' deleted successfully.")
	color.Unset()
	return nil
}

// listTeam validates the input command line arguments and lists the available teams
// within an organization by invoking the corresponding rpc/storage method
func listTeam(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list all the teams available in a specific organization.")
	var org string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
	case 1:
		org = args[0]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.TeamRequest{
		Organization: org,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListTeams(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Teams) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION\tMEMBERS\t")
	fmt.Fprintln(w, "----\t-----------\t-------\t")
	for _, info := range response.Teams {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", info.Name, info.Description, info.Members)
	}
	w.Flush()
	return nil
}

// addTeamMember validates the input command line arguments and adds new members to
// a team in an organization by invoking the corresponding rpc/storage method
func addTeamMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will add new members to the specified team within an organization.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	members := getMembers()
	request := &account.TeamMembershipsRequest{
		Organization: org,
		Name:         team,
		Members:      members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.AddTeamMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully added to team '", team, "'.")
	color.Unset()
	return nil
}

// removeTeamMember validates the input command line arguments and removes members from
// a team in an organization by invoking the corresponding rpc/storage method
func removeTeamMember(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will remove members from the specified team within an organization.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	members := getMembers()
	request := &account.TeamMembershipsRequest{
		Organization: org,
		Name:         team,
		Members:      members,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.DeleteTeamMemberships(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Members successfully removed from team '", team, "'.")
	color.Unset()
	return nil
}

// grantPermission validates the input command line arguments and grants permissions
// to a team in an organization by invoking the corresponding rpc/storage method
func grantPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will grant permissions to the specified resource.")
	var resId, level, org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		resId = args[0]
		level = args[1]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 3:
		resId = args[0]
		level = args[1]
		org = args[2]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 4:
		resId = args[0]
		level = args[1]
		org = args[2]
		team = args[3]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Level:        level,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.GrantPermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Grant Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// editPermission validates the input command line arguments and edits permissions
// of a team in an organization by invoking the corresponding rpc/storage method
func editPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will edit permissions of the specified resource.")
	var resId, level, org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify permission level")
		color.Unset()
		level = getPermissionLevel()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		resId = args[0]
		level = args[1]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 3:
		resId = args[0]
		level = args[1]
		org = args[2]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 4:
		resId = args[0]
		level = args[1]
		org = args[2]
		team = args[3]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Level:        level,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.EditPermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Edit Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// revokePermission validates the input command line arguments and revokes permissions
// of a team in an organization by invoking the corresponding rpc/storage method
func revokePermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will revoke permissions from the specified resource.")
	var resId, org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify resource id")
		color.Unset()
		resId = getResourceID()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		resId = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		resId = args[0]
		org = args[1]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 3:
		resId = args[0]
		org = args[1]
		team = args[2]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		ResourceId:   resId,
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	_, err = client.RevokePermission(context.Background(), request)
	if err != nil {
		return fmt.Errorf("server error: %v", err)
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Revoke Permission successful for team '", team, "' for resource '", resId, "'.")
	color.Unset()
	return nil
}

// listPermission validates the input command line arguments and lists the permissions
// by invoking the corresponding rpc/storage method
func listPermission(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
	fmt.Println("This will list the permissions of a team within an organization.")
	var org, team string
	switch len(args) {
	case 0:
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify organization")
		color.Unset()
		org = getOrganization()
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 1:
		org = args[0]
		color.Set(color.FgRed, color.Bold)
		fmt.Println("must specify team")
		color.Unset()
		team = getTeam()
	case 2:
		org = args[0]
		team = args[1]
	default:
		defer color.Set(color.FgRed, color.Bold)
		return errors.New("too many arguments - check again")
	}
	request := &account.PermissionRequest{
		Organization: org,
		Team:         team,
	}
	client := account.NewAccountServiceClient(amp.Conn)
	response, er := client.ListPermissions(context.Background(), request)
	if er != nil {
		return fmt.Errorf("server error: %v", er)
	}
	if response == nil || len(response.Permissions) == 0 {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("No information available!")
		color.Unset()
		return nil
	}
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "RESOURCE ID\tLEVEL\tTEAM\tORGANIZATION\t")
	fmt.Fprintln(w, "-----------\t-----\t----\t------------\t")
	for _, info := range response.Permissions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", info.ResourceId, info.Level, info.Team, info.Organization)
	}
	w.Flush()
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
		fmt.Println("Username is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getUserName()
	}
	return
}

func getCode() (code string) {
	fmt.Print("code: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&code)
	color.Unset()
	err := account.CheckVerificationCodeFormat(code)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Code is invalid. Code must be 8 characters long. Try again!")
		color.Unset()
		fmt.Println("")
		return getCode()
	}
	return
}

func getEmailAddress() (email string, err error) {
	fmt.Print("email: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scanln(&email)
	color.Unset()
	email, err = account.CheckEmailAddress(email)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Email in incorrect format. Try again!")
		color.Unset()
		fmt.Println("")
		return getEmailAddress()
	}
	return
}

func getPwd() (password string, err error) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Password is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getPwd()
	}
	password = string(pw)
	err = account.CheckPassword(password)
	if err != nil {
		fmt.Println(err)
	}

	err = account.CheckPasswordStrength(password)
	if err != nil {
		if strings.Contains(err.Error(), "password too weak") {
			color.Set(color.FgRed, color.Bold)
			fmt.Println("Password entered is too weak. Password must be at least 8 characters long. Try again!")
			color.Unset()
			fmt.Println("")
			return getPwd()
		} else {
			return
		}
	}
	return
}

func getConfirmNewPwd(newPwd string) {
	color.Set(color.FgGreen, color.Bold)
	fmt.Println("Enter Password again for confirmation.")
	color.Unset()
	confirmNewPwd, _ := getPwd()
	if confirmNewPwd != newPwd {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Password mismatch. Try again!")
		color.Unset()
		fmt.Println("")
		getConfirmNewPwd(newPwd)
	} else {
		return
	}
	return
}

func getResourceID() (id string) {
	fmt.Print("Resource ID: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&id)
	color.Unset()
	err := account.CheckResourceID(id)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Resource ID is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getResourceID()
	}
	return
}

func getPermissionLevel() (level string) {
	fmt.Print("Permission Level: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&level)
	color.Unset()
	err := account.CheckPermissionLevel(level)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Permission Level is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getPermissionLevel()
	}
	return
}

func getOrganization() (org string) {
	fmt.Print("Organization: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&org)
	color.Unset()
	err := account.CheckOrganizationName(org)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Organization Name is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getOrganization()
	}
	return
}

func getTeam() (team string) {
	fmt.Print("Team Name: ")
	color.Set(color.FgGreen, color.Bold)
	fmt.Scan(&team)
	color.Unset()
	err := account.CheckTeamName(team)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("Team name is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getTeam()
	}
	return
}

func getMembers() (memArr []string) {
	fmt.Print("Member name(s): ")
	color.Set(color.FgGreen, color.Bold)
	reader := bufio.NewReader(os.Stdin)
	members, _ := reader.ReadString('\n')
	memArr = strings.Fields(members)
	color.Unset()
	err := account.CheckMembers(memArr)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("At least one member is mandatory. Try again!")
		color.Unset()
		fmt.Println("")
		return getMembers()
	}
	return
}
