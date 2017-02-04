package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

// StackList command to list stack
var StackList = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List the stacks",
	Long:    `List the running stacks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stackList(AMP, cmd, args)
	},
}

func init() {
	StackCmd.AddCommand(StackList)
}

func stackList(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.Connect(); err != nil {
		manager.fatalf("Amp server is not available\n")
	}

	request := &stack.StackRequest{}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.ListStacks(ctx, request)
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	lines := strings.Split(reply.Answer, "\n")
	manager.printf(colRegular, "%s\n", lines[0])
	if len(lines) > 1 {
		lines = lines[1:]
		sort.Strings(lines)
		for _, line := range lines {
			if line != "" {
				manager.printf(colSuccess, "%s\n", line)
			}
		}
	}
	return nil
}
