package member

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type changeMemOrgOpts struct {
	name   string
	member string
	role   string
}

var (
	changeMemOrgOptions = &changeMemOrgOpts{}
)

// NewOrgChangeMemRoleCommand returns a new instance of the organization member role change command.
func NewOrgChangeMemRoleCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "change",
		Short:   "Change member role",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return changeOrgMemRole(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&changeMemOrgOptions.name, "org", "", "Organization name")
	flags.StringVar(&changeMemOrgOptions.member, "member", "", "Member name")
	flags.StringVar(&changeMemOrgOptions.role, "role", "", "Organization role")
	return cmd
}

func changeOrgMemRole(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		changeMemOrgOptions.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("member").Changed {
		changeMemOrgOptions.member = c.Console().GetInput("member name")
	}
	if !cmd.Flag("role").Changed {
		changeMemOrgOptions.role = c.Console().GetInput("organization role")
	}
	orgRole := accounts.OrganizationRole_ORGANIZATION_MEMBER
	switch changeMemOrgOptions.role {
	case "owner":
		orgRole = accounts.OrganizationRole_ORGANIZATION_OWNER
	case "member":
		orgRole = accounts.OrganizationRole_ORGANIZATION_MEMBER
	default:
		return fmt.Errorf("invalid organization role: %s. Please specify 'owner' or 'member' as role value.", changeMemOrgOptions.role)
	}
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ChangeOrganizationMemberRoleRequest{
		OrganizationName: changeMemOrgOptions.name,
		UserName:         changeMemOrgOptions.member,
		Role:             orgRole,
	}
	if _, err = client.ChangeOrganizationMemberRole(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Member role has been changed.")
	return nil
}
