package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewRemoveCommand returns a new instance of the remove command
func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the cluster",
		Run:   remove,
	}
	return cmd
}

func remove(cmd *cobra.Command, args []string) {
	var flag bool
	ctx := context.Background()
	resp, err := plugin.AWS.ListStack(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, stk := range resp.StackSummaries {
		if *stk.StackName != plugin.AWS.Config.StackName {
			flag = true
			continue
		} else {
			switch *stk.StackStatus {
			case cloudformation.StackStatusCreateInProgress, cloudformation.StackStatusCreateComplete, cloudformation.StackStatusCreateFailed, cloudformation.StackStatusRollbackFailed, cloudformation.StackStatusRollbackComplete:
				if _, err := plugin.AWS.DeleteStack(ctx); err != nil {
					log.Fatal(err)
				}

				if plugin.AWS.Config.Sync {
					err := plugin.AWS.ScanEvents(plugin.StackModeDelete)
					if err != nil {
						j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
						if jerr != nil {
							log.Println(jerr.Error())
							log.Fatal(err)
						}
						fmt.Println(j)
						os.Exit(1)
					}
				}
				return
			case cloudformation.StackStatusDeleteInProgress:
				log.Fatal("cluster deletion already in progress, check again in a few minutes")
			case cloudformation.StackStatusRollbackInProgress:
				log.Fatal("cluster deployment is performing a rollback, the deletion is not possible right now, please try again in in a few minutes")
			default:
				log.Fatalf("cluster deletion not possible with the current status [%s], try again in a few minutes", *stk.StackStatus)
			}
		}
	}
	if flag || len(resp.StackSummaries) == 0 {
		if j, err := plugin.PluginOutputToJSON(nil, nil, fmt.Errorf("stack [%s] doesn't seem to exist", plugin.AWS.Config.StackName)); err == nil {
			fmt.Println(j)
			os.Exit(1)
		}
		if plugin.AWS.Config.Region != "" {
			log.Fatal(fmt.Sprintf("stack [%s] doesn't seem to exist in region %s", plugin.AWS.Config.StackName, plugin.AWS.Config.Region))
		} else {
			log.Fatal(fmt.Sprintf("stack [%s] doesn't seem to exist", plugin.AWS.Config.StackName))
		}
	}
}
