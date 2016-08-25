package cli

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strconv"
)

// Logs fetches the logs
func Logs(amp *client.AMP, cmd *cobra.Command) error {
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
	if request.Size, err = strconv.ParseInt(cmd.Flag("n").Value.String(), 10, 64); err != nil {
		log.Fatalf("Unable to convert n parameter: %v\n", cmd.Flag("n").Value.String())
	}
	var short bool
	if short, err = strconv.ParseBool(cmd.Flag("short").Value.String()); err != nil {
		log.Fatalf("Unable to convert short parameter: %v\n", cmd.Flag("short").Value.String())
	}
	var follow bool
	if follow, err = strconv.ParseBool(cmd.Flag("f").Value.String()); err != nil {
		log.Fatalf("Unable to convert f parameter: %v\n", cmd.Flag("f").Value.String())
	}

	// Get logs from elasticsearch
	c := logs.NewLogsClient(amp.Connect())
	defer amp.Disconnect()
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
