package object_store

import (
	"errors"
	"fmt"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/object_store"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type CreateObjectStoreOptions struct {
	acl string
}

var (
	createOptions = CreateObjectStoreOptions{}
)

// NewCreateCommand returns a new instance of the object-store command.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [flags] NAME",
		Aliases: []string{"add", "new"},
		Short:   "Create an object store",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd, args[0])
		},
	}
	cmd.Flags().StringVar(&createOptions.acl, "acl", "private", "ACL (private, public-read, public-read-write, authenticated-read)")
	return cmd
}

func create(c cli.Interface, cmd *cobra.Command, name string) error {
	req := &object_store.CreateRequest{
		Name: name,
		Acl:  createOptions.acl,
	}

	client := object_store.NewObjectStoreClient(c.ClientConn())
	reply, err := client.Create(context.Background(), req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	w := tabwriter.NewWriter(c.Out(), 0, 0, cli.Padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tLOCATION")
	fmt.Fprintf(w, "%s\t%s\t%s\n", reply.Id, reply.Name, reply.Location)
	w.Flush()
	return nil
}
