package main

import (
	"context"
	"fmt"

	sk "github.com/appcelerator/amp/cluster/agent/swarm"
	"github.com/docker/swarmkit/api"
	"github.com/spf13/cobra"
)

func NewMonitorCommand() *cobra.Command {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor swarm events",
		RunE:   monitor,
	}
	return monitorCmd
}

func monitor(cmd *cobra.Command, args []string) error {
	c, conn, err := sk.Dial(sk.DefaultSocket())
	if err != nil {
		return err
	}

	// this is just to prove things are working...
	nodes, err := sk.ListNodes(sk.Context(), c, sk.LiveNodesFilter)
	for _, n := range nodes {
		fmt.Println("Node ID: ", n.GetID())
	}

	watcher := api.NewWatchClient(conn)
	watchEntry := sk.NewWatchRequestEntry("service", sk.WatchActionKindAll, nil)
	watchEntries := []*api.WatchRequest_WatchEntry{
		watchEntry,
	}

	// will probably need a cancelable context
	ctx := context.TODO()
	in := sk.NewWatchRequest(watchEntries, nil, true)
	w, err := watcher.Watch(ctx, in)
	if err != nil {
		return err
	}

	for {
		msg, err := w.Recv()
		if err != nil {
			return err
		}

		fmt.Println(msg.String())
	}

	return nil
}
