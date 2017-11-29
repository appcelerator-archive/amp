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

// NewForgetCommand returns a new instance of the forget command
func NewForgetCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "forget [flags] NAME or ID",
		Aliases: []string{"deregister"},
		Short:   "deregister object storage",
		PreRunE: cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return forget(c, cmd, args[0])
		},
	}
	return cmd
}

func forget(c cli.Interface, cmd *cobra.Command, name string) error {
	client := object_store.NewObjectStoreClient(c.ClientConn())
	req := &object_store.ForgetRequest{Name: name}
	reply, err := client.Forget(context.Background(), req)
	if err != nil {
		if s, ok := grpcStatus.FromError(err); ok {
			return errors.New(s.Message())
		}
	}
	fmt.Printf("Object store %s has been deregistered\n", reply.Name)
	return nil
}
