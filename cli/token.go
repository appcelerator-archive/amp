package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/appcelerator/amp/api/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	ampConfigFolder = filepath.Join(".config", "amp")
)

const (
	//Credentials suffix for authentication token file
	Credentials = ".credentials"
)

// SaveToken saves the authentication token to file
func SaveToken(headers metadata.MD, server string) error {
	tokens := headers[auth.TokenKey]
	ampTokenFile := strings.TrimSuffix(server, DefaultPort) + Credentials
	// Extract token from header
	if len(tokens) == 0 {
		return fmt.Errorf("invalid token")
	}
	token := tokens[0]
	if token == "" {
		return fmt.Errorf("invalid token")
	}
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %s", err)
	}
	if err := os.MkdirAll(filepath.Join(usr.HomeDir, ampConfigFolder), os.ModePerm); err != nil {
		return fmt.Errorf("cannot create folder: %s", err)
	}
	if err := ioutil.WriteFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile), []byte(token), 0600); err != nil {
		return fmt.Errorf("cannot write token: %s", err)
	}
	return nil
}

// ReadToken reads the authentication token from file
func ReadToken(server string) (string, error) {
	ampTokenFile := strings.TrimSuffix(server, DefaultPort) + Credentials
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("cannot get current user: %s", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile))
	if err != nil {
		return "", fmt.Errorf("cannot read token: %s", err)
	}
	return string(data), nil
}

// RemoveToken deletes the authentication token file
func RemoveToken(server string) error {
	ampTokenFile := strings.TrimSuffix(server, DefaultPort) + Credentials
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot get current user: %s", err)
	}
	filePath := filepath.Join(usr.HomeDir, ampConfigFolder, ampTokenFile)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("you are not logged in. Use `amp login` or `amp user signup`")
	}
	err = os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("cannot remove token: %s", err)
	}
	return nil
}

// LoginCredentials represents login credentials
type LoginCredentials struct {
	Token string
}

// GetRequestMetadata implements credentials.PerRPCCredentials
func (c *LoginCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		auth.AuthorizationHeader: auth.ForgeAuthorizationHeader(c.Token),
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (c *LoginCredentials) RequireTransportSecurity() bool {
	return true
}

// GetToken returns the stored token
func GetToken(server string) string {
	token, err := ReadToken(server)
	if err != nil {
		return ""
	}
	return token
}
