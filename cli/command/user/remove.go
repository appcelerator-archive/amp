package user

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// NewRemoveUserCommand returns a new instance of the remove user command.
func NewRemoveUserCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm USERNAME(S)",
		Short:   "Remove one or more users",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeUser(c, args)
		},
	}
}

func removeUser(c cli.Interface, args []string) error {
	var errs []string
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	token, err := cli.ReadToken(c.Server())
	if err != nil {
		return errors.New("you are not logged in. Use `amp login` or `amp user signup`.")
	}
	pToken, _ := jwt.ParseWithClaims(token, &auth.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte{}, nil
	})
	claims, ok := pToken.Claims.(*auth.AuthClaims)
	if !ok {
		return errors.New("you are not logged in. Use `amp login` or `amp user signup`.")
	}
	for _, name := range args {
		request := &account.DeleteUserRequest{
			Name: name,
		}
		if _, err := client.DeleteUser(context.Background(), request); err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
		if name == claims.AccountName {
			if err := cli.RemoveToken(c.Server()); err != nil {
				errs = append(errs, err.Error())
			}
		}
		c.Console().Println(name)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	if err := cli.RemoveFile(c.Server()); err != nil {
		return err
	}
	return nil
}
