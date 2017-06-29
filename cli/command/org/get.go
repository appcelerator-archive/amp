package org

import (
	"errors"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

// NewOrgGetCommand returns a new instance of the get organization command.
func NewOrgGetCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "get ORGANIZATION",
		Short:   "Get organization",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getOrg(c, args)
		},
	}
}

func getOrg(c cli.Interface, args []string) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: args[0],
	}
	reply, err := client.GetOrganization(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	c.Console().Printf("Organization: %s\n", reply.Organization.Name)
	c.Console().Printf("Email: %s\n", reply.Organization.Email)
	c.Console().Printf("Created On: %s\n", time.ConvertTime(reply.Organization.CreateDt))
	return nil
}
