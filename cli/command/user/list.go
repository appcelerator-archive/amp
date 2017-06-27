package user

import (
	"errors"
	"text/tabwriter"

	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type listUserOptions struct {
	quiet bool
}

// NewListUserCommand returns a new instance of the list user command.
func NewListUserCommand(c cli.Interface) *cobra.Command {
	opts := listUserOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List users",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(c, opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only display user names")
	return cmd
}

func listUser(c cli.Interface, opts listUserOptions) error {
	request := &account.ListUsersRequest{}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	reply, err := client.ListUsers(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if opts.quiet {
		for _, user := range reply.Users {
			c.Console().Println(user.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL\tCREATED ON")
	for _, user := range reply.Users {
		fmt.Fprintf(w, "%s\t%s\t%s\n", user.Name, user.Email, time.ConvertTime(user.CreateDt))
	}
	w.Flush()
	return nil
}
