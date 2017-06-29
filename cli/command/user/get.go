package user

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// NewGetUserCommand returns a new instance of the get user command.
func NewGetUserCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "get USERNAME",
		Short:   "Get user",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getUser(c, args)
		},
	}
}

func getUser(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetUserRequest{
		Name: args[0],
	}
	reply, err := client.GetUser(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("Username: %s\n", reply.User.Name)
	c.Console().Printf("Email: %s\n", reply.User.Email)
	c.Console().Printf("Verified?: %t\n", reply.User.IsVerified)
	c.Console().Printf("Created On: %s\n", time.ConvertTime(reply.User.CreateDt))
	return nil
}
