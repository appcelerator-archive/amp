package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
)

// PlatformPull is the main command for attaching platform subcommands.
var PlatformPull = &cobra.Command{
	Use:   "pull",
	Short: "Pull infrastructure stacks images locally or on the server host",
	Long:  `The pull command pulls on the server host or locally all the images available from AMP Infrastructure stack.`,
	Run: func(cmd *cobra.Command, args []string) {
		pullAMPImages(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformPull)
}

func pullAMPImages(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager("true")
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.Connect(); err != nil {
		if !amp.IsLocalhost() {
			manager.fatalf("Amp server is not available:\n", err)
		}
		if err := manager.connectDocker(); err != nil {
			manager.fatalf("Docker connect error: %v\n", err)
		}
		server := stack.NewServer(nil, manager.docker)
		request := &stack.StackRequest{Name: "local"}
		reply, err := server.GetImages(ctx, request)
		if err != nil {
			manager.fatalf("pull error: %v\n", err)
		}
		for _, name := range reply.Images {
			request.Name = name
			manager.printf(colSuccess, "pulling image %s: ", name)
			_, err := server.PullImage(ctx, request)
			if err != nil {
				manager.printf(colError, "pull error: %v\n", err)
			}
			manager.printf(colUser, "done\n")
		}
	} else {
		request := &stack.StackRequest{}
		client := stack.NewStackServiceClient(amp.Conn)
		reply, err := client.GetImages(ctx, request)
		if err != nil {
			manager.fatalf("pull error: %v\n", err)
		}
		for _, name := range reply.Images {
			request.Name = name
			manager.printf(colSuccess, "pulling image %s : ", name)
			_, err := client.PullImage(ctx, request)
			if err != nil {
				manager.printf(colError, "pull error: %v\n", err)
			}
			manager.printf(colUser, "done\n")
		}
	}
	return nil
}
