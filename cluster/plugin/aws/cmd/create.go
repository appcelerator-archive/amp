package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewCreateCommand returns a new instance of the create command
func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   create,
	}
	cmd.Flags().StringVar(&plugin.Config.OnFailure, "onfailure", "ROLLBACK", "action to take if stack creation fails")
	cmd.Flags().StringVarP(&plugin.Config.TemplateURL, "template", "t", plugin.DefaultTemplateURL, "cloud formation template url")
	return cmd
}

func create(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.AWS.CreateStack(ctx, 20)
	if err != nil {
		j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
		if jerr != nil {
			log.Println(jerr.Error())
			log.Fatal(err)
		}
		fmt.Println(j)
		os.Exit(1)
	}

	if plugin.AWS.Config.Sync {
		if err := plugin.AWS.ScanEvents(plugin.StackModeCreate); err != nil {
			j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
			if jerr != nil {
				log.Println(jerr.Error())
				log.Fatal(err)
			}
			fmt.Println(j)
			os.Exit(1)
		}
		Info()
	} else {
		event := plugin.StackEvent{
			EventId:              "NoSync-000",
			LogicalResourceId:    plugin.AWS.Config.StackName,
			ResourceType:         "AWS::CloudFormation:Stack",
			ResourceStatus:       cloudformation.StackStatusCreateInProgress,
			ResourceStatusReason: "The sync flag was not used, please check the status of the stack on the AWS console and read the output",
			Timestamp:            time.Now().Format(time.UnixDate),
			StackId:              *resp.StackId,
		}
		j, err := plugin.PluginOutputToJSON(&event, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(j)
	}
}
