package org

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
	"github.com/appcelerator/amp/pkg/time"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	cmd.Flags().BoolP("quiet", "q", false, "Only display organization name")
	return cmd
}

func listOrg(c cli.Interface, cmd *cobra.Command) error {
	conn, err := c.ClientConn()
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	client := account.NewAccountClient(conn)
	request := &account.ListOrganizationsRequest{}
	reply, err := client.ListOrganizations(context.Background(), request)
	if err != nil {
		c.Console().Fatalf(grpc.ErrorDesc(err))
	}
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		c.Console().Fatalf("unable to convert quiet parameter : %v", grpc.ErrorDesc(err))
	} else if quiet {
		for _, org := range reply.Organizations {
			c.Console().Println(org.Name)
		}
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', 0)
	fmt.Fprintln(w, "ORGANIZATION\tEMAIL\tCREATED\t")
	for _, org := range reply.Organizations {
		fmt.Fprintf(w, "%s\t%s\t%s\n", org.Name, org.Email, time.ConvertTime(org.CreateDt))
	}
	w.Flush()
	return nil
}
