package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack operations",
	Short: "Stack operations",
	Long:  `Manage stack-related operations.`,
}

var (
	upCmd = &cobra.Command{
		Use:   "up [-f FILE] [name]",
		Short: "Create and deploy a stack",
		Long:  `Create and deploy a stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			AMP.Connect()
			defer AMP.Disconnect()
			err := up(AMP, cmd, args)
			if err != nil {
				if AMP.Verbose() {
					log.Println(err)
				}
				log.Fatal("Failed to create and deploy stack")
			}
		},
	}
	// stack configuration file
	stackfile string
	startCmd  = &cobra.Command{
		Use:   "restart [stack name or id]",
		Short: "Start a stopped stack",
		Long:  `Start a stopped stack`,
		Run: func(cmd *cobra.Command, args []string) {
			AMP.Connect()
			defer AMP.Disconnect()
			err := start(AMP, cmd, args)
			if err != nil {
				if AMP.Verbose() {
					log.Println(err)
				}
				log.Fatal("Failed to start stack")
			}
		},
	}
	stopCmd = &cobra.Command{
		Use:   "stop [stack name or id]",
		Short: "Stop a stack",
		Long:  `Stop all services of a stack.`,
		Run: func(cmd *cobra.Command, args []string) {
			AMP.Connect()
			defer AMP.Disconnect()
			err := stop(AMP, cmd, args)
			if err != nil {
				if AMP.Verbose() {
					log.Println(err)
				}
				log.Fatal("Failed to stop stack")
			}
		},
	}
	rmCmd = &cobra.Command{
		Use:   "rm [stack name or id]",
		Short: "Remove a stack",
		Long:  `Remove a stack completly including ETCD data.`,
		Run: func(cmd *cobra.Command, args []string) {
			AMP.Connect()
			defer AMP.Disconnect()
			err := remove(AMP, cmd, args)
			if err != nil {
				if AMP.Verbose() {
					log.Println(err)
				}
				log.Fatal("Failed to remove stack")
			}
		},
	}
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "List available stacks",
		Long:  `List available stacks.`,
		Run: func(cmd *cobra.Command, args []string) {
			AMP.Connect()
			defer AMP.Disconnect()
			err := list(AMP, cmd, args)
			if err != nil {
				if AMP.Verbose() {
					log.Println(err)
				}
				log.Fatal("Failed to list stacks")
			}
		},
	}
	listQuiet  *bool
	listAll    *bool
	listLast   *int64
	listLatest *bool
)

func init() {
	RootCmd.AddCommand(StackCmd)
	flags := upCmd.Flags()
	flags.StringVarP(&stackfile, "file", "f", stackfile, "The name of the stackfile")
	rmCmd.Flags().BoolP("force", "f", false, "Remove the stack whatever condition")
	listQuiet = listCmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")
	listAll = listCmd.Flags().BoolP("all", "a", false, "Show all stacks (default shows just running)")
	listLast = listCmd.Flags().Int64P("last", "n", 0, "Show n last created stacks (includes all states)")
	listLatest = listCmd.Flags().BoolP("latest", "l", false, "Show the latest created stack (includes all states)")
	StackCmd.AddCommand(upCmd)
	StackCmd.AddCommand(startCmd)
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
		log.Fatal("Specify the stackfile with the --flag option")
	}

	if len(args) == 0 {
		log.Fatal("Must specify stack name")
	}
	name := args[0]
	if name == "" {
		log.Fatal("Must specify stack name")
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

	fmt.Println(reply.StackId)
	return nil
}

func start(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack id")
	}
	ident := args[0]
	if ident == "" {
		log.Fatal("Must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Start(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply.StackId)
	return nil
}

func stop(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack id")
	}
	ident := args[0]
	if ident == "" {
		log.Fatal("Must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Stop(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply.StackId)
	return nil
}

func remove(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack id")
	}
	ident := args[0]
	if ident == "" {
		log.Fatal("Must specify stack name or id")
	}

	force := false
	if cmd.Flag("force").Value.String() == "true" {
		force = true
	}
	for _, stackIdent := range args {
		request := &stack.RemoveRequest{
			StackIdent: stackIdent,
			Force:      force,
		}

		client := stack.NewStackServiceClient(amp.Conn)
		reply, err := client.Remove(context.Background(), request)
		if err != nil {
			return err
		}

		fmt.Println(reply.StackId)
	}
	return nil
}

func list(amp *client.AMP, cmd *cobra.Command, args []string) error {
	var limit = *listLast
	if *listLatest {
		limit = 1
	}
	request := &stack.ListRequest{
		All:   *listAll,
		Limit: limit,
	}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		return err
	}
	//Manage -q
	if *listQuiet {
		for _, info := range reply.List {
			fmt.Println(info.Id)
		}
		return nil
	}
	if reply == nil || len(reply.List) == 0 {
		fmt.Println("No stack is available")
		return nil
	}
	//Format output
	col1 := 10
	col2 := 20
	col3 := 10
	for _, info := range reply.List {
		if len(info.Name) > col1 {
			col1 = len(info.Name) + 2
		}
		if len(info.Id) > col2 {
			col2 = len(info.Id) + 2
		}
		if len(info.State) > col3 {
			col3 = len(info.State) + 2
		}
	}
	fmt.Printf("%s%s%s\n", col("NAME", col1), col("ID", col2), col("STATE", col3))
	fmt.Printf("%s%s%s\n", col("-", col1), col("-", col2), col("-", col3))
	for _, info := range reply.List {
		fmt.Printf("%s%s%s\n", col(info.Name, col1), col(info.Id, col2), col(info.State, col3))
	}
	return nil
}
