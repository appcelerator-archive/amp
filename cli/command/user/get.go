package user

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type getUserOpts struct {
	username string
}

var (
	getOptions = &getUserOpts{}
)

// NewGetUserCommand returns a new instance of the get user command.
func NewGetUserCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "get USERNAME",
		Short:   "Get user",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("username cannot be empty")
			}
			getOptions.username = args[0]
			return getUser(c, getOptions)
		},
	}
}

func getUser(c cli.Interface, opt *getUserOpts) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.GetUserRequest{
		Name: opt.username,
	}
	reply, err := client.GetUser(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Username: %s\n", reply.User.Name)
	c.Console().Printf("Email: %s\n", reply.User.Email)
	c.Console().Printf("Verified?: %t\n", reply.User.IsVerified)
	c.Console().Printf("Create Date: %s\n", time.ConvertTime(reply.User.CreateDt))
	return nil
}
