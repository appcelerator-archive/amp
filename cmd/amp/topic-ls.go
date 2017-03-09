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
		Long:    `The list command returns all available topics.`,
		Example: "amp topic ls -q",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTopic(AMP)
		},
	}
)

func init() {
	TopicCmd.AddCommand(listTopicCmd)
}

func listTopic(amp *cli.AMP) (err error) {
	request := &topic.ListRequest{}

	client := topic.NewTopicClient(amp.Conn)
	reply, er := client.List(context.Background(), request)
	if er != nil {
		manager.fatalf(grpc.ErrorDesc(er))
		return
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
