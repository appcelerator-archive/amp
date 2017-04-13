package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type remMemOrgOpts struct {
	name   string
	member string
}

var (
	remMemOrgOptions = &remMemOrgOpts{}
)

// NewOrgRemoveMemCommand returns a new instance of the remove organization member command.
func NewOrgRemoveMemCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [OPTIONS]",
		Short:   "Remove member from organization",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrgMem(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&remMemOrgOptions.name, "org", "", "Organization name")
	flags.StringVar(&remMemOrgOptions.member, "member", "", "Member name")
	return cmd
}

func removeOrgMem(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		remMemOrgOptions.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("member").Changed {
		remMemOrgOptions.member = c.Console().GetInput("member name")
	}

	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.RemoveUserFromOrganizationRequest{
		OrganizationName: remMemOrgOptions.name,
		UserName:         remMemOrgOptions.member,
	}
	if _, err := client.RemoveUserFromOrganization(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member has been removed from organization.")
	return nil
}
