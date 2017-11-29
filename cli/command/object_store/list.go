package object_store

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/object_store"
	"github.com/appcelerator/amp/cli"
	//"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	grpcStatus "google.golang.org/grpc/status"
)

// NewListCommand returns a new instance of the list command
func NewListCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list object storage",
		PreRunE: cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(c, cmd)
		},
	}
	return cmd
}

func list(c cli.Interface, cmd *cobra.Command) error {
	client := object_store.NewObjectStoreClient(c.ClientConn())
	req := &object_store.ListRequest{}
	reply, err := client.List(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	fmt.Printf("Found %d object stores\n", len(reply.Entries))
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tOWNER\tLOCATION\tREGION\t")
	for _, o := range reply.Entries {
		var e string
		if o.Missing {
			e = "Error: not found, you can ask AMP to forget it"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", o.ObjectStore.Id, o.ObjectStore.Name, o.ObjectStore.Owner.User, o.ObjectStore.Location, o.Region, e)
	}
	w.Flush()
	return nil
}
