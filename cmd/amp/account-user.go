package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
)

// UserCmd is the main command for attaching user sub-commands.
var (
	signUpCmd = &cobra.Command{
		Use:     "signup",
		Short:   "Signup for a new account",
		Example: "amp user signup --name=jdoe --email=jdoe@fakemail.me --password=p@s5wrd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP, cmd)
		},
	}

	verifyCmd = &cobra.Command{
		Use:     "verify",
		Short:   "Verify account",
		Example: "amp user verify --token=this-is-a-very-very-very-long-token-code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP, cmd)
		},
	}

	forgotLoginCmd = &cobra.Command{
		Use:     "forgot-login",
		Short:   "Retrieve account name",
		Example: "amp user forgot-login --email=jdoe@fakemail.me",
		RunE: func(cmd *cobra.Command, args []string) error {
			return forgotLogin(AMP, cmd)
		},
	}

	pwdCmd = &cobra.Command{
		Use:   "password",
		Short: "Account password operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwd()
		},
	}

	pwdChangeCmd = &cobra.Command{
		Use:     "change",
		Short:   "Change account password",
		Example: "amp user password change --name=jdoe --password=p@s5wrd --new-password=v@larm0rghuli$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd)
		},
	}

	pwdResetCmd = &cobra.Command{
		Use:     "reset",
		Short:   "Reset account password",
		Example: "amp user password reset --name=jdoe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd)
		},
	}

	pwdSetCmd = &cobra.Command{
		Use:     "set",
		Short:   "Set account password",
		Example: "amp user password set --token=this-is-a-token-sample --password=v@lard0haeri$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdSet(AMP, cmd)
		},
	}

	listUserCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List user",
		Example: "amp user ls -q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(AMP, cmd)
		},
	}

	deleteUserCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove user",
		Example: "amp user rm --name=hpotter \namp user del --name=hpotter",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteUser(AMP, cmd)
		},
	}

	getUserCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get user info",
		Example: "amp user get --name=rweasley",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getUser(AMP, cmd)
		},
	}

	email  string
	token  string
	newPwd string
)

func init() {
	UserCmd.AddCommand(signUpCmd)
	UserCmd.AddCommand(verifyCmd)
	UserCmd.AddCommand(forgotLoginCmd)
	UserCmd.AddCommand(pwdCmd)
	pwdCmd.AddCommand(pwdChangeCmd)
	pwdCmd.AddCommand(pwdResetCmd)
	pwdCmd.AddCommand(pwdSetCmd)
	UserCmd.AddCommand(listUserCmd)
	UserCmd.AddCommand(deleteUserCmd)
	UserCmd.AddCommand(getUserCmd)

	signUpCmd.Flags().StringVar(&username, "name", username, "Account Name")
	signUpCmd.Flags().StringVar(&email, "email", email, "Email ID")
	signUpCmd.Flags().StringVar(&password, "password", password, "Password")

	verifyCmd.Flags().StringVar(&token, "token", token, "Verification Token")

	forgotLoginCmd.Flags().StringVar(&email, "email", email, "Email ID")

	pwdSetCmd.Flags().StringVar(&token, "token", token, "Verification Token")
	pwdSetCmd.Flags().StringVar(&password, "password", password, "Password")

	pwdChangeCmd.Flags().StringVar(&username, "name", username, "Account Name")
	pwdChangeCmd.Flags().StringVar(&password, "password", password, "Current Password")
	pwdChangeCmd.Flags().StringVar(&newPwd, "new-password", newPwd, "New Password")

	pwdResetCmd.Flags().StringVar(&username, "name", username, "Account Name")

	listUserCmd.Flags().BoolP("quiet", "q", false, "Only display User Name")

	getUserCmd.Flags().StringVar(&name, "name", name, "Account Name")

	deleteUserCmd.Flags().StringVar(&name, "name", name, "Account Name")
}

// signUp validates the input command line arguments and creates a new account
// by invoking the corresponding rpc/storage method
func signUp(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = getName()
	}
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = getEmailAddress()
	}
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = getPassword()
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
		token = getToken()
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

// forgotLogin validates the input command line arguments and retrieves account name
// by invoking the corresponding rpc/storage method
func forgotLogin(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = getEmailAddress()
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

// pwd validates the input command line arguments and performs password-related operations
// by invoking the corresponding rpc/storage method
func pwd() (err error) {
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

// listUser validates the input command line arguments and lists all users
// by invoking the corresponding rpc/storage method
func listUser(amp *cli.AMP, cmd *cobra.Command) (err error) {
	request := &account.ListUsersRequest{}
	accClient := account.NewAccountClient(amp.Conn)
	reply, er := accClient.ListUsers(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		return fmt.Errorf("unable to convert quiet parameter : %v", err.Error())
	} else if quiet {
		for _, user := range reply.Users {
			fmt.Println(user.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\t")
	for _, user := range reply.Users {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Email)
	}
	w.Flush()
	return nil
}

// deleteUser validates the input command line arguments and deletes a user
// by invoking the corresponding rpc/storage method
func deleteUser(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		name = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		name = getName()
	}

	request := &account.DeleteUserRequest{
		Name: name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err = accClient.DeleteUser(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "Successfully deleted user.")
	return nil
}

// getUser validates the input command line arguments and retrieves info of a user
// by invoking the corresponding rpc/storage method
func getUser(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		name = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		name = getName()
	}

	request := &account.GetUserRequest{
		Name: name,
	}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.GetUser(context.Background(), request)
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\tVERIFIED\tCREATED\t")
	fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", reply.User.Name, reply.User.Email, reply.User.IsVerified, convertTime(reply.User.CreateDt))
	w.Flush()
	return nil
}
