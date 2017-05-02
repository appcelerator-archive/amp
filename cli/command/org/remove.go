package org

import (
	"errors"
	"strings"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// NewOrgRemoveCommand returns a new instance of the remove organization command.
func NewOrgRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm ORGANIZATION(S)",
		Short:   "Remove organization",
		Aliases: []string{"remove"},
		PreRunE: cli.AtLeastArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeOrg(c, args)
		},
	}
}

func removeOrg(c cli.Interface, args []string) error {
	var errs []string
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	for _, org := range args {
		request := &account.DeleteOrganizationRequest{
			Name: org,
		}
		if _, err := client.DeleteOrganization(context.Background(), request); err != nil {
			errs = append(errs, grpc.ErrorDesc(err))
			continue
		}
		c.Console().Println(org)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
