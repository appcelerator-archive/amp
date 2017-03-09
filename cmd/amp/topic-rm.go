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
	removeTopicCmd = &cobra.Command{
		Use:     "rm",
		Short:   "Remove topic",
		Example: "amp topic rm 7gstrgfgv \namp topic del 7gstrgfgv",
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeTopic(AMP, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(removeTopicCmd)
}

func removeTopic(amp *cli.AMP, args []string) (err error) {
	if len(args) == 0 {
		return errors.New("must specify topic id")
	}
	id := args[0]
	if id == "" {
		return errors.New("must specify topic id")
	}

	request := &topic.DeleteRequest{Id: id}

	client := topic.NewTopicClient(amp.Conn)
	reply, er := client.Delete(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
	}
	fmt.Println(reply.Topic.Id)
	return nil
}
