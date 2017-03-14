package main

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	createTopicCmd = &cobra.Command{
		Use:     "create",
		Short:   "Create a topic",
		Example: "dockerize",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTopic(AMP, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(createTopicCmd)
}

func createTopic(amp *cli.AMP, args []string) error {
	if len(args) == 0 {
		return errors.New("must specify topic name")
	}
	name := args[0]
	if name == "" {
		return errors.New("must specify topic name")
	}

	request := &topic.CreateRequest{Topic: &topic.TopicEntry{
		Name: name,
	}}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.Create(context.Background(), request)
	if err != nil {
		mgr.Fatal(grpc.ErrorDesc(err))
	}

	fmt.Println(reply.Topic.Id)
	return nil
}
