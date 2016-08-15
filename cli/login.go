package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"path"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/howeyc/gopass"
	"golang.org/x/oauth2"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

func getUsername() (username string) {
	color.Set(color.FgMagenta, color.Bold)
	fmt.Println("Login using github")
	fmt.Print("github username or email: ")
	color.Unset()
	fmt.Scanln(&username)
	return
}

func getPassword() (password string, err error) {
	color.Set(color.FgMagenta, color.Bold)
	defer color.Unset()
	fmt.Print("github password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		return
	}
	password = string(pw)
	fmt.Println("")
	return
}

func getAuth(client *github.Client, username, password string) (auth *github.Authorization, err error) {
	note := "note"
	auth, _, err = client.Authorizations.Create(&github.AuthorizationRequest{
		Scopes: []github.Scope{
			github.ScopeRepo,
			github.ScopeAdminRepoHook,
		},
		Note: &note,
	})
	if err != nil {
		otpError := regexp.MustCompile("Must specify two-factor authentication OTP code")
		if otpError.MatchString(err.Error()) {
			color.Set(color.FgBlue)
			defer color.Unset()
			fmt.Println("Two Factor Authentication required")
			fmt.Print("authentication code: ")
			authenticationCode, pwerr := gopass.GetPasswd()
			if pwerr != nil {
				return nil, pwerr
			}
			fmt.Println("")
			basicAuth := github.BasicAuthTransport{
				Username: username,
				Password: string(password),
				OTP:      string(authenticationCode),
			}
			auth, err = getAuth(github.NewClient(basicAuth.Client()), username, password)
		}
	}
	return
}

func getToken() (token string, err error) {
	username := getUsername()
	password, err := getPassword()
	if err != nil {
		return
	}
	basicAuth := github.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	auth, err := getAuth(github.NewClient(basicAuth.Client()), username, password)
	if err != nil {
		return
	}
	token = *auth.Token
	return
}

func getOauthClient(token string) (client *github.Client) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client = github.NewClient(tc)
	return
}

// Login creates a github access token and store it in your config file to authenticate further commands
func (a *AMP) Login() {
	fmt.Println(a.Config)
	homedir, err := homedir.Dir()
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
  fmt.Println(homedir)
	contents, err := ioutil.ReadFile(path.Join(homedir, ".ampswarm.yaml"))
	if (err != nil) {
		notFoundError := regexp.MustCompile("no such file or directory")
		if (!notFoundError.MatchString(err.Error())) {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println(contents)
	
	m := make(map[interface{}]interface{})
	
	err = yaml.Unmarshal(contents, &m)
	if err != nil {
    fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println(m)
	
	m["github"] = "foo"
	
	fmt.Println(m)
	
	contents, err = yaml.Marshal(&m)
	if err != nil {
    fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(contents))
	
	err = ioutil.WriteFile(path.Join(homedir, ".ampswarm.yaml"), contents, os.ModeAppend)
	if err != nil {
    fmt.Println(err)
		os.Exit(1)
	}
  contents, err = ioutil.ReadFile(path.Join(homedir, ".ampswarm.yaml"))
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	m = make(map[interface{}]interface{})
	
	err = yaml.Unmarshal(contents, &m)
	if err != nil {
    fmt.Println(err)
		os.Exit(1)
	}
	
	fmt.Println(m["github"])

	// color.Set(color.FgRed)
	// defer color.Unset()
	// token, err := getToken()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// client := getOauthClient(token)

	// user, _, err := client.Users.Get("")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// member, _, err := client.Organizations.IsMember("appcelerator", *user.Login)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// if !member {
	// 	fmt.Println("not a member of the organization")
	// 	os.Exit(1)
	// }
	// color.Set(color.FgCyan)
	// fmt.Println("Ok, now save the token")
}
