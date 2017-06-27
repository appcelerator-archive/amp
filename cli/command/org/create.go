package org

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type createOrgOptions struct {
	name  string
	email string
}

// NewOrgCreateCommand returns a new instance of the create organization command.
func NewOrgCreateCommand(c cli.Interface) *cobra.Command {
	opts := createOrgOptions{}
	cmd := &cobra.Command{
		Use:     "create [OPTIONS]",
		Short:   "Create organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrg(c, cmd, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "org", "", "Organization name")
	flags.StringVar(&opts.email, "email", "", "Email address")
	return cmd
}

func createOrg(c cli.Interface, cmd *cobra.Command, opts createOrgOptions) error {
	if !cmd.Flag("org").Changed {
		opts.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("email").Changed {
		opts.email = c.Console().GetInput("email")
	}
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.CreateOrganizationRequest{
		Name:  opts.name,
		Email: opts.email,
	}
	if _, err := client.CreateOrganization(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if err := cli.SaveOrg(opts.name, c.Server()); err != nil {
		return err
	}
	c.Console().Println("Organization has been created.")
	return nil
}
