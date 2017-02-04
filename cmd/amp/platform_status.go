package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"sort"
	"text/tabwriter"
)

type statusLine struct {
	stackName string
	status    string
}

// PlatformStatus is the main command for attaching platform subcommands.
var PlatformStatus = &cobra.Command{
	Use:   "status [OPTION...]",
	Short: "Display infrastrucure stacks status",
	Long: `The status command retrieves current status of AMP platform (stopped, partially running, running).
The command returns 1 if status is not running.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAMPStatus(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformStatus)
}

func getAMPStatus(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	statusList := []statusLine{}
	local := false
	if err := amp.Connect(); err != nil {
		if !amp.IsLocalhost() {
			manager.fatalf("Amp server is not available:\n", err)
		}
		if err := manager.connectDocker(); err != nil {
			manager.fatalf("Docker connect error: %v\n", err)
		}
		line := getLocalAmpcoreStatus(ctx, manager)
		statusList = append(statusList, line)
		local = true
	} else {
		sort.Strings(stack.InfraStackList)
		for _, stackName := range stack.InfraStackList {
			request := &stack.StackRequest{Name: stackName}
			client := stack.NewStackServiceClient(amp.Conn)
			reply, err := client.GetStackStatus(ctx, request)
			status := ""
			if err != nil {
				status = fmt.Sprintf("error: %v", err)
			} else {
				status = reply.Answer
			}
			statusList = append(statusList, statusLine{stackName: stackName, status: status})
		}
	}
	displayStatus(manager, statusList, local)
	return nil
}

func getLocalAmpcoreStatus(ctx context.Context, manager *ampManager) statusLine {
	request := &stack.StackRequest{Name: "ampcore"}
	server := stack.NewServer(nil, manager.docker)
	reply, err := server.GetStackStatus(ctx, request)
	status := ""
	if err != nil {
		status = fmt.Sprintf("error: %v", err)
	} else {
		status = reply.Answer
	}
	return statusLine{stackName: "ampcore", status: status}
}

func displayStatus(manager *ampManager, statusList []statusLine, local bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	if local {
		fmt.Fprintf(w, "%s\n", manager.fcolRegular("STACK\tSTATUS (local)"))
	} else {
		fmt.Fprintf(w, "%s\n", manager.fcolRegular("STACK\tSTATUS"))
	}
	for _, line := range statusList {
		status := line.status
		stackName := line.stackName
		if status == "running" {
			fmt.Fprintf(w, "%s\n", manager.fcolSuccess(fmt.Sprintf("%s\t%s", stackName, status)))
		} else if status == "starting" {
			fmt.Fprintf(w, "%s\n", manager.fcolRegular(fmt.Sprintf("%s\t%s", stackName, status)))
		} else if status == "failling" {
			fmt.Fprintf(w, "%s\n", manager.fcolWarn(fmt.Sprintf("%s\t%s", stackName, status)))
		} else if status == "stopped" {
			fmt.Fprintf(w, "%s\n", manager.fcolInfo(fmt.Sprintf("%s\t%s", stackName, status)))
		} else {
			fmt.Fprintf(w, "%s\n", manager.fcolError(fmt.Sprintf("%s\t%s", stackName, status)))
		}
	}
	w.Flush()
}
