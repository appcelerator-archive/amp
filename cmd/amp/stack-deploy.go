package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
)

// StackDeploy command to start stack
var StackDeploy = &cobra.Command{
	Use:     "deploy [OPTIONS] STACK",
	Aliases: []string{"up", "start"},
	Short:   "Deploy a new stack or update an existing stack",
	Long:    `Start a stack, create all stacks services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stackDeploy(AMP, cmd, args)
	},
}

func init() {
	StackDeploy.Flags().StringP("compose-file", "c", "", "Path to compose file")
	StackDeploy.Flags().String("variables-file", "", "Path to variables file")
	StackDeploy.Flags().Bool("with-registry-auth", false, "Send registry authentication details to Swarm agents")
	StackDeploy.Flags().Bool("debug", false, "Display the the final compose file with resolved variables")
	StackCmd.AddCommand(StackDeploy)
}

func stackDeploy(amp *client.AMP, cmd *cobra.Command, args []string) error {
	manager := newManager(cmd.Flag("verbose").Value.String())
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	if err := amp.Connect(); err != nil {
		manager.fatalf("Amp server is not available\n")
	}
	if len(args) < 1 {
		manager.fatalf("Need stack name as first argument")
	}

	stackName := args[0]
	stackFile := cmd.Flag("compose-file").Value.String()
	varFile := cmd.Flag("variables-file").Value.String()
	if stackFile == "" {
		manager.fatalf("Need compose file path using -c option")
	}
	registryAuth := false
	if cmd.Flag("with-registry-auth").Value.String() == "true" {
		registryAuth = true
	}
	debug := false
	if cmd.Flag("debug").Value.String() == "true" {
		debug = true
	}
	data, erru := stack.ResolvedComposeFileVariables(stackFile, varFile, "")
	if erru != nil {
		return erru
	}
	if debug {
		manager.printf(colRegular, "-------------------------------------------------------------------")
		manager.printf(colRegular, "%s\n", string(data))
		manager.printf(colRegular, "-------------------------------------------------------------------")
	}
	stackInstance := &stack.Stack{
		Name:     stackName,
		FileData: data,
	}
	request := &stack.StackDeployRequest{Stack: stackInstance, RegistryAuth: registryAuth}
	client := stack.NewStackServiceClient(amp.Conn)
	reply, err := client.Deploy(ctx, request)
	if err != nil {
		manager.fatalf("Error: %v\n", err)
		return err
	}
	manager.printf(colSuccess, "%s\n", reply.Answer)
	return nil
}
