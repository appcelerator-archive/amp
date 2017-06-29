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

type addMemOrgOptions struct {
	name string
}

// NewOrgAddMemCommand returns a new instance of the add organization member command.
func NewOrgAddMemCommand(c cli.Interface) *cobra.Command {
	opts := addMemOrgOptions{}
	cmd := &cobra.Command{
		Use:     "add [OPTIONS] MEMBER",
		Short:   "Add member to organization",
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addOrgMem(c, cmd, args, opts)
		},
	}
	cmd.Flags().StringVar(&opts.name, "org", "", "Organization name")
	return cmd
}

func addOrgMem(c cli.Interface, cmd *cobra.Command, args []string, opts addMemOrgOptions) error {
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
		request := &account.AddUserToOrganizationRequest{
			OrganizationName: opts.name,
			UserName:         member,
		}
		if _, err := client.AddUserToOrganization(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	if err := cli.SaveOrg(opts.name, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Member(s) have been added to organization.")
	return nil
}
