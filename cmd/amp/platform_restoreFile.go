package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
)

// PlatformRestoreFile is the main command for attaching platform subcommands.
var PlatformRestoreFile = &cobra.Command{
	Use:   "restoreFile [STACK]",
	Short: "Restore on the server the previous compose stack file or a variables file saved",
	Long:  `Restore on the server the previous compose stack file or a variables file saved`,
	Run: func(cmd *cobra.Command, args []string) {
		restoreStackFile(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformRestoreFile)
}

func restoreStackFile(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if len(args) == 0 {
		manager.fatalf("First argument should be an infrastructure stack name or 'variable' to send a variables file path\n")
	}
	fileName := args[0]

	if err := amp.Connect(); err != nil {
		manager.fatalf("Server not yet ready\n")
	}
	request := &stack.StackRestoreRequest{FileName: fileName}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.RestoreFile(ctx, request)
	if err != nil {
		manager.fatalf("update error: %v\n", err)
	}
	manager.printf(colSuccess, "%s\n", reply.Answer)
	return nil
}
