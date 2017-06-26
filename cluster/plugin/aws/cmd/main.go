package main

import (
	"context"
	"fmt"
	"log"

	plugin "github.com/appcelerator/amp/cluster/plugin/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
)

var (
	opts = &plugin.RequestOptions{
		OnFailure: "DO_NOTHING",
		Params:    []string{},
		TemplateURL: plugin.DefaultTemplateURL,
	}

	sess *session.Session
	svc *cf.CloudFormation
)

func initClient(cmd *cobra.Command, args []string) {
	sess = session.Must(session.NewSession())

	// Create the service's client with the session.
	svc = cf.New(sess, aws.NewConfig().WithRegion(opts.Region).WithLogLevel(aws.LogOff))
}

func provision(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	resp, err := plugin.CreateStack(ctx, svc, opts, 20)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Sync {
		input := &cf.DescribeStacksInput{
			StackName: aws.String(opts.StackName),
		}
		if err := svc.WaitUntilStackCreateCompleteWithContext(ctx, input); err != nil {
			log.Fatal(err)
		}
		// use the info command to print json cluster info to stdout
		info(cmd, args)
	} else {
		// only print to stdout if not sync; otherwise stdout is used to display json stack output information now
		log.Printf("stack created: %s\n", opts.StackName)
		log.Println(awsutil.StringValue(resp.StackId))
	}
}

func update(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.UpdateStack(ctx, svc, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))

	input := &cf.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	}
	if opts.Sync {
		if err := svc.WaitUntilStackUpdateCompleteWithContext(ctx, input); err != nil {
			log.Fatal(err)
		}
		log.Printf("stack updated: %s\n", opts.StackName)
	}
}

func destroy(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.DeleteStack(ctx, svc, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))

	input := &cf.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	}
	if opts.Sync {
		if err := svc.WaitUntilStackDeleteCompleteWithContext(ctx, input); err != nil {
			log.Fatal(err)
		}
		log.Printf("stack deleted: %s\n", opts.StackName)
	}
}

func info(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	resp, err := plugin.InfoStack(ctx, svc, opts)
	if err != nil {
		log.Fatal(err)
	}

	j, err := plugin.StackOutputToJSON(resp)
	if err != nil {
		log.Fatal(err)
	}

	// print json result to stdout
	fmt.Print(j)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "awsplugin",
		Short: "init/update/destroy an AWS cluster in Docker swarm mode",
		PersistentPreRun: initClient,
	}
	rootCmd.PersistentFlags().StringVarP(&opts.Region, "region", "r", "", "aws region")
	rootCmd.PersistentFlags().StringVarP(&opts.StackName, "stackname", "n", "", "aws stack name")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.Params, "parameter", "p", []string{}, "parameter")
	rootCmd.PersistentFlags().BoolVarP(&opts.Sync, "sync", "s", false, "block until operation is complete")
	rootCmd.PersistentFlags().StringVarP(&opts.TemplateURL, "template", "t", plugin.DefaultTemplateURL, "cloud formation template url")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   provision,
	}
	initCmd.Flags().StringVar(&opts.OnFailure, "onfailure", "ROLLBACK", "action to take if stack creation fails")

	infoCmd := &cobra.Command{
		Use: "info",
		Short: "get information about the cluster",
		Run: info,
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update the cluster",
		Run:   update,
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the cluster",
		Run:   destroy,
	}

	rootCmd.AddCommand(initCmd, infoCmd, updateCmd, destroyCmd)

	_ = rootCmd.Execute()
}
