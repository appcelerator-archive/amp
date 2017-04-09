package org

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listOrgOpts struct {
	quiet bool
}

var (
	listOrgOptions = &listOrgOpts{}
)

// NewOrgListCommand returns a new instance of the list organization command.
func NewOrgListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Short:   "List organization",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrg(c, cmd)
		},
	}
	cmd.Flags().BoolVarP(&listOrgOptions.quiet, "quiet", "q", false, "Only display organization name")
	return cmd
}

func listOrg(c cli.Interface, cmd *cobra.Command) error {
	conn, err := c.ClientConn()
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ListOrganizationsRequest{}
	reply, err := client.ListOrganizations(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if listOrgOptions.quiet {
		for _, org := range reply.Organizations {
			c.Console().Println(org.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	for _, org := range reply.Organizations {
		fmt.Fprintf(w, "%s\t%s\t%s\n", org.Name, org.Email, cli.ConvertTime(c, org.CreateDt))
	}
	w.Flush()
	return nil
}
