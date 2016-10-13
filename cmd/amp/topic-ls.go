package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"text/tabwriter"
)

const (
	padding = 3
)

var (
	listTopicCmd = &cobra.Command{
		Use:   "ls [OPTIONS]",
		Short: "List topics",
		Long:  `List topics.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := listTopic(AMP, cmd, args)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	TopicCmd.AddCommand(listTopicCmd)
}

func listTopic(amp *client.AMP, cmd *cobra.Command, args []string) error {
	request := &topic.ListRequest{}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tPartitions\tReplication Factor\t")
	for _, topic := range reply.Topics {
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t\n", topic.Id, topic.Name, topic.Partitions, topic.ReplicationFactor)
	}
	w.Flush()

	return nil
}
