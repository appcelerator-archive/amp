package aws

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	DefaultTemplateURL = "https://editions-us-east-1.s3.amazonaws.com/aws/edge/Docker.tmpl"
)

// RequestOptions stores raw request input options before transformation into a AWS SDK specific
// structs  used by the cloudformation api.
type RequestOptions struct {
	// OnFailure determines what happens if stack creations fails.
	// Valid values are: "DO_NOTHING", "ROLLBACK", "DELETE"
	// Default: "ROLLBACK"
	OnFailure string

	Params []string

	Region string

	StackName string

	Sync bool

	TemplateURL string
}

// StackOutput contains the converted output from the create stack operation
type StackOutput struct {
	// Description is the user defined description associated with the output
	Description string

	// OutputKey is the key associated with the output
	OutputKey string

	// OutputValue is the value associated with the output
	OutputValue string
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

func DeleteStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions) (*cf.DeleteStackOutput, error) {
	input := &cf.DeleteStackInput{
		StackName: aws.String(opts.StackName),
	}

	return svc.DeleteStackWithContext(ctx, input)
}

func describeStack(ctx context.Context, svc *cf.CloudFormation, id string, page int) ([]StackOutput, error) {
	input := &cf.DescribeStacksInput{
		StackName: aws.String(id),
		NextToken: aws.String(strconv.Itoa(page)),
	}
	output, err := svc.DescribeStacksWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var stack *cf.Stack
	for _, stack = range output.Stacks {
		n := stack.StackName
		fmt.Println(n)
		if aws.StringValue(stack.StackName) == id {
			break
		}
		stack = nil
	}

	if stack == nil {
		return nil, errors.New("stack not found: " + id)
	}

	stackOutput := []StackOutput{}
	for _, o := range stack.Outputs {
		stackOutput = append(stackOutput, StackOutput{
			Description: aws.StringValue(o.Description),
			OutputKey: aws.StringValue(o.OutputKey),
			OutputValue: aws.StringValue(o.OutputValue),
		})
	}

	return stackOutput, nil
}
