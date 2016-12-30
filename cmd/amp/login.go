package main

import (
	"fmt"
	"regexp"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login via GitHub",
	Long:  `Create a GitHub access token and store it in your Config file to authenticate further commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := AMP.Connect()
		if err != nil {
			return err
		}
		return Login(AMP)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}

// Login creates a github access token and store it in your Config file to authenticate further commands
func Login(a *client.AMP) error {
	username, password, lastEight, name, err := basicLogin(a)
	if err != nil {
		otpError := regexp.MustCompile("Must specify two-factor authentication OTP code")
		if otpError.MatchString(err.Error()) {
			lastEight, name, err = otpLogin(a, username, password)
		}
		if err != nil {
			cli.PrintErr(err)
		}
	}

	a.Configuration.GitHub = lastEight
	cli.SaveConfiguration(a.Configuration)

	welcomeUser(name)
	return nil
}

// GetUsername prompts for username and returns the result
func getUsername() (username string) {
	color.Set(color.FgMagenta, color.Bold)
	fmt.Println("Login using GitHub")
	fmt.Print("GitHub username or email: ")
	color.Unset()
	fmt.Scanln(&username)
	return
}

// GetPassword prompts for password and returns the result
func getPassword() (password string, err error) {
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset()
	fmt.Print("GitHub password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		return
	}
	password = string(pw)
	fmt.Println("")
	return
}

// GetOTP prompts for two factor auth and returns the result
func getOTP() (otp string, err error) {
	color.Set(color.FgBlue)
	defer color.Unset()
	fmt.Println("Two Factor Authentication required")
	fmt.Print("Authentication code: ")
	otpraw, err := gopass.GetPasswd()
	if err != nil {
		return
	}
	fmt.Println("")
	otp = string(otpraw)
	return
}

// WelcomeUser prints a welcome message
func welcomeUser(name string) {
	color.Set(color.FgCyan)
	defer color.Unset()
	fmt.Println("Welcome", name)
	fmt.Println("You are logged in")
}

func basicLogin(a *client.AMP) (username, password, lastEight, name string, err error) {
	username = getUsername()
	password, err = getPassword()
	if err != nil {
		return
	}
	lastEight, name, err = a.GithubOauth(username, password, "")
	if err != nil {
		badCredError := regexp.MustCompile("Bad credentials")
		if badCredError.MatchString(err.Error()) {
			return basicLogin(a)
		}
	}
	return
}

func otpLogin(a *client.AMP, username, password string) (lastEight, name string, err error) {
	otp, err := getOTP()
	if err != nil {
		return
	}
	lastEight, name, err = a.GithubOauth(username, password, otp)
	if err != nil {
		otpError := regexp.MustCompile("Must specify two-factor authentication OTP code")
		if otpError.MatchString(err.Error()) {
			return otpLogin(a, username, password)
		}
	}
	return
}
