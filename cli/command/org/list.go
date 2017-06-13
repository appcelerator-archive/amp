package org

import (
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type listOrgOptions struct {
	quiet bool
}

// NewOrgListCommand returns a new instance of the list organization command.
func NewOrgListCommand(c cli.Interface) *cobra.Command {
	opts := listOrgOptions{}
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Short:   "List organization",
		Aliases: []string{"list"},
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listOrg(c, opts)
		},
	}
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Only display organization names")
	return cmd
}

func listOrg(c cli.Interface, opts listOrgOptions) error {
	conn := c.ClientConn()
	client := account.NewAccountClient(conn)
	request := &account.ListOrganizationsRequest{}
	reply, err := client.ListOrganizations(context.Background(), request)
	if err != nil {
		return fmt.Errorf("%s", grpc.ErrorDesc(err))
	}
	if opts.quiet {
		for _, org := range reply.Organizations {
			c.Console().Println(org.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED ON")
	for _, org := range reply.Organizations {
		fmt.Fprintf(w, "%s\t%s\t%s\n", org.Name, org.Email, time.ConvertTime(org.CreateDt))
	}
	w.Flush()
	return nil
}
