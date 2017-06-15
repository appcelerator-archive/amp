package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
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

	Params []*cf.Parameter

	Region    string

	StackName string
}

func CreateStack(svc *cf.CloudFormation, stackSpec *StackSpec, timeout int64) {
	input := &cf.CreateStackInput{
		StackName: aws.String(stackSpec.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String(stackSpec.OnFailure),
		Parameters:       stackSpec.Params,
		TemplateURL:      aws.String(templateURL),
		TimeoutInMinutes: aws.Int64(timeout),
	}

	resp, err := svc.CreateStack(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))
}

