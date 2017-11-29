package object_store

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/object_store"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	grpcStatus "google.golang.org/grpc/status"
)

type RemoveObjectStoreOptions struct {
	force bool
}

var (
	removeOptions = RemoveObjectStoreOptions{}
)

// NewRemoveCommand returns a new instance of the forget command
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove [flags] NAME or ID",
		Aliases: []string{"rm", "delete", "del"},
		Short:   "delete object storage",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, cmd, args[0])
		},
	}
	cmd.Flags().BoolVar(&removeOptions.force, "force", false, "Force remove non empty object stores")
	return cmd
}

func remove(c cli.Interface, cmd *cobra.Command, name string) error {
	client := object_store.NewObjectStoreClient(c.ClientConn())
	req := &object_store.RemoveRequest{Name: name, Force: removeOptions.force}
	reply, err := client.Remove(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	fmt.Printf("Object store %s has been removed\n", reply.Name)
	return nil
}
