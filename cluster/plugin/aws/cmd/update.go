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

// NewUpdateCommand returns a new instance of the update command
func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update the cluster",
		Run:   update,
	}
	cmd.Flags().StringVarP(&plugin.Config.TemplateURL, "template", "t", plugin.DefaultTemplateURL, "cloud formation template url")
	return cmd
}

func update(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.AWS.UpdateStack(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if plugin.AWS.Config.Sync {
		err := plugin.AWS.ScanEvents(plugin.StackModeUpdate)
		if err != nil {
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
