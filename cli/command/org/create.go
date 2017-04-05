package org

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type createOrgOpts struct {
	name  string
	email string
}

var (
	createOrgOptions = &createOrgOpts{}
)

// NewOrgCreateCommand returns a new instance of the create organization command.
func NewOrgCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createOrg(c, cmd)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&createOrgOptions.name, "org", "", "Organization name")
	flags.StringVar(&createOrgOptions.email, "email", "", "Email address")
	return cmd
}

func createOrg(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		createOrgOptions.name = c.Console().GetInput("organization name")
	}
	if !cmd.Flag("email").Changed {
		createOrgOptions.email = c.Console().GetInput("email")
	}
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.CreateOrganizationRequest{
		Name:  createOrgOptions.name,
		Email: createOrgOptions.email,
	}
	if _, err = client.CreateOrganization(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("Organization has been created.")
	return nil
}
