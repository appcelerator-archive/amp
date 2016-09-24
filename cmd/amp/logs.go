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
	Use:   "logs [OPTIONS] SERVICE",
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
	logsCmd.Flags().String("container", "", "Filter by the given container")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	logsCmd.Flags().String("message", "", "Filter the message content by the given pattern")
	logsCmd.Flags().BoolP("meta", "m", false, "Display entry metadata")
	logsCmd.Flags().String("node", "", "Filter by the given node")
	logsCmd.Flags().StringP("number", "n", "100", "Number of results")
	logsCmd.Flags().String("stack", "", "Filter by the given stack")

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
		fmt.Printf("stack: %v\n", cmd.Flag("stack").Value)
		fmt.Printf("message: %v\n", cmd.Flag("message").Value)
		fmt.Printf("container: %v\n", cmd.Flag("container_id").Value)
		fmt.Printf("node: %v\n", cmd.Flag("node_id").Value)
		fmt.Printf("n: %v\n", cmd.Flag("n").Value)
		fmt.Printf("meta: %v\n", cmd.Flag("meta").Value)
	}

	request := logs.GetRequest{}
	if len(args) > 0 {
		request.Service = args[0]
	}
	request.Container = cmd.Flag("container").Value.String()
	request.Node = cmd.Flag("node").Value.String()
	request.Message = cmd.Flag("message").Value.String()
	request.Stack = cmd.Flag("stack").Value.String()
	if request.Size, err = strconv.ParseInt(cmd.Flag("number").Value.String(), 10, 64); err != nil {
		log.Fatalf("Unable to convert n parameter: %v\n", cmd.Flag("n").Value.String())
	}
	var meta bool
	if meta, err = strconv.ParseBool(cmd.Flag("meta").Value.String()); err != nil {
		log.Fatalf("Unable to convert meta parameter: %v\n", cmd.Flag("meta").Value.String())
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
		displayLogEntry(entry, meta)
	}
	if !follow {
		return nil
	}

	// If follow is requested, get subsequent logs and stream it
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
		displayLogEntry(entry, meta)
	}
	return nil
}

func displayLogEntry(entry *logs.LogEntry, meta bool) {
	if meta {
		fmt.Printf("%+v\n", entry)
	} else {
		fmt.Printf("%s\n", entry.Message)
	}
}
