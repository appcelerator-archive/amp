package user

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewRemoveUserCommand returns a new instance of the remove user command.
func NewRemoveUserCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm USERNAME(S)",
		Short:   "Remove user",
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
	for _, name := range args {
		request := &account.DeleteUserRequest{
			Name: name,
		}
		if _, err := client.DeleteUser(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(name)
	}
	if err := cli.RemoveToken(); err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
