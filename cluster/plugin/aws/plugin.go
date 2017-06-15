package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	templateURL = "https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl"
)

// StackSpec stores raw configuration options before transformation into a CreateStackInput struct
// used by the cloudformation api.
type StackSpec struct {
	// OnFailure determines what happens if stack creations fails.
	// Valid values are: "DO_NOTHING", "ROLLBACK", "DELETE"
	// Default: "ROLLBACK"
	OnFailure string

	Params []string

	Region string

	StackName string
}

func parseParam(s string) *cf.Parameter {
	p := &cf.Parameter{}

	// split string to at most 2 substrings
	// if there is only 1 substring, then assume UsePreviousValue=true
	kv := strings.SplitN(s, "=", 2)
	p.SetParameterKey(kv[0])
	if len(kv) == 1 {
		p.SetUsePreviousValue(true)
	} else {
		p.SetParameterValue(kv[1])
	}

	return p
}

func CreateStack(svc *cf.CloudFormation, stackSpec *StackSpec, timeout int64) (*cf.CreateStackOutput, error) {
	sp := stackSpec.Params
	params := make([]*cf.Parameter, len(sp))
	for i := range sp {
		params[i] = parseParam(sp[i])
	}

	input := &cf.CreateStackInput{
		StackName: aws.String(stackSpec.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String(stackSpec.OnFailure),
		Parameters:       params,
		TemplateURL:      aws.String(templateURL),
		TimeoutInMinutes: aws.Int64(timeout),
	}

	return svc.CreateStack(input)
}
