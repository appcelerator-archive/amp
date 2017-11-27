package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/appcelerator/amp/cluster/plugin/aws/plugin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
)

var (
	Version string
	Build   string
	opts    = &plugin.RequestOptions{
		OnFailure:       "DO_NOTHING",
		Params:          []string{},
		TemplateURL:     plugin.DefaultTemplateURL,
		AccessKeyId:     "",
		SecretAccessKey: "",
		Profile:         "default",
	}

	sess *session.Session
	svc  *cf.CloudFormation
)

const (
	STACK_MODE_CREATE = "create"
	STACK_MODE_DELETE = "delete"
	STACK_MODE_UPDATE = "update"
)

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("Version: %s - Build: %s\n", Version, Build)
}

func initClient(cmd *cobra.Command, args []string) {
	// export vars if creds are passed as arguments
	if opts.AccessKeyId != "" && opts.SecretAccessKey != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", opts.AccessKeyId)
		os.Setenv("AWS_SECRET_ACCESS_KEY", opts.SecretAccessKey)
	} else {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		os.Setenv("AWS_PROFILE", opts.Profile)
	}
	config := aws.NewConfig().WithLogLevel(aws.LogOff)
	// region can be set with a CLI option, but if not set it can be set by the config file
	if opts.Region != "" {
		config.Region = aws.String(opts.Region)
	}
	sess = session.Must(session.NewSession())

	// Create the service's client with the session.
	svc = cf.New(sess, config)
}

