package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
)

// StackRemove command to remove stack
var StackRemove = &cobra.Command{
	Use:     "rm STACK",
	Aliases: []string{"remove", "down", "stop"},
	Short:   "Remove the stack",
	Long:    `Remove the stack, stop all its services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stackRemove(AMP, cmd, args)
	},
}

func init() {
	StackCmd.AddCommand(StackRemove)
}

func stackRemove(amp *client.AMP, cmd *cobra.Command, args []string) error {
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

	stackName := args[0]
	request := &stack.StackRequest{Name: stackName}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Remove(ctx, request)
	if err != nil {
		manager.fatalf("%v\n", err)
	}

	manager.printf(colSuccess, "%s\n", reply.Answer)
	return nil
}
