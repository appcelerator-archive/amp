package oauth

import (
	"fmt"
	"regexp"

	"github.com/appcelerator/amp/data/storage"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const (
	organization = "appcelerator"
)

// Oauth handles third party authentication
type Oauth struct {
	Store        storage.Interface
	ClientID     string
	ClientSecret string
}

// Create creates new oauth credencials and stores them
func (o *Oauth) Create(ctx context.Context, in *AuthRequest) (out *AuthReply, err error) {
	// try to create an authorization
	auth, err := o.getAuth(in.Username, in.Password, in.Otp)
	if err != nil {
		// handle user errors
		badCredentials := regexp.MustCompile("401 Bad credentials")
		if badCredentials.MatchString(err.Error()) {
			return nil, grpc.Errorf(codes.NotFound, "Bad credentials")
		}
		otpMissing := regexp.MustCompile("Must specify two-factor authentication OTP code")
		if otpMissing.MatchString(err.Error()) {
			return nil, grpc.Errorf(codes.Unauthenticated, "Must specify two-factor authentication OTP code")
		}
		return
	}

	// handle variations in the flow
	new := false
	token := *auth.Token

	// handle get case

	// try to retrieve token from store
	if token == "" {
		token, err = getToken(ctx, o.Store, *auth.TokenLastEight)
		if err != nil {
			return
		}
	}

	// delete the authorization and create a new one
	if token == "" {
		err = deleteAuth(auth, in.Username, in.Password, in.Otp)
		if err != nil {
			return
		}
		auth, err = o.getAuth(in.Username, in.Password, in.Otp)
		if err != nil {
			return
		}
		new = true
		token = *auth.Token
	}
	if token == "" {
		return nil, grpc.Errorf(codes.Unknown, "Could not retrieve token")
	}

	// confirm organization
	name, err := checkOrganization(token)
	if err != nil {
		return
	}

	// save token to store if not in place
	if new {
		saved := Token{}
		err = o.Store.Create(ctx, *auth.TokenLastEight, &Token{
			Id:             int32(*auth.ID),
			Token:          token,
			TokenLastEight: *auth.TokenLastEight,
		}, &saved, 0)
		if err != nil {
			return
		}
	}

	// return sessionkey and name
	out = &AuthReply{
		SessionKey: *auth.TokenLastEight,
		Name:       name,
	}
	return
}

func deleteAuth(auth *github.Authorization, username, password, otp string) (err error) {
	client := getBasicClient(username, password, otp)
	_, err = client.Authorizations.Delete(*auth.ID)
	return
}

func getToken(ctx context.Context, store storage.Interface, key string) (token string, err error) {
	out := Token{}
	err = store.Get(ctx, key, &out, true)
	if err != nil {
		return
	}
	token = out.Token
	return
}

func getBasicClient(username, password, otp string) (client *github.Client) {
	basicAuth := github.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	if otp != "" {
		basicAuth.OTP = otp
	}
	client = github.NewClient(basicAuth.Client())
	return
}

func (o *Oauth) getAuth(username, password, otp string) (auth *github.Authorization, err error) {
	client := getBasicClient(username, password, otp)
	fingerPrint := "amplifier"
	auth, _, err = client.Authorizations.GetOrCreateForApp(o.ClientID, &github.AuthorizationRequest{
		Scopes: []github.Scope{
			github.ScopeRepo,
			github.ScopeAdminRepoHook,
		},
		ClientSecret: &o.ClientSecret,
		Fingerprint:  &fingerPrint,
	})
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

func checkOrganization(token string) (name string, err error) {
	client := getOauthClient(token)

	user, _, err := client.Users.Get("")
	if err != nil {
		return
	}
	member, _, err := client.Organizations.IsMember(organization, *user.Login)
	if err != nil {
		return
	}
	if !member {
		return "", grpc.Errorf(codes.PermissionDenied, "Not a member of the organization")
	}
	name = *user.Name
	return
}

// CheckAuthorization verifies that the user is logged in and has organization access
func CheckAuthorization(ctx context.Context, store storage.Interface) (name string, err error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "Metadata not ok")
	}
	fmt.Println(md)
	keys := md["sessionkey"]
	if len(keys) == 0 {
		return "", grpc.Errorf(codes.Unauthenticated, "SessionKey missing")
	}
	key := keys[0]
	if key == "" {
		return "", grpc.Errorf(codes.Unauthenticated, "SessionKey missing")
	}
	token, err := getToken(ctx, store, key)
	if err != nil {
		return
	}
	return checkOrganization(token)
}
