package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/appcelerator/amp/api/rpc/config"
	"github.com/appcelerator/amp/api/rpc/types"
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

// NewCreateCommand returns a new instance of the create command for creating a config.
func NewCreateCommand(c cli.Interface) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create [OPTIONS] CONFIG FILE|-",
		Short:   "Create a config from a file or STDIN as content",
		PreRunE: cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(c, cmd, args)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&createOpts.Labels, "labels", "l", []string{}, "Config labels")

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
			return fmt.Errorf("Error reading content from STDIN: %s", err.Error())
		}

	} else {
		data, err = ioutil.ReadFile(source)
		if err != nil {
			return fmt.Errorf("Error reading from file '%s': %s", source, err.Error())
		}

	}

	conn := c.ClientConn()
	client := config.NewConfigServiceClient(conn)
	spec := &config.ConfigSpec{
		Annotations: &types.Annotations{Name: name},
		Data:        data,
	}
	request := &config.CreateConfigRequest{
		Spec: spec,
	}
	resp, err := client.CreateConfig(context.Background(), request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return errors.New(s.Message())
		}
		return fmt.Errorf("Error creating config: %s", err)
	}
	fmt.Println(resp.GetConfig().GetId())

	return nil
}
