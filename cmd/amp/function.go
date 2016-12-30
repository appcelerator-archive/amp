package main

import (
	"errors"
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var (
	functionCmd = &cobra.Command{
		Use:     "function",
		Short:   "function operations",
		Aliases: []string{"fn"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	createFunctionCmd = &cobra.Command{
		Use:   "create NAME IMAGE",
		Short: "Create a function",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createFunction(AMP, cmd, args)
		},
	}

	listFunctionCmd = &cobra.Command{
		Use:   "ls",
		Short: "List functions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFunction(AMP, cmd, args)
		},
	}

	removeFunctionCmd = &cobra.Command{
		Use:   "rm FUNCTION",
		Short: "Remove a function",
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeFunction(AMP, cmd, args)
		},
	}
)

func init() {
	listFunctionCmd.Flags().BoolP("quiet", "q", false, "Only display IDs")

	functionCmd.AddCommand(createFunctionCmd)
	functionCmd.AddCommand(listFunctionCmd)
	functionCmd.AddCommand(removeFunctionCmd)
	RootCmd.AddCommand(functionCmd)
}

func createFunction(amp *client.AMP, cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("must specify function name and docker image")
	case 1:
		return errors.New("must specify docker image")
	case 2:
	// OK
	default:
		return errors.New("too many arguments")
	}

	name, image := strings.TrimSpace(args[0]), strings.TrimSpace(args[1])
	if name == "" {
		return errors.New("function name cannot be empty")
	}
	if image == "" {
		return errors.New("docker image cannot be empty")
	}

	// Create function
	request := &function.CreateRequest{Function: &function.FunctionEntry{
		Name:  name,
		Image: image,
	}}
	reply, err := function.NewFunctionClient(amp.Conn).Create(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply.Function.Id)
	return nil
}

func listFunction(amp *client.AMP, cmd *cobra.Command, args []string) error {
	// List functions
	request := &function.ListRequest{}
	reply, err := function.NewFunctionClient(amp.Conn).List(context.Background(), request)
	if err != nil {
		return err
	}

	// --quiet only display IDs
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		return fmt.Errorf("Unable to convert quiet parameter: %v", cmd.Flag("f").Value.String())
	} else if quiet {
		for _, fn := range reply.Functions {
			fmt.Println(fn.Id)
		}
		return nil
	}

	// Table view
	w := tabwriter.NewWriter(os.Stdout, 0, 0, tablePadding, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tImage")
	for _, fn := range reply.Functions {
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", fn.Id, fn.Name, fn.Image)
	}
	w.Flush()

	return nil
}

func removeFunction(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("rm requires at least one argument")
	}

	client := function.NewFunctionClient(amp.Conn)
	for _, arg := range args {
		if arg == "" {
			continue
		}

		request := &function.DeleteRequest{Id: arg}
		_, err := client.Delete(context.Background(), request)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(arg)
		}
	}

	return nil
}
