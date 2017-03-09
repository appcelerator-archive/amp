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
		Use:   "create TOPIC-NAME",
		Short: "Create a topic",
		Long:  `The create command creates a topic with specified name.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTopic(AMP, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(createTopicCmd)
}

func createTopic(amp *cli.AMP, args []string) (err error) {
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
	reply, er := client.Create(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}

	fmt.Println(reply.Topic.Id)
	return nil
}
