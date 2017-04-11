package user

import (
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
	getUserOptions = &getUserOpts{}
)

// NewGetUserCommand returns a new instance of the get user command.
func NewGetUserCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get user",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getUser(c, cmd)
		},
	}
	cmd.Flags().StringVar(&getUserOptions.username, "name", "", "User name")
	return cmd
}

func getUser(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("name").Changed {
		getUserOptions.username = c.Console().GetInput("username")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.GetUserRequest{
		Name: getUserOptions.username,
	}
	reply, err := client.GetUser(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Printf("Username: %s\n", reply.User.Name)
	c.Console().Printf("Email: %s\n", reply.User.Email)
	c.Console().Printf("Verified?: %t\n", reply.User.IsVerified)
	c.Console().Printf("Create Date: %s\n", time.ConvertTime(reply.User.CreateDt))
	return nil
}
