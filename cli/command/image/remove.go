package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/image"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

// NewRemoveCommand returns a new instance of the image remove command
func NewRemoveCommand(c cli.Interface) *cobra.Command {
	return &cobra.Command{
		Use:     "rm REPO DIGEST",
		Short:   "Remove an image",
		Aliases: []string{"remove", "del", "delete"},
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(c, cmd, args)
		},
	}
}

func remove(c cli.Interface, cmd *cobra.Command, args []string) error {
	conn := c.ClientConn()
	client := image.NewImageClient(conn)
	request := &image.RemoveRequest{Name: args[0], Digest: args[1]}
	_, err := client.ImageRemove(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New("ImageRemove API call failed: " + s.Message())
		}
		return fmt.Errorf("error removing image: %s", err)
	}
	fmt.Printf("Image successfully removed")
	return nil
}
