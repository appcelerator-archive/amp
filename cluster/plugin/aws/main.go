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

// StackSpec stores raw configuration options before transformation into a CreateStackInput struct
// used by the cloudformation api.
type StackSpec struct {
	KeyPair   string

	// OnFailure determines what happens if stack creations fails.
	// Valid values are: "DO_NOTHING", "ROLLBACK", "DELETE"
	// Default: "ROLLBACK"
	OnFailure string

	Region    string

	StackName string
}

var (
	stackSpec = &StackSpec{
		OnFailure: "ROLLBACK",
	}
)

func createStack(svc *cf.CloudFormation, params []*cf.Parameter, stackName string, timeout int64) {
	input := &cf.CreateStackInput{
		StackName: aws.String(stackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String(stackSpec.OnFailure),
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
		aws.NewConfig().WithRegion(stackSpec.Region).WithLogLevel(aws.LogOff))

	params := []*cf.Parameter{
		{ParameterKey: aws.String("KeyName"), ParameterValue: aws.String(stackSpec.KeyPair)},
	}

	createStack(svc, params, stackSpec.StackName, 20)

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
		Use:   "awsplugin",
		Short: "init/update/destroy an AWS cluster in Docker swarm mode",
	}
	rootCmd.PersistentFlags().StringVarP(&stackSpec.KeyPair, "keypair", "k", "", "aws keypair name")
	rootCmd.PersistentFlags().StringVarP(&stackSpec.StackName, "stackname", "n", "", "aws stack name")
	rootCmd.PersistentFlags().StringVarP(&stackSpec.Region, "region", "r", "", "aws region")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster in swarm mode",
		Run:   provision,
	}
	initCmd.PersistentFlags().StringVar(&stackSpec.OnFailure, "onfailure", "", "")

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
