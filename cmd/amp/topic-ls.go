package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const (
	padding = 3
)

var (
	listTopicCmd = &cobra.Command{
		Use:   "ls",
		Short: "List topics",
		Long:  `The list command returns all available topics.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTopic(AMP, cmd, args)
		},
	}
)

func init() {
	TopicCmd.AddCommand(listTopicCmd)
}

func listTopic(amp *cli.AMP, cmd *cobra.Command, args []string) error {
	request := &topic.ListRequest{}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\t")
	fmt.Fprintln(w, "--\t----\t")
	for _, topic := range reply.Topics {
		fmt.Fprintf(w, "%s\t%s\t\n", topic.Id, topic.Name)
	}
	w.Flush()

	return nil
}
