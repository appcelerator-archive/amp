package main

import (
	"fmt"

	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	removeTopicCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove topic",
		Example: "7gstrgfgv",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTopic(AMP, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(removeTopicCmd)
}

func removeTopic(amp *cli.AMP, args []string) error {
	if len(args) == 0 {
		mgr.Error("must specify topic id")
	}
	id := args[0]
	if id == "" {
		mgr.Error("must specify topic id")
	}

	request := &topic.DeleteRequest{Id: id}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.Delete(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}
	fmt.Println(reply.Topic.Id)
	return nil
}
