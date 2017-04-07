package user

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listUserOpts struct {
	quiet bool
}

var (
	listUserOptions = &listUserOpts{}
)

// NewListUserCommand returns a new instance of the list user command.
func NewListUserCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List users",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listUser(c, cmd)
		},
	}
	cmd.Flags().BoolVarP(&listUserOptions.quiet, "quiet", "q", false, "Only display user names")
	return cmd
}

func listUser(c cli.Interface, cmd *cobra.Command) error {
	request := &account.ListUsersRequest{}
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	reply, err := client.ListUsers(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listUserOptions.quiet {
		for _, user := range reply.Users {
			c.Console().Println(user.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "USERNAME\tEMAIL")
	for _, user := range reply.Users {
		fmt.Fprintf(w, "%s\t%s\n", user.Name, user.Email)
	}
	w.Flush()
	return nil
}