// used by the create function to parse the events and send meaningful information to the CLI
func scanEvents(mode string) error {
	var emptyReason string
	significantEventTypes := map[string]bool{
		"AWS::CloudFormation::Stack":         true,
		"AWS::EC2::VPC":                      true,
		"AWS::CloudFormation::WaitCondition": true,
		"AWS::AutoScaling::AutoScalingGroup": true,
		"AWS::EFS::FileSystem":               true,
	}
	significantEventStatuses := map[string]bool{
		cf.StackStatusCreateInProgress:   false,
		cf.StackStatusCreateComplete:     true,
		cf.StackStatusCreateFailed:       true,
		cf.StackStatusRollbackInProgress: true,
		cf.StackStatusRollbackComplete:   true,
		cf.StackStatusDeleteInProgress:   true,
	}
	eventInput := &cf.DescribeStackEventsInput{
		StackName: aws.String(opts.StackName),
		NextToken: nil,
	}
	eventIds := map[string]bool{}
	for {
		resp, err := svc.DescribeStackEvents(eventInput)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling: Rate exceeded") {
				// ignore it, and continue processing the events
				time.Sleep(time.Second)
				continue
			} else if mode == STACK_MODE_DELETE && strings.Contains(err.Error(), fmt.Sprintf("Stack [%s] does not exist", opts.StackName)) {
				// stack does not exist, probably because we've just successfuly deleted it
				event := plugin.StackEvent{
					EventId:           "NoSync-999",
					LogicalResourceId: opts.StackName,
					ResourceType:      "AWS::CloudFormation:Stack",
					ResourceStatus:    cf.StackStatusDeleteComplete,
					Timestamp:         time.Now().Format(time.UnixDate),
				}
				j, err := plugin.PluginOutputToJSON(&event, nil, nil)
				if err != nil {
					return err
				}
				fmt.Println(j)
				return nil
			}
			return err
		}
		events := resp.StackEvents
		sort.Slice(events, func(i, j int) bool { return events[i].Timestamp.Unix() < events[j].Timestamp.Unix() })
		for _, se := range events {
			if eventIds[*se.EventId] {
				// already processed, ignore it
				continue
			}
			eventIds[*se.EventId] = true

			// filtering: stack events and creation / deletion of a few types
			if *se.ResourceType == "AWS::CloudFormation::Stack" ||
				(significantEventTypes[*se.ResourceType] == true && significantEventStatuses[*se.ResourceStatus] == true) {
				if se.ResourceStatusReason == nil {
					// we'll have an indirection, so better secure it
					se.ResourceStatusReason = &emptyReason
				}
				eventOutput := plugin.StackEvent{
					EventId:              *se.EventId,
					LogicalResourceId:    *se.LogicalResourceId,
					ResourceStatus:       *se.ResourceStatus,
					ResourceStatusReason: *se.ResourceStatusReason,
					ResourceType:         *se.ResourceType,
					Timestamp:            se.Timestamp.Format(time.UnixDate),
				}
				j, err := plugin.PluginOutputToJSON(&eventOutput, nil, nil)
				if err != nil {
					return err
				}
				fmt.Println(j)
			}
			// check end condition, based on the action
			if *se.ResourceType != "AWS::CloudFormation::Stack" {
				continue
			}
			switch mode {
			case STACK_MODE_CREATE:
				switch *se.ResourceStatus {
				case cf.StackStatusCreateComplete, cf.StackStatusDeleteComplete, cf.StackStatusRollbackComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			case STACK_MODE_DELETE:
				switch *se.ResourceStatus {
				case cf.StackStatusDeleteFailed, cf.StackStatusDeleteComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			case STACK_MODE_UPDATE:
				switch *se.ResourceStatus {
				case cf.ResourceStatusUpdateFailed, cf.ResourceStatusUpdateComplete, cf.StackStatusRollbackComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			default:
				return fmt.Errorf("unknown mode: %s", mode)
			}
		}
		// waiting before reading next batch of events (if too short, the API throttles)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func create(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.CreateStack(ctx, svc, opts, 20)
	if err != nil {
		j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
		if jerr != nil {
			log.Println(jerr.Error())
			log.Fatal(err)
		}
		fmt.Println(j)
		os.Exit(1)
	}

	if opts.Sync {
		err := scanEvents(STACK_MODE_CREATE)
		if err != nil {
			j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
			if jerr != nil {
				log.Println(jerr.Error())
				log.Fatal(err)
			}
			fmt.Println(j)
			os.Exit(1)
		}
		info(cmd, args)
	} else {
		event := plugin.StackEvent{
			EventId:              "NoSync-000",
			LogicalResourceId:    opts.StackName,
			ResourceType:         "AWS::CloudFormation:Stack",
			ResourceStatus:       cf.StackStatusCreateInProgress,
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

func update(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.UpdateStack(ctx, svc, opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Sync {
		err := scanEvents(STACK_MODE_UPDATE)
		if err != nil {
			j, jerr := plugin.PluginOutputToJSON(nil, nil, err)
			if jerr != nil {
				log.Println(jerr.Error())
				log.Fatal(err)
			}
			fmt.Println(j)
			os.Exit(1)
		}
		info(cmd, args)
	} else {
		event := plugin.StackEvent{
			EventId:              "NoSync-000",
			LogicalResourceId:    opts.StackName,
			ResourceType:         "AWS::CloudFormation:Stack",
			ResourceStatus:       cf.StackStatusUpdateInProgress,
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

func delete(cmd *cobra.Command, args []string) {
	var flag bool
	ctx := context.Background()
	resp, err := plugin.ListStack(ctx, svc)
	if err != nil {
		log.Fatal(err)
	}

	for _, stk := range resp.StackSummaries {
		if *stk.StackName != opts.StackName {
			flag = true
			continue
		} else {
			switch *stk.StackStatus {
			case cf.StackStatusCreateInProgress, cf.StackStatusCreateComplete, cf.StackStatusCreateFailed, cf.StackStatusRollbackFailed, cf.StackStatusRollbackComplete:
				if _, err := plugin.DeleteStack(ctx, svc, opts); err != nil {
					log.Fatal(err)
				}

				if opts.Sync {
					err := scanEvents(STACK_MODE_DELETE)
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
			case cf.StackStatusDeleteInProgress:
				log.Fatal("cluster deletion already in progress, check again in a few minutes")
			case cf.StackStatusRollbackInProgress:
				log.Fatal("cluster deployment is performing a rollback, the deletion is not possible right now, please try again in in a few minutes")
			default:
				log.Fatalf("cluster deletion not possible with the current status [%s], try again in a few minutes", *stk.StackStatus)
			}
		}
	}
	if flag || len(resp.StackSummaries) == 0 {
		if j, err := plugin.PluginOutputToJSON(nil, nil, fmt.Errorf("stack [%s] doesn't seem to exist", opts.StackName)); err == nil {
			fmt.Println(j)
			os.Exit(1)
		}
		if opts.Region != "" {
			log.Fatal(fmt.Sprintf("stack [%s] doesn't seem to exist in region %s", opts.StackName, opts.Region))
		} else {
			log.Fatal(fmt.Sprintf("stack [%s] doesn't seem to exist", opts.StackName))
		}
	}
}

func info(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.InfoStack(ctx, svc, opts)
	if err != nil {
		if j, jerr := plugin.PluginOutputToJSON(nil, nil, err); jerr == nil {
			// print json error to stdout
			fmt.Println(j)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}

	j, err := plugin.PluginOutputToJSON(nil, resp, nil)
	if err != nil {
		log.Fatal(err)
	}

	// print json result to stdout
	fmt.Println(j)
}

func deprecationWarning(cmd *cobra.Command, args []string) {
	fmt.Println("Deprecated, update your CLI")
}

func main() {
	rootCmd := &cobra.Command{
		Use:              "awsplugin",
		Short:            "init/update/destroy an AWS cluster in Docker swarm mode",
		PersistentPreRun: initClient,
	}
	rootCmd.PersistentFlags().StringVarP(&opts.Region, "region", "r", "", "aws region")
	rootCmd.PersistentFlags().StringVarP(&opts.StackName, "stackname", "n", "", "aws stack name")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.Params, "parameter", "p", []string{}, "parameter")
	rootCmd.PersistentFlags().BoolVarP(&opts.Sync, "sync", "s", true, "block until operation is complete")
	rootCmd.PersistentFlags().StringVar(&opts.AccessKeyId, "access-key-id", "", "access key id (for example, AKIAIOSFODNN7EXAMPLE)")
	rootCmd.PersistentFlags().StringVar(&opts.SecretAccessKey, "secret-access-key", "", "secret access key (for example, wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY)")
	rootCmd.PersistentFlags().StringVar(&opts.Profile, "profile", "default", "credential profile")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   create,
	}
	initCmd.Flags().StringVar(&opts.OnFailure, "onfailure", "ROLLBACK", "action to take if stack creation fails")
	initCmd.Flags().StringVarP(&opts.TemplateURL, "template", "t", plugin.DefaultTemplateURL, "cloud formation template url")

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "get information about the cluster (deprecated)",
		Run:   deprecationWarning,
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update the cluster",
		Run:   update,
	}
	updateCmd.Flags().StringVarP(&opts.TemplateURL, "template", "t", plugin.DefaultTemplateURL, "cloud formation template url")

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the cluster",
		Run:   delete,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "version of the plugin",
		Run:   version,
	}
	rootCmd.AddCommand(versionCmd, initCmd, infoCmd, updateCmd, destroyCmd)

	_ = rootCmd.Execute()
}
