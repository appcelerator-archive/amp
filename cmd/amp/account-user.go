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

	email string
	token string
)

func init() {
	UserCmd.AddCommand(signUpCmd)
	UserCmd.AddCommand(verifyCmd)
	UserCmd.AddCommand(forgotLoginCmd)
	UserCmd.AddCommand(listUserCmd)
	UserCmd.AddCommand(deleteUserCmd)
	UserCmd.AddCommand(getUserCmd)

	signUpCmd.Flags().StringVar(&username, "name", username, "Account Name")
	signUpCmd.Flags().StringVar(&email, "email", email, "Email ID")
	signUpCmd.Flags().StringVar(&password, "password", password, "Password")

	verifyCmd.Flags().StringVar(&token, "token", token, "Verification Token")

	forgotLoginCmd.Flags().StringVar(&email, "email", email, "Email ID")

	listUserCmd.Flags().BoolP("quiet", "q", false, "Only display User Name")

	getUserCmd.Flags().StringVar(&name, "name", name, "Account Name")

	deleteUserCmd.Flags().StringVar(&name, "name", name, "Account Name")
}

// signUp validates the input command line arguments and creates a new account
// by invoking the corresponding rpc/storage method
func signUp(amp *cli.AMP, cmd *cobra.Command) error {
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
	_, err := accClient.SignUp(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Hi %s! Please check your email to complete the signup process.", username)
	return nil
}

// verify validates the input command line arguments and verifies an account
// by invoking the corresponding rpc/storage method
func verify(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flag("token").Changed {
		token = cmd.Flag("token").Value.String()
	} else {
		token = getToken()
	}

	request := &account.VerificationRequest{
		Token: token,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.Verify(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Your account has now been activated.")
	return nil
}

// forgotLogin validates the input command line arguments and retrieves account name
// by invoking the corresponding rpc/storage method
func forgotLogin(amp *cli.AMP, cmd *cobra.Command) error {
	if cmd.Flag("email").Changed {
		email = cmd.Flag("email").Value.String()
	} else {
		email = getEmailAddress()
	}

	request := &account.ForgotLoginRequest{
		Email: email,
	}
	accClient := account.NewAccountClient(amp.Conn)
	_, err := accClient.ForgotLogin(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Your login name has been sent to the address: %s", email)
	return nil
}

// listUser validates the input command line arguments and lists all users
// by invoking the corresponding rpc/storage method
func listUser(amp *cli.AMP, cmd *cobra.Command) error {
	request := &account.ListUsersRequest{}
	accClient := account.NewAccountClient(amp.Conn)
	reply, err := accClient.ListUsers(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		mgr.Error("unable to convert quiet parameter : %v", err.Error())
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
func deleteUser(amp *cli.AMP, cmd *cobra.Command) error {
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
	_, err := accClient.DeleteUser(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	mgr.Success("Successfully deleted user.")
	return nil
}

// getUser validates the input command line arguments and retrieves info of a user
// by invoking the corresponding rpc/storage method
func getUser(amp *cli.AMP, cmd *cobra.Command) error {
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
		mgr.Error(grpc.ErrorDesc(err))
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\tVERIFIED\tCREATED\t")
	fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", reply.User.Name, reply.User.Email, reply.User.IsVerified, convertTime(reply.User.CreateDt))
	w.Flush()
	return nil
}
