package main

import (
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

var (
	stackSpec = &plugin.StackSpec{
		OnFailure: "ROLLBACK",
		Params: []string{},
	}
)

func provision(cmd *cobra.Command, args []string) {
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc := cf.New(sess,
		aws.NewConfig().WithRegion(stackSpec.Region).WithLogLevel(aws.LogOff))

	resp, err := plugin.CreateStack(svc, stackSpec, 20)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))
}

func update(cmd *cobra.Command, args []string) {
	log.Println("update command not implemented yet")
}

func destroy(cmd *cobra.Command, args []string) {
	log.Println("destroy command not implemented yet")
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "awsplugin",
		Short: "init/update/destroy an AWS cluster in Docker swarm mode",
	}
	rootCmd.PersistentFlags().StringVarP(&stackSpec.Region, "region", "r", "", "aws region")
	rootCmd.PersistentFlags().StringVarP(&stackSpec.StackName, "stackname", "n", "", "aws stack name")
	rootCmd.PersistentFlags().StringSliceVarP(&stackSpec.Params, "parameter", "p", []string{}, "parameter")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   provision,
	}
	initCmd.Flags().StringVar(&stackSpec.OnFailure, "onfailure", "ROLLBACK", "action to take if stack creation fails")

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

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
