package member

import (
	"errors"

	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type changeMemOrgOptions struct {
	name   string
	member string
	role   string
}

// NewOrgChangeMemRoleCommand returns a new instance of the organization member role change command.
func NewOrgChangeMemRoleCommand(c cli.Interface) *cobra.Command {
	opts := changeMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "role [OPTIONS]",
		Short:   "Change member role",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return changeOrgMemRole(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "org", "", "Organization name")
	flags.StringVar(&opts.member, "member", "", "Member name")
	flags.StringVar(&opts.role, "role", "member", "Organization role (value can be 'member' or 'role')")
	return cmd
}

func changeOrgMemRole(c cli.Interface, cmd *cobra.Command, opts changeMemOrgOptions) error {
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
	if !cmd.Flag("member").Changed {
		opts.member = c.Console().GetInput("member name")
	}
	if !cmd.Flag("role").Changed {
		opts.role = c.Console().GetInput("organization role")
	}
	orgRole := accounts.OrganizationRole_ORGANIZATION_MEMBER
	switch opts.role {
	case "owner":
		orgRole = accounts.OrganizationRole_ORGANIZATION_OWNER
	case "member":
		orgRole = accounts.OrganizationRole_ORGANIZATION_MEMBER
	default:
		return fmt.Errorf("invalid organization role: %s. Please specify 'owner' or 'member' as role value.", opts.role)
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ChangeOrganizationMemberRoleRequest{
		OrganizationName: opts.name,
		UserName:         opts.member,
		Role:             orgRole,
	}
	if _, err := client.ChangeOrganizationMemberRole(context.Background(), request); err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	if err := cli.SaveOrg(opts.name, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Member role has been changed.")
	return nil
}
