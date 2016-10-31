package client

import (
	"github.com/appcelerator/amp/api/rpc/oauth"
	"golang.org/x/net/context"
)

// GithubOauth does an rpc call to login via github
func (a *AMP) GithubOauth(username, password, otp string) (lastEight, name string, err error) {
	defer a.Disconnect()
	_, err = a.Connect()
	if err != nil {
		return
	}
	c := oauth.NewGithubClient(a.Conn)
	authReply, err := c.Create(context.Background(), &oauth.AuthRequest{
		Username: username,
		Password: password,
		Otp:      otp,
	})
	if err != nil {
		return
	}
	lastEight = authReply.SessionKey
	name = authReply.Name
	return
}
