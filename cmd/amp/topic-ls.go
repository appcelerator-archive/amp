package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/appcelerator/amp/api/rpc/topic"
	"github.com/appcelerator/amp/cmd/amp/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	padding = 3
)

var (
	listTopicCmd = &cobra.Command{
		Use:     "ls",
		Short:   "List topics",
		Example: "-q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTopic(AMP)
		},
	}
)

func init() {
	TopicCmd.AddCommand(listTopicCmd)
}

func listTopic(amp *cli.AMP) error {
	request := &topic.ListRequest{}

	client := topic.NewTopicClient(amp.Conn)
	reply, err := client.List(context.Background(), request)
	if err != nil {
		mgr.Error(grpc.ErrorDesc(err))
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\t")
	for _, topic := range reply.Topics {
		fmt.Fprintf(w, "%s\t%s\t\n", topic.Id, topic.Name)
	}
	w.Flush()

	return nil
}
