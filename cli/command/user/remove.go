package user

import (
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
	removeUserOptions = &removeUserOpts{}
)

// NewRemoveUserCommand returns a new instance of the remove user command.
func NewRemoveUserCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove user",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeUser(c, cmd)
		},
	}
	cmd.Flags().StringVar(&removeUserOptions.username, "name", "", "User name")
	return cmd
}

func removeUser(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("name").Changed {
		removeUserOptions.username = c.Console().GetInput("username")
	}

	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.DeleteUserRequest{
		Name: removeUserOptions.username,
	}
	if _, err := client.DeleteUser(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("User removed.")
	return nil
}
