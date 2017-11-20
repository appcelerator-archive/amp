package secret

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/appcelerator/amp/api/rpc/secret"
	"github.com/appcelerator/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

type CreateOpts struct {
	Labels []string
}

var createOpts = &CreateOpts{
	Labels: []string{},
}

// NewCreateCommand returns a new instance of the create command for creating a secret.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS] SECRET FILE|-",
		Short:   "Create a secret from a file or STDIN as content",
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&createOpts.Labels, "labels", "l", []string{}, "Secret labels")

	return cmd
}

func create(c cli.Interface, cmd *cobra.Command, args []string) error {
	name := args[0]
	source := args[1]

	var data []byte
	var err error

	if source == "-" {
		data, err = ioutil.ReadAll(c.In())
		if err != nil {
			return fmt.Errorf("error reading content from STDIN: %s", err.Error())
		}

	} else {
		data, err = ioutil.ReadFile(source)
		if err != nil {
			return fmt.Errorf("error reading from file '%s': %s", source, err.Error())
		}

	}

	conn := c.ClientConn()
	client := secret.NewSecretClient(conn)
	request := &secret.CreateRequest{
		Name: name,
		Data: data,
	}
	reply, err := client.SecretCreate(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
		return fmt.Errorf("error creating secret: %s", err)
	}
	fmt.Println(reply.GetId())

	return nil
}
