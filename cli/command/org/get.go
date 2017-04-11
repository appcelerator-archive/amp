package org

import (
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type getOrgOpts struct {
	name string
}

var (
	getOrgOptions = &getOrgOpts{}
)

// NewOrgGetCommand returns a new instance of the get organization command.
func NewOrgGetCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getOrg(c, cmd)
		},
	}
	cmd.Flags().StringVar(&getOrgOptions.name, "org", "", "Organization name")
	return cmd
}

func getOrg(c cli.Interface, cmd *cobra.Command) error {
	if !cmd.Flag("org").Changed {
		getOrgOptions.name = c.Console().GetInput("organization name")
	}
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: getOrgOptions.name,
	}
	reply, err := client.GetOrganization(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	c.Console().Printf("Organization: %s\n", reply.Organization.Name)
	c.Console().Printf("Email: %s\n", reply.Organization.Email)
	c.Console().Printf("Created: %s\n", time.ConvertTime(reply.Organization.CreateDt))
	return nil
}
