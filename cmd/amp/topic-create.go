package main

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	createTopicCmd = &cobra.Command{
		Use:   "create [OPTIONS] NAME",
		Short: "Create a topic",
		Long:  `Create a topic.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createTopic(AMP, cmd, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(createTopicCmd)
}

func createTopic(amp *client.AMP, cmd *cobra.Command, args []string) error {
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
		return err
	}

	fmt.Println(reply.Topic.Id)
	return nil
}
