package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/authn"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Cobra definitions for account management related commands
var (
	signUpCmd = &cobra.Command{
		Use:     "signup",
		Short:   "Signup for a new account",
		Example: "amp account signup --name=jdoe --email=jdoe@fakemail.me --password=p@s5wrd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return signUp(AMP, cmd)
		},
	}

	verifyCmd = &cobra.Command{
		Use:     "verify",
		Short:   "Verify account",
		Example: "amp account verify --token=this-is-a-very-very-very-long-token-code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify(AMP, cmd)
		},
	}

	loginCmd = &cobra.Command{
		Use:     "login",
		Short:   "Login to account",
		Example: "amp account login --name=jdoe --password=p@s5wrd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(AMP, cmd)
		},
	}

	forgotLoginCmd = &cobra.Command{
		Use:     "forgot-login",
		Short:   "Retrieve account name",
		Example: "amp account forgot-login --email=jdoe@fakemail.me",
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
		Example: "amp account password change --name=jdoe --password=p@s5wrd --new-password=v@larm0rghuli$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdChange(AMP, cmd)
		},
	}

	pwdResetCmd = &cobra.Command{
		Use:     "reset",
		Short:   "Reset account password",
		Example: "amp account password reset --name=jdoe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdReset(AMP, cmd)
		},
	}

	pwdSetCmd = &cobra.Command{
		Use:     "set",
		Short:   "Set account password",
		Example: "amp account password set --token=this-is-a-token-sample --password=v@lard0haeri$",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pwdSet(AMP, cmd)
		},
	}

	switchCmd = &cobra.Command{
		Use:     "switch",
		Short:   "Switch account",
		Example: "amp account switch --name=swag",
		RunE: func(cmd *cobra.Command, args []string) error {
			return switchAccount(AMP, cmd)
		},
	}

	whoAmICmd = &cobra.Command{
		Use:     "whoami",
		Short:   "Display currently logged-in user",
		Example: "amp account whoami",
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoAmI()
		},
	}

	logoutCmd = &cobra.Command{
		Use:     "logout",
		Short:   "Logout current user",
		Example: "amp account logout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return logout()
		},
	}

	username string
	email    string
	password string
	token    string
	newPwd   string

	//TODO: pass verbose as arg
	manager = newCmdManager("")
)

//Adding account management commands to the account command
func init() {
	AccountCmd.AddCommand(signUpCmd)
	AccountCmd.AddCommand(verifyCmd)
	AccountCmd.AddCommand(loginCmd)
	AccountCmd.AddCommand(forgotLoginCmd)
	AccountCmd.AddCommand(pwdCmd)
	pwdCmd.AddCommand(pwdChangeCmd)
	pwdCmd.AddCommand(pwdResetCmd)
	pwdCmd.AddCommand(pwdSetCmd)
	AccountCmd.AddCommand(switchCmd)
	AccountCmd.AddCommand(whoAmICmd)
	AccountCmd.AddCommand(logoutCmd)

	signUpCmd.Flags().StringVar(&username, "name", username, "Account Name")
	signUpCmd.Flags().StringVar(&email, "email", email, "Email ID")
	signUpCmd.Flags().StringVar(&password, "password", password, "Password")

	verifyCmd.Flags().StringVar(&token, "token", token, "Verification Token")

	loginCmd.Flags().StringVar(&username, "name", username, "Account Name")
	loginCmd.Flags().StringVar(&password, "password", password, "Password")

	forgotLoginCmd.Flags().StringVar(&email, "email", email, "Email ID")

	pwdSetCmd.Flags().StringVar(&token, "token", token, "Verification Token")
	pwdSetCmd.Flags().StringVar(&password, "password", password, "Password")

	pwdChangeCmd.Flags().StringVar(&username, "name", username, "Account Name")
	pwdChangeCmd.Flags().StringVar(&password, "password", password, "Current Password")
	pwdChangeCmd.Flags().StringVar(&newPwd, "new-password", newPwd, "New Password")

	pwdResetCmd.Flags().StringVar(&username, "name", username, "Account Name")

	switchCmd.Flags().StringVar(&username, "name", username, "Account Name")

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

// login validates the input command line arguments and allows login to an existing account
// by invoking the corresponding rpc/storage method
func login(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("username: ")
		username = getName()
	}
	if cmd.Flag("password").Changed {
		password = cmd.Flag("password").Value.String()
	} else {
		password = getPassword()
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

// switchAccount validates the input command line arguments and switches from personal account to an organization account
// by invoking the corresponding rpc/storage method
func switchAccount(amp *cli.AMP, cmd *cobra.Command) (err error) {
	if cmd.Flag("name").Changed {
		username = cmd.Flag("name").Value.String()
	} else {
		fmt.Print("account: ")
		username = getName()
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

// whoAmI validates the input command line arguments and displays the current account
// by invoking the corresponding rpc/storage method
func whoAmI() (err error) {
	token, err := cli.ReadToken()
	if err != nil {
		manager.fatalf("You are not logged in.")
		return
	}
	pToken, _ := jwt.ParseWithClaims(token, &authn.AccountClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
	if claims, ok := pToken.Claims.(*authn.AccountClaims); ok {
		if claims.ActiveOrganization != "" {
			manager.printf(colSuccess, "Logged in as organization %s (on behalf of user %s).", claims.ActiveOrganization, claims.AccountName)
		} else {
			manager.printf(colSuccess, "Logged in as user %s.", claims.AccountName)
		}
	}
	return nil
}

// logout validates the input command line arguments and logs out of the current account
// by invoking the corresponding rpc/storage method
func logout() (err error) {
	err = cli.RemoveToken()
	if err != nil {
		manager.fatalf(grpc.ErrorDesc(err))
		return
	}
	manager.printf(colSuccess, "You have been successfully logged out!")
	return nil
}
