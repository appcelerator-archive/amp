package user

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type removeUserOpts struct {
	username string
}

var (
	rmOptions = &removeUserOpts{}
)

// NewRemoveUserCommand returns a new instance of the remove user command.
func NewRemoveUserCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm USERNAME",
		Short:   "Remove user",
		Aliases: []string{"del"},
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("username cannot be empty")
			}
			rmOptions.username = args[0]
			return removeUser(c, rmOptions)
		},
	}
}

func removeUser(c cli.Interface, opt *removeUserOpts) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.DeleteUserRequest{
		Name: opt.username,
	}
	if _, err := client.DeleteUser(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("User removed.")
	return nil
}
