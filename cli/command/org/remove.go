package org

import (
	"errors"
	"fmt"

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
	return &cobra.Command{
		Use:     "rm ORGANIZATION",
		Short:   "Remove organization",
		Aliases: []string{"remove"},
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("organization name cannot be empty")
			}
			removeOrgOptions.name = args[0]
			return removeOrg(c, removeOrgOptions)
		},
	}
}

func removeOrg(c cli.Interface, opt *removeOrgOpts) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.DeleteOrganizationRequest{
		Name: opt.name,
	}
	if _, err := client.DeleteOrganization(context.Background(), request); err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	c.Console().Println("Organization has been removed.")
	return nil
}
