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
	keyName string
	region  string

	sess *session.Session
	svc *cf.CloudFormation
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

	stackSpec := &StackSpec{
		Region:    region,
		StackName: stackName,
		OnFailure: "DELETE",
		Params: []string{
			fmt.Sprintf("KeyName=%s", keyName),
		},
	}

	_, err := CreateStack(svc, stackSpec, 20)
	if err != nil {
		t.Fatal(err)
	}
}
