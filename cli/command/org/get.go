package org

import (
	"errors"
	"fmt"

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
	return &cobra.Command{
		Use:     "get ORGANIZATION",
		Short:   "Get organization",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("organization name cannot be empty")
			}
			getOrgOptions.name = args[0]
			return getOrg(c, getOrgOptions)
		},
	}
}

func getOrg(c cli.Interface, opt *getOrgOpts) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.GetOrganizationRequest{
		Name: opt.name,
	}
	reply, err := client.GetOrganization(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Printf("Organization: %s\n", reply.Organization.Name)
	c.Console().Printf("Email: %s\n", reply.Organization.Email)
	c.Console().Printf("Created: %s\n", time.ConvertTime(reply.Organization.CreateDt))
	return nil
}
