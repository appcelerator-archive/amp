package main

import (
	"context"
	"log"

	plugin "github.com/appcelerator/amp/cluster/plugin/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
)

const (
	templateURL = "https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl"
)

func initSession() {

}

var (
	opts = &plugin.RequestOptions{
		OnFailure: "ROLLBACK",
		Params:    []string{},
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

	log.Println(awsutil.StringValue(resp))

	input := &cf.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	}
	if opts.Sync {
		if err := svc.WaitUntilStackCreateCompleteWithContext(ctx, input); err != nil {
			log.Fatal(err)
		}
		log.Printf("stack created: %s\n", opts.StackName)
	}
}

func update(cmd *cobra.Command, args []string) {
	log.Println("update command not implemented yet")
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

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   provision,
	}
	initCmd.Flags().StringVar(&opts.OnFailure, "onfailure", "ROLLBACK", "action to take if stack creation fails")

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

	rootCmd.AddCommand(initCmd, updateCmd, destroyCmd)

	_ = rootCmd.Execute()
}
