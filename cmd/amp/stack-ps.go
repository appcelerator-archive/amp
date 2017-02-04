package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"strings"
)

// StackPs command to list stack tasks
var StackPs = &cobra.Command{
	Use:   "ps [OPTIONS] STACK",
	Short: "List the tasks in the stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		return stackPs(AMP, cmd, args)
	},
}

func init() {
	flags := StackPs.Flags()
	flags.Bool("no-trunc", false, "Do not truncate output")
	flags.Bool("no-resolve", false, "Do not map IDs to Names")
	flags.StringP("filter", "f", "", "Filter output based on conditions provided")
	StackCmd.AddCommand(StackPs)
}

func stackPs(amp *client.AMP, cmd *cobra.Command, args []string) error {
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
	noTrunc := false
	if cmd.Flag("no-trunc").Value.String() == "true" {
		noTrunc = true
	}
	noResolve := false
	if cmd.Flag("no-resolve").Value.String() == "true" {
		noResolve = true
	}
	filter := cmd.Flag("filter").Value.String()

	stackName := args[0]
	request := &stack.StackPsRequest{Name: stackName, NoTrunc: noTrunc, NoResolve: noResolve, Filter: filter}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.ListTasks(ctx, request)
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
