package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// PlatformStop is the main command for attaching platform subcommands.
var PlatformStop = &cobra.Command{
	Use:     "rm infra-STACK",
	Aliases: []string{"remove", "down", "stop"},
	Short:   "Remove infrastruccture stacks, one by one or all",
	Long:    `The stop command stops all infrastructure stacks (default) or a dedicate one.`,
	Run: func(cmd *cobra.Command, args []string) {
		stopAMP(AMP, cmd, args)
	},
}

func init() {
	PlatformStop.Flags().BoolP("quiet", "q", false, "Suppress terminal output")
	PlatformCmd.AddCommand(PlatformStop)
}

func stopAMP(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	stackName := "ampcore"
	if len(args) >= 1 {
		stackName = args[0]
	}
	if stackName == "ampcore" {
		if err := amp.Connect(); err == nil {
			for _, name := range stack.InfraShortStackList {
				if !isStackAlreadyStopped(ctx, amp, manager, name) {
					stopStack(ctx, amp, manager, name)
				}
			}
		}
		if err := manager.connectDocker(); err != nil {
			manager.fatalf("Docker connect error: %v\n", err)
		}
		manager.printf(colRegular, "removing stack: ampcore\n")
		server := stack.NewServer(nil, manager.docker)
		request := &stack.StackRequest{Name: stackName}
		reply, err := server.Remove(ctx, request)
		if err != nil {
			manager.fatalf("remove ampCore error: %v\n", err)
		}
		manager.printf(colSuccess, reply.Answer)
		return nil
	}
	if err := amp.Connect(); err != nil {
		manager.fatalf("Amp server is not available\n", err)
	}
	return stopStack(ctx, amp, manager, stackName)
}

func stopStack(ctx context.Context, amp *client.AMP, manager *ampManager, stackName string) error {
	manager.printf(colRegular, "removing stack: %s\n", stackName)
	request := &stack.StackRequest{Name: stackName}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Remove(ctx, request)
	if err != nil {
		manager.printf(colError, "remove %s error: %v\n", stackName, err)
		return err
	}
	manager.printf(colSuccess, reply.Answer)
	return nil
}

func isStackAlreadyStopped(ctx context.Context, amp *client.AMP, manager *ampManager, stackName string) bool {
	request := &stack.StackRequest{Name: stackName}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.GetStackStatus(ctx, request)
	if err != nil {
		return false
	}
	if reply.Answer == "stopped" {
		return true
	}
	return false
}
