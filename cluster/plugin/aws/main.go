package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
)

const (
	templateURL = "https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl"
)

type StackSpec struct {
	stackName string
	region string
}

var (
	stackSpec = &StackSpec{}
)

func createStack(svc *cf.CloudFormation, params []*cf.Parameter, stackName string, timeout int64) {
	input := &cf.CreateStackInput{
		StackName: aws.String(stackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String("DELETE"),
		Parameters:       params,
		TemplateURL:      aws.String(templateURL),
		TimeoutInMinutes: aws.Int64(timeout),
	}

	resp, err := svc.CreateStack(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))
}

func provision(cmd *cobra.Command, args []string) {
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc := cf.New(sess,
		aws.NewConfig().WithRegion(stackSpec.region).WithLogLevel(aws.LogOff))

	params := []*cf.Parameter{
		{ParameterKey: aws.String("KeyName"), ParameterValue: aws.String("tony-amp-dev") },
	}

	createStack(svc, params, stackSpec.stackName, 20)

	log.Println(svc.APIVersion)
}

func update(cmd *cobra.Command, args []string) {
	log.Println("update command not implemented yet")
}

func destroy(cmd *cobra.Command, args []string) {
	log.Println("destroy command not implemented yet")
}

func main() {
	rootCmd := &cobra.Command{
		Use: "awsplugin",
		Short: "init/update/destroy an AWS cluster in Docker swarm mode",
	}
	rootCmd.PersistentFlags().StringVarP(&stackSpec.stackName, "name", "n", "", "stack name")
	rootCmd.PersistentFlags().StringVarP(&stackSpec.region, "region", "r", "", "aws region")

	initCmd := &cobra.Command{
		Use: "init",
		Short: "init cluster in swarm mode",
		Run: provision,
	}

	updateCmd := &cobra.Command{
		Use: "update",
		Short: "update the cluster",
		Run: update,
	}

	destroyCmd := &cobra.Command{
		Use: "destroy",
		Short: "destroy the cluster",
		Run: destroy,
	}

	rootCmd.AddCommand(initCmd, updateCmd, destroyCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
