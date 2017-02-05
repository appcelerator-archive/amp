package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"strings"
)

// StackServices command to list stack services
var StackServices = &cobra.Command{
	Use:   "services [OPTIONS] STACK",
	Short: "List the services in the stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		return stackServices(AMP, cmd, args)
	},
}

func init() {
	flags := StackServices.Flags()
	flags.BoolP("quiet", "q", false, "Only display IDs")
	flags.StringP("filter", "f", "", "Filter output based on conditions provided")
	StackCmd.AddCommand(StackServices)
}

func stackServices(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.Connect(); err != nil {
		manager.fatalf("Amp server is not available\n")
	}
	if len(args) < 1 {
		manager.fatalf("Need stack name as first argument")
	}
	quiet := false
	if cmd.Flag("quiet").Value.String() == "true" {
		quiet = true
	}
	filter := cmd.Flag("filter").Value.String()

	stackName := args[0]
	request := &stack.StackServicesRequest{Name: stackName, Quiet: quiet, Filter: filter}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Services(ctx, request)
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	lines := strings.Split(reply.Answer, "\n")
	for i, line := range lines {
		if i == 0 {
			manager.printf(colRegular, "%s\n", line)
		} else {
			manager.printf(colSuccess, "%s\n", line)
		}
	}
	return nil
}
