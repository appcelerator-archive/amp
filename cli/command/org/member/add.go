package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addMemOrgOptions struct {
	name   string
	member string
}

// NewOrgAddMemCommand returns a new instance of the add organization member command.
func NewOrgAddMemCommand(c cli.Interface) *cobra.Command {
	opts := addMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS]",
		Short:   "Add member to organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMem(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "org", "", "Organization name")
	flags.StringVar(&opts.member, "member", "", "Member name")
	return cmd
}

func addOrgMem(c cli.Interface, cmd *cobra.Command, opts addMemOrgOptions) error {
	if !cmd.Flag("org").Changed {
		opts.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("member").Changed {
		opts.member = c.Console().GetInput("member name")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.AddUserToOrganizationRequest{
		OrganizationName: opts.name,
		UserName:         opts.member,
	}
	if _, err := client.AddUserToOrganization(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been added to organization.")
	return nil
}
