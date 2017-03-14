package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var (
	functionCmd = &cobra.Command{
		Use:     "function",
		Short:   "Function operations",
		Aliases: []string{"fn"},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return AMP.Connect()
		},
	}

	createFunctionCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a function",
		Example: "sample-func samples/function-test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createFunction(AMP, args)
		},
	}

	listFunctionCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List functions",
		Example: "-q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listFunction(AMP, cmd)
		},
	}

	removeFunctionCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove a function",
		Aliases: []string{"del"},
		Example: "ujyhjdb656",
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeFunction(AMP, args)
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

func createFunction(amp *cli.AMP, args []string) error {
	switch len(args) {
	case 0:
		//return errors.New("must specify function name and docker image")
		mgr.Fatal("must specify function name and docker image")
	case 1:
		//return errors.New("must specify docker image")
		mgr.Fatal("must specify docker image")
	case 2: // OK
	default:
		//return errors.New("too many arguments")
		mgr.Fatal("too many arguments")
	}

	name, image := strings.TrimSpace(args[0]), strings.TrimSpace(args[1])
	if name == "" {
		//return errors.New("function name cannot be empty")
		mgr.Fatal("function name cannot be empty")
	}
	if image == "" {
		//return errors.New("docker image cannot be empty")
		mgr.Fatal("docker image cannot be empty")
	}

	// Create function
	request := &function.CreateRequest{
		Name:  name,
		Image: image,
	}
	reply, err := function.NewFunctionClient(amp.Conn).Create(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}

	fmt.Println(reply.Function.Id)
	return nil
}

func listFunction(amp *cli.AMP, cmd *cobra.Command) error {
	// List functions
	request := &function.ListRequest{}
	reply, err := function.NewFunctionClient(amp.Conn).List(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}

	// --quiet only display IDs
	if quiet, err := strconv.ParseBool(cmd.Flag("quiet").Value.String()); err != nil {
		mgr.Fatal("Unable to convert quiet parameter: %v", cmd.Flag("f").Value.String())
	} else if quiet {
		for _, fn := range reply.Functions {
			fmt.Println(fn.Id)
		}
		return nil
	}

	// Table view
	w := tabwriter.NewWriter(os.Stdout, 0, 0, tablePadding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tIMAGE\tOWNER")
	for _, fn := range reply.Functions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", fn.Id, fn.Name, fn.Image, fn.Owner.Name)
	}
	w.Flush()

	return nil
}

func removeFunction(amp *cli.AMP, args []string) error {
	if len(args) == 0 {
		mgr.Fatal("rm requires at least one argument")
	}

	client := function.NewFunctionClient(amp.Conn)
	for _, arg := range args {
		if arg == "" {
			continue
		}

		request := &function.DeleteRequest{Id: arg}
		_, err := client.Delete(context.Background(), request)
		if err != nil {
			mgr.Fatal(grpc.ErrorDesc(err))
		} else {
			fmt.Println(arg)
		}
	}

	return nil
}
