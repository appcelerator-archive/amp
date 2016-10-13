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
		Run: func(cmd *cobra.Command, args []string) {
			err := createTopic(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	partitions        uint64
	replicationFactor uint64
)

func init() {
	flags := createCmd.Flags()
	flags.Uint64Var(&partitions, "partitions", 1, "Number of partitions")
	flags.Uint64Var(&replicationFactor, "replicationFactor", 1, "Replication factor")
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
		Name:              name,
		Partitions:        partitions,
		ReplicationFactor: replicationFactor,
	}}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.Create(context.Background(), request)
	if err != nil {
		return err
	}

	fmt.Println(reply.Topic.Id)
	return nil
}
