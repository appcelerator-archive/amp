package main

import (
	"fmt"
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

// PlatformMonitor is the main command for attaching platform subcommands.
var PlatformMonitor = &cobra.Command{
	Use:   "monitor",
	Short: "Display infrastructure tasks services",
	Long:  `The monitor command displays information about AMP infrastructure stacks services.`,
	Run: func(cmd *cobra.Command, args []string) {
		stackMonitor(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformMonitor)
}

func stackMonitor(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager("true")
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if amp.IsLocalhost() {
		if err := manager.connectDocker(); err != nil {
			manager.fatalf("Docker connect error:\n", err)
		}
	}
	connected := false
	lines := []*stack.MonitorService{}
	for {
		// Monitor work locally or remotelly depending on ampcore started or starting
		if !connected {
			if err := amp.Connect(); err != nil {
				//if amp connect arreur then ampcore is stopped or starting
				if !amp.IsLocalhost() {
					//if the server is remote and ampcore not started then no way to monitor
					manager.fatalf("Amp server is not available\n")
				}
			} else {
				connected = true
			}
		}
		if !connected {
			//ampcore is stopped or starting and server address is local, we use local docker to get monitor information
			server := stack.NewServer(nil, manager.docker)
			ret, err := server.GetMonitorLines(ctx)
			if err != nil {
				manager.fatalf("monitor error: %v\n", err)
			}
			lines = ret
		} else {
			//server (amplifier and haproxy) is started locally or remote, we use grpc connection to get monitor information
			request := &stack.StackRequest{}
			client := stack.NewStackServiceClient(amp.Conn)
			reply, err := client.Monitor(ctx, request)
			if err != nil {
				connected = false
			} else {
				lines = reply.Lines
			}
		}
		monitorDisplay(manager, lines, connected)
		time.Sleep(1 * time.Second)
	}
}

// display monitor information using []string
func monitorDisplay(manager *ampManager, lines []*stack.MonitorService, connected bool) {
	fmt.Println("\033[0;0H")
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	if connected {
		fmt.Fprintf(w, "%s\n", manager.fcolRegular("ID\tSERVICE\tSTATUS\tMODE\tREPLICAS\tTASK FAILED"))
	} else {
		fmt.Fprintf(w, "%s\n", manager.fcolRegular("ID\tSERVICE (local)\tSTATUS\tMODE\tREPLICAS\tTASK FAILED"))
	}
	stackName := ""
	for _, serv := range lines {
		if serv.Stack == "" {
			serv.Stack = "free"
		}
		if stackName != serv.Stack {
			line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", serv.Stack+":", "", "", "", "", "")
			fmt.Fprintf(w, "%s\n", manager.fcolUser(line))
			stackName = serv.Stack
		}
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", serv.Id, serv.Service, serv.Status, serv.Mode, serv.Replicas, serv.FailedTask)
		if serv.Status == "running" {
			fmt.Fprintf(w, "%s\n", manager.fcolSuccess(line))
		} else if serv.Status == "failing" {
			fmt.Fprintf(w, "%s\n", manager.fcolWarn(line))
		} else if serv.Status == "stopped" {
			fmt.Fprintf(w, "%s\n", manager.fcolInfo(line))
		} else if serv.Status == "starting" {
			fmt.Fprintf(w, "%s\n", manager.fcolRegular(line))
		} else {
			fmt.Fprintf(w, "%s              \n", manager.fcolInfo(line))
		}
	}
	fmt.Println("\033[2J\033[0;0H")
	w.Flush()
}
