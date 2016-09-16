package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	upCmd = &cobra.Command{
		Use:   "up [-f FILE] [name]",
		Short: "Create and deploy a stack",
		Long:  `Create and deploy a stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := up(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	// stack configuration file
	stackfile string
)

func init() {
	flags := upCmd.Flags()
	flags.StringVarP(&stackfile, "file", "f", stackfile, "the name of the stackfile")

	StackCmd.AddCommand(upCmd)
}

func up(amp *client.AMP, cmd *cobra.Command, args []string) error {
	stackfile, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	// TODO: note: currently --file is *not* an optional flag event though it's intended to be
	if stackfile == "" {
		return errors.New("specify the stackfile with the --flag option")
	}

	if len(args) == 0 {
		return errors.New("must specify stack name")
	}
	name := args[0]
	if name == "" {
		return errors.New("must specify stack name")
	}

	b, err := ioutil.ReadFile(stackfile)
	if err != nil {
		return err
	}

	contents := string(b)
	request := &stack.UpRequest{StackName: name, Stackfile: contents}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Up(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)
	return nil
}
