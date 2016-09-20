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

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack operations",
	Short: "stack operations",
	Long:  `Manage stack-related operations.`,
}

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
	stopCmd   = &cobra.Command{
		Use:   "stop [stack name or id]",
		Short: "Stop a stack",
		Long:  `Stop all services of a stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := stop(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	rmCmd = &cobra.Command{
		Use:   "rm [stack name or id]",
		Short: "Remove a stack",
		Long:  `Remove a stack completly including ETCD data.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := remove(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List available stacks",
		Long:  `List available stacks.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := list(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(StackCmd)
	flags := upCmd.Flags()
	flags.StringVarP(&stackfile, "file", "f", stackfile, "the name of the stackfile")
	StackCmd.AddCommand(upCmd)
	StackCmd.AddCommand(stopCmd)
	StackCmd.AddCommand(rmCmd)
	StackCmd.AddCommand(listCmd)
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

func stop(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		return errors.New("must specify stack id")
	}
	ident := args[0]
	if ident == "" {
		return errors.New("must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Stop(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)
	return nil
}

func remove(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		return errors.New("must specify stack id")
	}
	ident := args[0]
	if ident == "" {
		return errors.New("must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Remove(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply)
	return nil
}

func list(amp *client.AMP, cmd *cobra.Command, args []string) error {
	request := &stack.ListRequest{}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		return err
	}
	if reply == nil || len(reply.List) == 0 {
		fmt.Println("No stack is avaialble")
		return nil
	}
	fmt.Println("Stack list")
	for _, info := range reply.List {
		fmt.Println(info)
	}
	return nil
}
