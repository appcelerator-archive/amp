package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	// replace latest by the version when releasing the plugin
	DefaultTemplateURL = "https://s3.amazonaws.com/io-amp-binaries/templates/latest/aws-swarm-asg.yml"
)

// Configuration stores raw request input options before transformation into a AWS SDK specific structs  used by the cloudformation api.
type Configuration struct {
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
	// AWS Access Key ID
	AccessKeyId string
	// AWS Secret Access Key
	SecretAccessKey string
	// AWS Profile
	Profile string
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
	Event  *StackEvent   `json:"Event"`
	Error  string        `json:"Error"`
}

// AWSPlugin is the AWS plugin
type AWSPlugin struct {
	Config *Configuration
	cf     *cloudformation.CloudFormation
}

var (
	Config = &Configuration{
		OnFailure:   "DO_NOTHING",
		Params:      []string{},
		TemplateURL: DefaultTemplateURL,
		Profile:     "default",
	}
	AWS *AWSPlugin
)

// New returns a new AWS Plugin instance
func New(config *Configuration) *AWSPlugin {
	// export vars if creds are passed as arguments
	if config.AccessKeyId != "" && config.SecretAccessKey != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", config.AccessKeyId)
		os.Setenv("AWS_SECRET_ACCESS_KEY", config.SecretAccessKey)
	} else {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		os.Setenv("AWS_PROFILE", config.Profile)
	}

	awsConfig := aws.NewConfig().WithLogLevel(aws.LogOff)

	// region can be set with a CLI option, but if not set it can be set by the Config file
	if config.Region != "" {
		awsConfig.Region = aws.String(config.Region)
	}

	// Create the service's client with the session.
	return &AWSPlugin{
		Config: config,
		cf:     cloudformation.New(session.Must(session.NewSession()), awsConfig),
	}
}

func parseParam(s string) *cloudformation.Parameter {
	p := &cloudformation.Parameter{}

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

func toParameters(sa []string) []*cloudformation.Parameter {
	params := make([]*cloudformation.Parameter, len(sa))
	for i := range sa {
		params[i] = parseParam(sa[i])
	}
	return params
}

// CreateStack starts the AWS stack creation operation
// The operation will return immediately
func (p *AWSPlugin) CreateStack(ctx context.Context, timeout int64) (*cloudformation.CreateStackOutput, error) {
	input := &cloudformation.CreateStackInput{
		StackName: aws.String(p.Config.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String(p.Config.OnFailure),
		Parameters:       toParameters(p.Config.Params),
		TemplateURL:      aws.String(p.Config.TemplateURL),
		TimeoutInMinutes: aws.Int64(timeout),
	}
	return p.cf.CreateStackWithContext(ctx, input)
}

// InfoStack returns the output information that was produced when the stack was created or updated
func (p *AWSPlugin) InfoStack(ctx context.Context) ([]StackOutput, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: aws.String(p.Config.StackName),
		NextToken: aws.String(strconv.Itoa(p.Config.Page)),
	}
	output, err := p.cf.DescribeStacksWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	var stack *cloudformation.Stack
	for _, stack = range output.Stacks {
		//if stack.StackName == input.StackName {
		if aws.StringValue(stack.StackName) == p.Config.StackName {
			break
		}
		stack = nil
	}

	if stack == nil {
		return nil, errors.New("stack not found: " + p.Config.StackName)
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
func (p *AWSPlugin) UpdateStack(ctx context.Context) (*cloudformation.UpdateStackOutput, error) {
	input := &cloudformation.UpdateStackInput{
		StackName: aws.String(p.Config.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		Parameters:  toParameters(p.Config.Params),
		TemplateURL: aws.String(p.Config.TemplateURL),
	}
	return p.cf.UpdateStackWithContext(ctx, input)
}

// DeleteStack starts the delete operation
// The operation will return immediately
func (p *AWSPlugin) DeleteStack(ctx context.Context) (*cloudformation.DeleteStackOutput, error) {
	input := &cloudformation.DeleteStackInput{
		StackName: aws.String(p.Config.StackName),
	}
	return p.cf.DeleteStackWithContext(ctx, input)
}

// PluginOutputToJSON is a helper function that wrap an event, an error or a slice of StackOutput
// to a JSON string representation of PluginOutput
func PluginOutputToJSON(ev *StackEvent, so []StackOutput, e error) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9:;,.\\-()\\[\\]'\" ]+")
	if err != nil {
		return "", err
	}

	po := PluginOutput{
		Output: so,
		Event:  ev,
	}
	if e != nil {
		po.Error = reg.ReplaceAllString(e.Error(), "%")
	}
	j, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// ListStack lists the stacks based on the filter on stack status
func (p *AWSPlugin) ListStack(ctx context.Context) (*cloudformation.ListStacksOutput, error) {
	statusFilter := []string{
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusRollbackComplete,
	}
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice(statusFilter),
	}
	return p.cf.ListStacksWithContext(ctx, input)
}
