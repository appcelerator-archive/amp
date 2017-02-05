package main

import (
	"github.com/appcelerator/amp/api/client"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// PlatformStart is the main command for attaching platform subcommands.
var PlatformStart = &cobra.Command{
	Use:     "deploy [OPTIONS] STACK",
	Aliases: []string{"up", "start"},
	Short:   "Start infrastructure stacks, default ampcore, specific one or all, using 'all'",
	Long:    `The start command starts the infrastructure stacks, ampcore by default, 'all' start all available infrastructure stacks`,
	Run: func(cmd *cobra.Command, args []string) {
		startAMP(AMP, cmd, args)
	},
}

func init() {
	PlatformStart.Flags().BoolP("local", "l", false, "Use local amp image")
	PlatformStart.Flags().Bool("with-registry-auth", false, "Send registry authentication details to Swarm agents")
	PlatformCmd.AddCommand(PlatformStart)
}

func startAMP(amp *client.AMP, cmd *cobra.Command, args []string) error {
	//init manager for colors and docker connection
	manager := newManager(cmd.Flag("verbose").Value.String())
	//init ctx checking authorization
	ctx, err := amp.GetAuthorizedContext()
	if err != nil {
		manager.fatalf("%v\n", err)
	}
	//set registryAuth
	registryAuth := false
	if cmd.Flag("with-registry-auth").Value.String() == "true" {
		registryAuth = true
	}
	//set amptag: local for developement tests
	ampTag := ""
	if cmd.Flag("local").Value.String() == "true" {
		ampTag = "local"
	}
	data := ""
	//get the stackname, default ampcore, "all" to start all the stack
	stackName := "ampcore"
	if len(args) >= 1 {
		stackName = args[0]
	}

	if stackName == "ampcore" || stackName == "all" {
		//ampcore is launch only locally. It can't be launched remotely because amplifier is not yet started
		if !amp.IsLocalhost() {
			//if the target is remote then error
			if stackName == "ampcore" {
				manager.fatalf("ampcore can be started only locally on swarm manager machine\n")
			}
		} else {
			//connection to the local docker engine
			if err := manager.connectDocker(); err != nil {
				manager.fatalf("Needs AMP installed to start ampcore")
			}
			//check prerequisite
			if err := manager.systemPrerequisites(); err != nil {
				manager.fatalf("Prerequiste error: %v\n", err)
			}
			manager.printf(colRegular, "starting stack: ampcore\n")
			//load ampcore compose file replacing first the variable names by their values using amp.manifest file
			data, err = stack.ResolvedComposeFileVariables("ampcore.yml", "amp.manifest", ampTag)
			if err != nil {
				manager.fatalf("start ampCore error: %v\n", err)
			}
			//execute Deploy command using stack code directelly (not going through grpc)
			stackInstance := &stack.Stack{
				Name:     "ampcore",
				FileData: data,
			}
			server := stack.NewServer(nil, manager.docker)
			request := &stack.StackDeployRequest{Stack: stackInstance, RegistryAuth: registryAuth}
			reply, err := server.Deploy(ctx, request)
			if err != nil {
				manager.fatalf("start ampCore error: %v\n", err)
			}
			manager.printf(colSuccess, reply.Answer)
			//if not all then it's finished
			if stackName != "all" {
				return nil
			}
			//wait for haproxy ready, if not, the other stacks won't be able to up because they use amplifier and so need to be route to it
			if manager.verbose {
				manager.printf(colRegular, "Waiting for amplifier and haproxy available\n")
			}
			if err := server.WaitForServiceReady(ctx, "ampcore_haproxy", 120); err != nil {
				manager.printf(colWarn, "haproxy starting timeout\n")
				return nil
			}
			if manager.verbose {
				manager.printf(colSuccess, "haproxy is ready, remove command are availables\n")
			}
		}
	}
	//connect to the remote amplifier (routed by haproxy), this amplifier can be a local service for dev env., it it's used as it was remote.
	if err := amp.Connect(); err != nil {
		manager.fatalf("Server not yet ready\n")
	}
	//start all the stacks
	if stackName == "all" {
		for _, name := range stack.InfraShortStackList {
			startStack(ctx, amp, manager, name, ampTag)
		}
		return nil
	}
	//start only the asked stack
	return startStack(ctx, amp, manager, stackName, ampTag)
}

//start one stack
func startStack(ctx context.Context, amp *client.AMP, manager *ampManager, stackName string, ampTag string) error {
	manager.printf(colRegular, "starting stack: %s\n", stackName)
	//start a stack using stack code through grpc.
	stackInstance := &stack.Stack{
		Name:   stackName,
		AmpTag: ampTag,
	}
	retry := 0
	request := &stack.StackDeployRequest{Stack: stackInstance, RegistryAuth: registryAuth}
	for {
		client := stack.NewStackServiceClient(amp.Conn)
		reply, err := client.Deploy(ctx, request)
		if err != nil {
			manager.printf(colWarn, "start %s error: %v retry=%d\n", stackName, err, retry)
			retry++
			if retry > 3 {
				manager.printf(colError, "start %s error: %v retry=%d\n", stackName, err)
				return err
			}
		}
		manager.printf(colSuccess, reply.Answer)
		return nil
	}
}
