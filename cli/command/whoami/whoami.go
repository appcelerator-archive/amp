package whoami

import (
	"errors"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/cli"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

// NewWhoAmICommand returns a new instance of the whoami command.
func NewWhoAmICommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "whoami",
		Short:   "Display currently logged-in user",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoami(c)
		},
	}
}

func whoami(c cli.Interface) error {
	token, err := cli.ReadToken()
	if err != nil {
		return errors.New("you are not logged in. Use `amp login` or `amp user signup`.")
	}
	pToken, _ := jwt.ParseWithClaims(token, &auth.LoginClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
	if claims, ok := pToken.Claims.(*auth.LoginClaims); ok {
		if claims.ActiveOrganization != "" {
			c.Console().Printf("Logged in as organization: %s (on behalf of user %s)\n", claims.ActiveOrganization, claims.AccountName)
		} else {
			c.Console().Printf("Logged in as user: %s\n", claims.AccountName)
		}
	}
	return nil
}
