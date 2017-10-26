package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	// replace latest by the version when releasing the plugin
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
// Meant to be included in PluginOutput
type StackOutput struct {
	// Description is the user defined description associated with the output
	Description string `json:"description"`
	// OutputKey is the key associated with the output
	OutputKey string `json:"key"`
	// OutputValue is the value associated with the output
	OutputValue string `json:"value"`
}

// StackEvent contains the converted event from the create stack operation
// Similar to the structure from the aws SDK
// Meant to be included in PluginOutput
type StackEvent struct {
	ClientRequestToken   string `json:"ClientRequestToken"`
	EventId              string `json:"EventId"`
	LogicalResourceId    string `json:"LogicalResourceId"`
	PhysicalResourceId   string `json:"PhysicalResourceId"`
	ResourceProperties   string `json:"ResourceProperties"`
	ResourceStatus       string `json:"ResourceStatus"`
	ResourceStatusReason string `json:"ResourceStatusReason"`
	ResourceType         string `json:"ResourceType"`
	StackId              string `json:"StackId"`
	StackName            string `json:"StackName"`
	Timestamp            string `json:"Timestamp"`
}

// PluginOutput contains the stack output, the stack events and the errors that the plugin
// will transmit to the CLI. it's not a single envelop, it can be repeated
type PluginOutput struct {
	Output []StackOutput `json:"Output"`
	Event  *StackEvent   `"json:"Event"`
	Error  string        `"json:"Error"`
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
// The operation will return immediately
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
// The operation will return immediately
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
// The operation will return immediately
func DeleteStack(ctx context.Context, svc *cf.CloudFormation, opts *RequestOptions) (*cf.DeleteStackOutput, error) {
	input := &cf.DeleteStackInput{
		StackName: aws.String(opts.StackName),
	}

	return svc.DeleteStackWithContext(ctx, input)
}

// PluginOutputToJSON is a helper function that wrap an event, an error or a slice of StackOutput
// to a JSON string representation of PluginOutput
func PluginOutputToJSON(ev *StackEvent, so []StackOutput, e error) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9:;,. ]+")
	if err != nil {
		return "", err
	}

	po := PluginOutput{
		Output: so,
		Event:  ev,
	}
	if e != nil {
		po.Error = reg.ReplaceAllString(e.Error(), "")
	}
	j, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// ListStack lists the stacks based on the filter on stack status
func ListStack(ctx context.Context, svc *cf.CloudFormation) (*cf.ListStacksOutput, error) {
	statusFilter := []string{cf.StackStatusCreateFailed, cf.StackStatusCreateInProgress, cf.ChangeSetStatusCreateComplete, cf.StackStatusRollbackInProgress, cf.StackStatusRollbackFailed, cf.StackStatusDeleteInProgress}
	input := &cf.ListStacksInput{
		StackStatusFilter: aws.StringSlice(statusFilter),
	}

	return svc.ListStacksWithContext(ctx, input)
}
