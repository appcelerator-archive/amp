package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"text/tabwriter"
)

// StackCmd is the main command for attaching stack subcommands.
var StackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Stack operations",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return AMP.Connect()
	},
}

var (
	stackCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a stack",
		Example: "amp stack create -f examples/stacks/micro/stack.yml micro-stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackCreate(AMP, cmd, args)
		},
	}
	stackUpCmd = &cobra.Command{
		Use:     "up",
		Short:   "Create and deploy a stack",
		Example: "amp stack up -f examples/stacks/micro/stack.yml micro-stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackUp(AMP, cmd, args)
		},
	}
	// stack configuration file
	stackfile     string
	stackStartCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start a stopped stack",
		Example: "amp stack start micro-stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackStart(AMP, args)
		},
	}
	stackStopCmd = &cobra.Command{
		Use:     "stop",
		Short:   "Stop a stack",
		Example: "amp stack stop micro-stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackStop(AMP, args)
		},
	}
	stackRmCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove a stack",
		Example: "amp stack rm micro-stack \namp stack del micro-stack",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackRm(AMP, cmd, args)
		},
	}
	stackListCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List available stacks",
		Example: "amp stack ls -q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackList(AMP)
		},
	}
	stackTasksCmd = &cobra.Command{
		Use:     "ps",
		Short:   "List the tasks of a stack",
		Example: "amp stack ps micro-stack",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackTasks(AMP, args)
		},
	}
	stackUrlsCmd = &cobra.Command{
		Use:     "urls",
		Short:   "List the urls for a stack",
		Example: "amp stack urls pinger",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stackUrls(AMP, args)
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

func stackCreate(amp *cli.AMP, cmd *cobra.Command, args []string) error {
	stackfile, err := cmd.Flags().GetString("file")
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	// TODO: note: currently --file is *not* an optional flag event though it's intended to be
	if stackfile == "" {
		mgr.Error("specify the stackfile with the --flag option")
	}

	if len(args) == 0 {
		mgr.Error("must specify stack name")
	}
	name := args[0]
	if name == "" {
		mgr.Error("must specify stack name")
	}

	b, err := ioutil.ReadFile(stackfile)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	contents := string(b)
	ctx := context.Background()

	s, err := stack.ParseStackfile(ctx, contents)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	s.Name = name

	request := &stack.StackFileRequest{Stack: s}

	// only send auth if flag was set
	if registryAuth {
		dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
		opts := cliflags.NewClientOptions()
		err = dockerCli.Initialize(opts)
		if err != nil {
			mgr.Error(grpc.ErrorDesc(err))
		}
		for _, service := range s.Services {
			// Retrieve encoded auth token from the image reference
			encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, service.Image)
			if err != nil {
				mgr.Error(grpc.ErrorDesc(err))
			}
			service.RegistryAuth = encodedAuth
		}
	}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Create(ctx, request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	fmt.Println(reply.StackId)
	return nil
}

func stackUp(amp *cli.AMP, cmd *cobra.Command, args []string) error {
	stackfile, err := cmd.Flags().GetString("file")
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	// TODO: note: currently --file is *not* an optional flag event though it's intended to be
	if stackfile == "" {
		mgr.Error("specify the stackfile with the --flag option")
	}

	if len(args) == 0 {
		mgr.Error("must specify stack name")
	}
	name := args[0]
	if name == "" {
		mgr.Error("must specify stack name")
	}

	b, err := ioutil.ReadFile(stackfile)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	contents := string(b)
	ctx := context.Background()

	s, err := stack.ParseStackfile(ctx, contents)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	s.Name = name

	request := &stack.StackFileRequest{Stack: s}

	// only send auth if flag was set
	if registryAuth {
		dockerCli := command.NewDockerCli(os.Stdin, os.Stdout, os.Stderr)
		opts := cliflags.NewClientOptions()
		err := dockerCli.Initialize(opts)
		if err != nil {
			mgr.Error(grpc.ErrorDesc(err))
		}
		for _, service := range s.Services {
			// Retrieve encoded auth token from the image reference
			encodedAuth, err := command.RetrieveAuthTokenFromImage(ctx, dockerCli, service.Image)
			if err != nil {
				mgr.Error(grpc.ErrorDesc(err))
			}
			service.RegistryAuth = encodedAuth
		}
	}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Up(ctx, request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	fmt.Println(reply.StackId)
	return nil
}

func stackStart(amp *cli.AMP, args []string) error {

	if len(args) == 0 {
		mgr.Error("must specify stack name or id")
	}
	ident := args[0]
	if ident == "" {
		mgr.Error("must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Start(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	fmt.Println(reply)
	return nil
}

func stackStop(amp *cli.AMP, args []string) error {

	if len(args) == 0 {
		mgr.Error("must specify stack name or id")
	}
	ident := args[0]
	if ident == "" {
		mgr.Error("must specify stack name or id")
	}

	request := &stack.StackRequest{StackIdent: ident}

	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Stop(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	fmt.Println(reply.StackId)
	return nil
}

func stackRm(amp *cli.AMP, cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		mgr.Error("must specify stack name or id")
	}
	ident := args[0]
	if ident == "" {
		mgr.Error("must specify stack name or id")
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
			mgr.Error(grpc.ErrorDesc(err))
		}

		fmt.Println(reply.StackId)
	}
	return nil
}

func stackList(amp *cli.AMP) error {
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
		mgr.Error(grpc.ErrorDesc(err))
	}
	//Manage -q
	if *listQuiet {
		for _, info := range reply.List {
			fmt.Println(info.Id)
		}
		return nil
	}
	if reply == nil || len(reply.List) == 0 {
		mgr.Warn("no stack available")
	}
	//Format output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tID\tSTATE")
	for _, info := range reply.List {
		fmt.Fprintf(w, "%s\t%s\t%s\n", info.Name, info.Id, info.State)
	}
	w.Flush()
	return nil
}

func stackTasks(amp *cli.AMP, args []string) error {
	if len(args) == 0 {
		//log.Fatal("Must specify stack name or id")
		mgr.Error("must specify stack name or id")
	}
	if len(args) > 1 {
		//log.Fatal("Must specify single stack name or id")
		mgr.Error("must specify only one stack name or id")
	}
	ident := args[0]
	if ident == "" {
		//log.Fatal("Must specify stack name or id")
		mgr.Error("must specify stack name or id")
	}
	request := &stack.TasksRequest{
		StackIdent: ident,
	}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Tasks(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	fmt.Print(reply.Message)
	return nil
}

func stackUrls(amp *cli.AMP, args []string) error {
	if len(args) == 0 {
		//log.Fatal("Must specify stack name or id")
		mgr.Error("must specify stack name or id")
	}
	for _, ident := range args {
		if ident == "" {
			//log.Fatal("Must specify stack name or id")
			mgr.Error("must specify stack name or id")
		}
		request := &stack.StackRequest{
			StackIdent: ident,
		}
		client := stack.NewStackServiceClient(amp.Conn)
		reply, err := client.Get(context.Background(), request)
		if err != nil {
			mgr.Error(grpc.ErrorDesc(err))
		}
		for _, service := range reply.Stack.Services {
			for _, spec := range service.PublishSpecs {
				fmt.Printf("http://%s.%s.local.atomiq.io/\n", spec.Name, reply.Stack.Name)
			}
		}
	}
	return nil
}
