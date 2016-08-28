package main

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch the logs",
	Long:  `Search through all the logs of the system and fetch entries matching provided criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Logs(AMP, cmd, args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	// TODO logsCmd.Flags().String("timestamp", "", "filter by the given timestamp")
	logsCmd.Flags().String("service_id", "", "Filter by the given service id")
	logsCmd.Flags().String("service_name", "", "Filter by the given service name")
	logsCmd.Flags().String("message", "", "Filter the message content by the given pattern")
	logsCmd.Flags().String("container_id", "", "Filter by the given container id")
	logsCmd.Flags().String("node_id", "", "Filter by the given node id")
	logsCmd.Flags().String("from", "-1", "Fetch from the given index")
	logsCmd.Flags().StringP("number", "n", "100", "Number of results")
	logsCmd.Flags().BoolP("short", "s", false, "Display message content only")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")

	RootCmd.AddCommand(logsCmd)
}

// Logs fetches the logs
func Logs(amp *client.AMP, cmd *cobra.Command, args []string) error {
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		return err
	}
	if amp.Verbose() {
		fmt.Println("Logs")
		fmt.Printf("service_id: %v\n", cmd.Flag("service_id").Value)
		fmt.Printf("service_name: %v\n", cmd.Flag("service_name").Value)
		fmt.Printf("message: %v\n", cmd.Flag("message").Value)
		fmt.Printf("container_id: %v\n", cmd.Flag("container_id").Value)
		fmt.Printf("node_id: %v\n", cmd.Flag("node_id").Value)
		fmt.Printf("from: %v\n", cmd.Flag("from").Value)
		fmt.Printf("n: %v\n", cmd.Flag("n").Value)
		fmt.Printf("short: %v\n", cmd.Flag("short").Value)
	}

	request := logs.GetRequest{}
	request.ServiceId = cmd.Flag("service_id").Value.String()
	request.ServiceName = cmd.Flag("service_name").Value.String()
	request.Message = cmd.Flag("message").Value.String()
	request.ContainerId = cmd.Flag("container_id").Value.String()
	request.NodeId = cmd.Flag("node_id").Value.String()
	if request.From, err = strconv.ParseInt(cmd.Flag("from").Value.String(), 10, 64); err != nil {
		log.Fatalf("Unable to convert from parameter: %v\n", cmd.Flag("from").Value.String())
	}
	if request.Size, err = strconv.ParseInt(cmd.Flag("number").Value.String(), 10, 64); err != nil {
		log.Fatalf("Unable to convert n parameter: %v\n", cmd.Flag("n").Value.String())
	}
	var short bool
	if short, err = strconv.ParseBool(cmd.Flag("short").Value.String()); err != nil {
		log.Fatalf("Unable to convert short parameter: %v\n", cmd.Flag("short").Value.String())
	}
	var follow bool
	if follow, err = strconv.ParseBool(cmd.Flag("follow").Value.String()); err != nil {
		log.Fatalf("Unable to convert f parameter: %v\n", cmd.Flag("f").Value.String())
	}

	// Get logs from elasticsearch
	c := logs.NewLogsClient(amp.Conn)
	r, err := c.Get(ctx, &request)
	if err != nil {
		return err
	}
	for _, entry := range r.Entries {
		displayLogEntry(entry, short)
	}
	if !follow {
		return nil
	}

	// If follow is requested, get subsequent logs from Kafka and stream it
	stream, err := c.GetStream(ctx, &request)
	if err != nil {
		return err
	}
	for {
		entry, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		displayLogEntry(entry, short)
	}
	return nil
}

func displayLogEntry(entry *logs.LogEntry, short bool) {
	if short {
		fmt.Printf("%s\n", entry.Message)
	} else {
		fmt.Printf("%+v\n", entry)
	}
}
