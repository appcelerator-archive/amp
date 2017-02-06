package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// PlatformUpdateFile is the main command for attaching platform subcommands.
var PlatformUpdateFile = &cobra.Command{
	Use:   "updateFile [STACK] [OPTION]",
	Short: "Update on the server a compose stack file or a variables file",
	Long:  `Send new stack compose file or a new variable file to server`,
	Run: func(cmd *cobra.Command, args []string) {
		updateStackFile(AMP, cmd, args)
	},
}

func init() {
	PlatformCmd.AddCommand(PlatformUpdateFile)
	PlatformUpdateFile.Flags().StringP("compose-file", "c", "", "Path to new compose file or new variables file")
}

func updateStackFile(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	stackFile := cmd.Flag("compose-file").Value.String()
	if stackFile == "" {
		manager.fatalf("file option is needed -c or --compose-file\n")
	}
	if len(args) == 0 {
		manager.fatalf("First argument should be an infrastructure stack name or 'variable' to send a variables file path\n")
	}
	fileName := args[0]

	if err := amp.Connect(); err != nil {
		manager.fatalf("Server not yet ready\n")
	}
	data, err := ioutil.ReadFile(stackFile)
	if err != nil {
		manager.fatalf("Error reading file: %v\n", err)
	}
	request := &stack.StackUpdateRequest{FileName: fileName, FileData: string(data)}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.UpdateFile(ctx, request)
	if err != nil {
		manager.fatalf("update error: %v\n", err)
	}
	manager.printf(colSuccess, "%s\n", reply.Answer)
	return nil
}
