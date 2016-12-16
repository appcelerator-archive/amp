package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Stack operations",
	Long:  `Manage stack-related operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return AMP.Connect()
	},
}

var (
	stackCreateCmd = &cobra.Command{
		Use:   "create [-f FILE] [name]",
		Short: "Create a stack",
		Long:  `Create a stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackCreate(AMP, cmd, args)
		},
	}
	stackUpCmd = &cobra.Command{
		Use:   "up [-f FILE] [name]",
		Short: "Create and deploy a stack",
		Long:  `Create and deploy a stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackUp(AMP, cmd, args)
		},
	}
	// stack configuration file
	stackfile     string
	stackStartCmd = &cobra.Command{
		Use:   "start [stack name or id]",
		Short: "Start a stopped stack",
		Long:  `Start a stopped stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackStart(AMP, cmd, args)
		},
	}
	stackStopCmd = &cobra.Command{
		Use:   "stop [stack name or id]",
		Short: "Stop a stack",
		Long:  `Stop all services of a stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackStop(AMP, cmd, args)
		},
	}
	stackRmCmd = &cobra.Command{
		Use:   "rm [stack name or id]",
		Short: "Remove a stack",
		Long:  `Remove a stack completly including ETCD data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackRm(AMP, cmd, args)
		},
	}
	stackListCmd = &cobra.Command{
		Use:   "ls",
		Short: "List available stacks",
		Long:  `List available stacks.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackList(AMP, cmd, args)
		},
	}
	stackTasksCmd = &cobra.Command{
		Use:   "ps [stack name or id]",
		Short: "List the tasks of a stack",
		Long:  `List the tasks of a stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackTasks(AMP, cmd, args)
		},
	}
	stackUrlsCmd = &cobra.Command{
		Use:   "urls [stack name or id . . .]",
		Short: "List the urls for a stack",
		Long:  `List the urls for a stack.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackUrls(AMP, cmd, args)
		},
	}
	listQuiet  *bool
	listAll    *bool
	listLast   *int64
	listLatest *bool
)

func init() {
	RootCmd.AddCommand(StackCmd)
	stackCreateCmd.Flags().StringVarP(&stackfile, "file", "f", stackfile, "The name of the stackfile")
	stackCreateCmd.Flags().BoolVar(&registryAuth, "with-registry-auth", false, "Send registry authentication details to Swarm agents")
	stackUpCmd.Flags().StringVarP(&stackfile, "file", "f", stackfile, "The name of the stackfile")
	stackUpCmd.Flags().BoolVar(&registryAuth, "with-registry-auth", false, "Send registry authentication details to Swarm agents")
	stackRmCmd.Flags().BoolP("force", "f", false, "Remove the stack whatever condition")
	listQuiet = stackListCmd.Flags().BoolP("quiet", "q", false, "Only display numeric IDs")
	listAll = stackListCmd.Flags().BoolP("all", "a", false, "Show all stacks (default shows just running)")
	listLast = stackListCmd.Flags().Int64P("last", "n", 0, "Show n last created stacks (includes all states)")
	listLatest = stackListCmd.Flags().BoolP("latest", "l", false, "Show the latest created stack (includes all states)")
	StackCmd.AddCommand(stackCreateCmd)
	StackCmd.AddCommand(stackUpCmd)
	StackCmd.AddCommand(stackStartCmd)
	StackCmd.AddCommand(stackStopCmd)
	StackCmd.AddCommand(stackRmCmd)
	StackCmd.AddCommand(stackListCmd)
	StackCmd.AddCommand(stackTasksCmd)
	StackCmd.AddCommand(stackUrlsCmd)
}

func stackCreate(amp *client.AMP, cmd *cobra.Command, args []string) (err error) {
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
	ctx := context.Background()

	s, err := stack.ParseStackfile(ctx, contents)
	if err != nil {
		return err
	}
	s.Name = name

	request := &stack.StackFileRequest{Stack: s}

	// only send auth if flag was set
	if registryAuth {
		dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
		opts := cliflags.NewClientOptions()
		err = dockerCli.Initialize(opts)
		if err != nil {
			return err
		}
		for _, service := range s.Services {
			// Retrieve encoded auth token from the image reference
			encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, service.Image)
			if err != nil {
				return err
			}
			service.RegistryAuth = encodedAuth
		}
	}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Create(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(reply.StackId)
	return nil
}

func stackUp(amp *client.AMP, cmd *cobra.Command, args []string) error {
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
	ctx := context.Background()

	s, err := stack.ParseStackfile(ctx, contents)
	if err != nil {
		return err
	}
	s.Name = name

	request := &stack.StackFileRequest{Stack: s}

	// only send auth if flag was set
	if registryAuth {
		dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
		opts := cliflags.NewClientOptions()
		err = dockerCli.Initialize(opts)
		if err != nil {
			return err
		}
		for _, service := range s.Services {
			// Retrieve encoded auth token from the image reference
			encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, service.Image)
			if err != nil {
				return err
			}
			service.RegistryAuth = encodedAuth
		}
	}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Up(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(reply.StackId)
	return nil
}

func stackStart(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack name or id")
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

	fmt.Println(reply)
	return nil
}

func stackStop(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack name or id")
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

func stackRm(amp *client.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		log.Fatal("Must specify stack name or id")
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

func stackList(amp *client.AMP, cmd *cobra.Command, args []string) error {
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

func stackTasks(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		log.Fatal("Must specify stack name or id")
	}
	if len(args) > 1 {
		log.Fatal("Must specify single stack name or id")
	}
	ident := args[0]
	if ident == "" {
		log.Fatal("Must specify stack name or id")
	}
	request := &stack.TasksRequest{
		StackIdent: ident,
	}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Tasks(context.Background(), request)
	if err != nil {
		return err
	}
	fmt.Print(reply.Message)
	return nil
}

func stackUrls(amp *client.AMP, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		log.Fatal("Must specify stack name or id")
	}
	for _, ident := range args {
		if ident == "" {
			log.Fatal("Must specify stack name or id")
		}
		request := &stack.StackRequest{
			StackIdent: ident,
		}
		client := stack.NewStackServiceClient(amp.Conn)
		reply, err := client.Get(context.Background(), request)
		if err != nil {
			return err
		}
		for _, service := range reply.Stack.Services {
			for _, spec := range service.PublishSpecs {
				fmt.Printf("http://%s.%s.local.atomiq.io/\n", spec.Name, reply.Stack.Name)
			}
		}
	}
	return nil
}
