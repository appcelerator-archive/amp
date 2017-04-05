package member

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type addMemOrgOpts struct {
	name   string
	member string
}

var (
	addMemOrgOptions = &addMemOrgOpts{}
)

// NewOrgAddMemCommand returns a new instance of the add organization member command.
func NewOrgAddMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add member to organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&addMemOrgOptions.name, "org", "", "Organization name")
	flags.StringVar(&addMemOrgOptions.member, "member", "", "Member name")
	return cmd
}

func addOrgMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		addMemOrgOptions.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("member").Changed {
		addMemOrgOptions.member = c.Console().GetInput("member name")
	}
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.AddUserToOrganizationRequest{
		OrganizationName: addMemOrgOptions.name,
		UserName:         addMemOrgOptions.member,
	}
	if _, err = client.AddUserToOrganization(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been added to organization.")
	return nil
}
