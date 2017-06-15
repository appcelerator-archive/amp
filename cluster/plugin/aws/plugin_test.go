package aws

import (
	"fmt"
	"log"
	"os"
	"testing"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/satori/go.uuid"
)

var (
	keyPair string
	region string

	sess *session.Session
	svc *cf.CloudFormation
)

// REQUIREMENTS:
// AWS credentials in path
// Environment variables set for the following:
//   KEYPAIR=<your-aws-key-pair>
//   REGION=<aws-region>
func init() {
	keyPair = os.Getenv("KEYPAIR")
	if keyPair == "" {
		log.Fatal("KEYPAIR environment variable not set")
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
	stackName := fmt.Sprintf("%s-plugin-test-%s", keyPair, uuid.NewV4())

	stackSpec := &StackSpec{
		KeyPair: keyPair,
		Region: region,
		StackName: stackName,
		OnFailure: "DELETE",
	}

	stackSpec.Params = []*cf.Parameter{
		{ParameterKey: aws.String("KeyName"), ParameterValue: aws.String(keyPair)},
	}

	_, err := CreateStack(svc, stackSpec, 20)
	if err != nil {
		t.Fatal(err)
	}
}
