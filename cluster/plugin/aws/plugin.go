package aws

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	DefaultTemplateURL = "https://s3.amazonaws.com/io-amp-binaries/templates/latest/aws-swarm-asg.yml"
)

// RequestOptions stores raw request input options before transformation into a AWS SDK specific
// structs  used by the cloudformation api.
type RequestOptions struct {
	// OnFailure determines what happens if stack creations fails.
	// Valid values are: "DO_NOTHING", "ROLLBACK", "DELETE"
	// Default: "ROLLBACK"
	OnFailure string

	// Params are for parameters supported by the CloudFormation template that will be used
	Params []string

	// Page is for aws requests that return paged information (pages start at 1)
	Page int

	// Region is the AWS region, ex: us-west-2
	Region string

	// StackName is the user-supplied name for identifying the stack
	StackName string

	// Sync, if true, causes the create operation to block until finished
	Sync bool

	// TemplateURL is the URL for the AWS CloudFormation to use
	TemplateURL string
}

// StackOutput contains the converted output from the create stack operation
type StackOutput struct {
	// Description is the user defined description associated with the output
	Description string `json:"description"`

	// OutputKey is the key associated with the output
	OutputKey string `json:"key"`

	// OutputValue is the value associated with the output
	OutputValue string `json:"value"`
}

// StackOutputList is used as a container for output by the StackOutputToJSON helper function
type StackOutputList struct {
	// Output is a slice of StackOutput
	Output []StackOutput `json:"output""`
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

func toParameters(sa []string) []*cf.Parameter {
	params := make([]*cf.Parameter, len(sa))
	for i := range sa {
		params[i] = parseParam(sa[i])
	}
	return params
}

// CreateStack starts the AWS stack creation operation
// The operation will return immediately unless opts.Sync is true
func CreateStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions, timeout int64) (*cf.CreateStackOutput, error) {
	input := &cf.CreateStackInput{
		StackName: aws.String(opts.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String(opts.OnFailure),
		Parameters:       toParameters(opts.Params),
		TemplateURL:      aws.String(opts.TemplateURL),
		TimeoutInMinutes: aws.Int64(timeout),
	}

	return svc.CreateStackWithContext(ctx, input)
}

// InfoStack returns the output information that was produced when the stack was created or updated
func InfoStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions) ([]StackOutput, error) {
	input := &cf.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
		NextToken: aws.String(strconv.Itoa(opts.Page)),
	}
	output, err := svc.DescribeStacksWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var stack *cf.Stack
	for _, stack = range output.Stacks {
		//if stack.StackName == input.StackName {
		if aws.StringValue(stack.StackName) == opts.StackName {
			break
		}
		stack = nil
	}

	if stack == nil {
		return nil, errors.New("stack not found: " + opts.StackName)
	}

	stackOutputs := []StackOutput{}
	for _, o := range stack.Outputs {
		stackOutputs = append(stackOutputs, StackOutput{
			Description: aws.StringValue(o.Description),
			OutputKey:   aws.StringValue(o.OutputKey),
			OutputValue: aws.StringValue(o.OutputValue),
		})
	}

	return stackOutputs, nil
}

// UpdateStack starts the update operation
// The operation will return immediately unless opts.Sync is true
func UpdateStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions) (*cf.UpdateStackOutput, error) {
	input := &cf.UpdateStackInput{
		StackName: aws.String(opts.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		Parameters:  toParameters(opts.Params),
		TemplateURL: aws.String(opts.TemplateURL),
	}

	return svc.UpdateStackWithContext(ctx, input)
}

// DeleteStack starts the delete operation
// The operation will return immediately unless opts.Sync is true
func DeleteStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions) (*cf.DeleteStackOutput, error) {
	input := &cf.DeleteStackInput{
		StackName: aws.String(opts.StackName),
	}

	return svc.DeleteStackWithContext(ctx, input)
}

// StackOutputToJSON is a helper function that converts a slice of StackOutput to a JSON string representation
// of StackOutputList
func StackOutputToJSON(so []StackOutput) (string, error) {
	list := StackOutputList{
		Output: so,
	}
	j, err := json.Marshal(list)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
