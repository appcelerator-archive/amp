package member

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remMemOrgOptions struct {
	name string
}

// NewOrgRemoveMemCommand returns a new instance of the remove organization member command.
func NewOrgRemoveMemCommand(c cli.Interface) *cobra.Command {
	opts := remMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS] MEMBER(S)",
		Short:   "Remove member from organization",
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
	if !cmd.Flag("org").Changed {
		opts.name = c.Console().GetInput("organization name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, member := range args {
		request := &account.RemoveUserFromOrganizationRequest{
			OrganizationName: opts.name,
			UserName:         member,
		}
		if _, err := client.RemoveUserFromOrganization(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(member)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
