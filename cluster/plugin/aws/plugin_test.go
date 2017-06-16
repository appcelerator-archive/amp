package aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/satori/go.uuid"
)

var (
	keyName string
	region  string

	sess *session.Session
	svc  *cf.CloudFormation
)

// REQUIREMENTS:
// AWS credentials in path
// Environment variables set for the following:
//   KEYPAIR=<your-aws-key-pair>
//   REGION=<aws-region>
func init() {
	keyName = os.Getenv("KEYNAME")
	if keyName == "" {
		log.Fatal("KEYNAME environment variable not set")
	}
	region = os.Getenv("REGION")
	if region == "" {
		log.Fatal("REGION environment variable not set")
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	sess = session.Must(session.NewSession())

	// Create the service's client with the session.
	svc = cf.New(sess,
		aws.NewConfig().WithRegion(region).WithLogLevel(aws.LogOff))
}

func teardown() {
	// TODO destroy test stacks
}

func TestCreate(t *testing.T) {
	stackName := fmt.Sprintf("%s-plugin-test-%s", keyName, uuid.NewV4())

	opts := &RequestOptions{
		Region:    region,
		StackName: stackName,
		OnFailure: "DELETE",
		Params: []string{
			fmt.Sprintf("KeyName=%s", keyName),
		},
	}

	// create stack
	// ============
	ctxCreate := context.Background()
	respCreate, err := CreateStack(ctxCreate, svc, opts, 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(awsutil.StringValue(respCreate))
	t.Log("waiting...")
	input := &cf.DescribeStacksInput{
		StackName: aws.String(opts.StackName),
	}
	if err := svc.WaitUntilStackUpdateCompleteWithContext(ctxCreate, input); err != nil {
		t.Fatal(err)
	}
	t.Logf("stack created: %s\n", opts.StackName)

	// update stack
	// ============
	opts.Params = append(opts.Params, fmt.Sprintf("ClusterSize=%s", "1"))
	ctxUpdate := context.Background()
	respUpdate, err := UpdateStack(ctxUpdate, svc, opts)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(awsutil.StringValue(respUpdate))
	t.Log("waiting...")
	if err := svc.WaitUntilStackUpdateCompleteWithContext(ctxUpdate, input); err != nil {
		t.Fatal(err)
	}
	t.Logf("stack update: %s\n", opts.StackName)

	// delete stack
	// ============
	ctxDelete := context.Background()
	respDelete, err := DeleteStack(ctxDelete, svc, opts)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(awsutil.StringValue(respDelete))
	t.Log("waiting...")
	if err := svc.WaitUntilStackDeleteCompleteWithContext(ctxDelete, input); err != nil {
		t.Fatal(err)
	}
	t.Logf("stack deleted: %s\n", opts.StackName)
}
