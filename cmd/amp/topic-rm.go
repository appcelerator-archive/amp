package main

import (
	"errors"
	"fmt"

	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	removeTopicCmd = &cobra.Command{
		Use:   "rm TOPIC-ID",
		Short: "Remove topic",
		Long:  `The remove command removes the specified topic id.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTopic(AMP, cmd, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(removeTopicCmd)
}

func removeTopic(amp *cli.AMP, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("must specify topic id")
	}
	id := args[0]
	if id == "" {
		return errors.New("must specify topic id")
	}

	request := &topic.DeleteRequest{Id: id}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.Delete(context.Background(), request)
	if err != nil {
		return err
	}
	fmt.Println(reply.Topic.Id)
	return nil
}
