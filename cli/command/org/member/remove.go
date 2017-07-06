package member

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type remMemOrgOptions struct {
	name string
}

// NewOrgRemoveMemCommand returns a new instance of the remove organization member command.
func NewOrgRemoveMemCommand(c cli.Interface) *cobra.Command {
	opts := remMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] MEMBER(S)",
		Short:   "Remove one or more members",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrgMem(c, cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.name, "org", "", "Organization name")
	return cmd
}

func removeOrgMem(c cli.Interface, cmd *cobra.Command, args []string, opts remMemOrgOptions) error {
	var errs []string
	org, err := cli.ReadOrg(c.Server())
	if !cmd.Flag("org").Changed {
		switch {
		case err == nil:
			opts.name = org
			c.Console().Println("organization name:", opts.name)
		default:
			opts.name = c.Console().GetInput("organization name")
		}
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, member := range args {
		request := &account.RemoveUserFromOrganizationRequest{
			OrganizationName: opts.name,
			UserName:         member,
		}
		if _, err := client.RemoveUserFromOrganization(context.Background(), request); err != nil {
			if s, ok := status.FromError(err); ok {
				errs = append(errs, s.Message())
				continue
			}
		}
		c.Console().Println(member)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
