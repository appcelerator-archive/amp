package org

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type removeOrgOpts struct {
	name string
}

var (
	removeOrgOptions = &removeOrgOpts{}
)

// NewOrgRemoveCommand returns a new instance of the remove organization command.
func NewOrgRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove organization",
		Aliases: []string{"del"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrg(c, cmd)
		},
	}
	cmd.Flags().StringVar(&removeOrgOptions.name, "org", "", "Organization name")
	return cmd
}

func removeOrg(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		removeOrgOptions.name = c.Console().GetInput("organization name")
	}
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.DeleteOrganizationRequest{
		Name: removeOrgOptions.name,
	}
	if _, err = client.DeleteOrganization(context.Background(), request); err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Println("Organization has been removed.")
	return nil
}
