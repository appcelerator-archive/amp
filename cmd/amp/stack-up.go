package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
)

var (
	upCmd = &cobra.Command{
		Use:   "up [-f FILE]",
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

	fmt.Printf("stackfile: %s\n", stackfile)
	b, err := ioutil.ReadFile(stackfile)
	if err != nil {
		return err
	}

	contents := string(b)
	fmt.Println(contents)

	request := &stack.UpRequest{Stackfile: contents}

	client := stack.NewStackClient(amp.Conn)
	reply, err := client.Up(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)

	return nil
}
