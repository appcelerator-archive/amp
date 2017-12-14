package aws

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cf "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

type MetadataKey string

const (
	AwsAvailabilityZoneURL        = MetadataKey("http://169.254.169.254/latest/meta-data/placement/availability-zone/")
	AwsInstanceIdURL              = MetadataKey("http://169.254.169.254/latest/meta-data/instance-id")
	AwsStackNameTag               = "aws:cloudformation:stack-name"
	ErrMsgBucketAlreadyOwnedByYou = s3.ErrCodeBucketAlreadyOwnedByYou
	ErrMsgBucketAlreadyExists     = s3.ErrCodeBucketAlreadyExists
	ErrMsgNoSuchBucket            = s3.ErrCodeNoSuchBucket
	ErrMsgInvalidACL              = "Invalid ACL"
)

func getTag(ec2svc *ec2.EC2, instanceId string, tagname string) (string, error) {
	instanceIdKey := "resource-id"
	ec2TagKey := "key"
	instanceFilter := ec2.Filter{Name: &instanceIdKey, Values: []*string{&instanceId}}
	tagFilter := ec2.Filter{Name: &ec2TagKey, Values: []*string{&tagname}}
	input := ec2.DescribeTagsInput{Filters: []*ec2.Filter{&instanceFilter, &tagFilter}}
	output, err := ec2svc.DescribeTags(&input)
	if err != nil {
		return "", err
	}
	if len(output.Tags) != 1 {
		return "", errors.New("couldn't read the tag on the aws instance")
	}
	tagDescription := *output.Tags[0]
	return *tagDescription.Value, nil
}

// getMetadata returns the value of the metadata by key
func getMetadata(key MetadataKey) (string, error) {
	timeout := time.Duration(5 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}
	resp, err := httpClient.Get(string(key))
	if err != nil {
		return "", err
	}
	buff, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

// Region get the region from the metadata service
func Region() (string, error) {
	az, err := getMetadata(AwsAvailabilityZoneURL)
	if err != nil {
		return "", err
	}
	return az[:len(az)-1], nil
}

// StackInfo updates the receiver with cloudformation stack information
func StackInfo(ctx context.Context, outputs *map[string]string) error {
	region, err := Region()
	if err != nil {
		return err
	}
	// get the instance id
	instanceId, err := getMetadata(AwsInstanceIdURL)
	if err != nil {
		return err
	}
	// create an aws api session
	awsconfig := awssdk.NewConfig().WithRegion(region).WithLogLevel(awssdk.LogOff)
	sess := session.Must(session.NewSession())
	// instance of cloudformation and ec2 services
	cfsvc := cf.New(sess, awsconfig)
	ec2svc := ec2.New(sess, awsconfig)
	// get the stack name
	stackName, err := getTag(ec2svc, instanceId, AwsStackNameTag)
	if err != nil {
		return err
	}
	page := 1
	input := &cf.DescribeStacksInput{
		StackName: awssdk.String(stackName),
		NextToken: awssdk.String(strconv.Itoa(page)),
	}
	output, err := cfsvc.DescribeStacksWithContext(ctx, input)
	if err != nil {
		return err
	}

	var stack *cf.Stack
	for _, stack = range output.Stacks {
		if awssdk.StringValue(stack.StackName) == stackName {
			break
		}
		stack = nil
	}

	if stack == nil {
		return errors.New("stack not found: " + stackName)
	}

	*outputs = map[string]string{
		"StackName": stackName,
		"Provider":  "AWS", // can't import constants from package cloud
		"Region":    region,
	}
	for _, o := range stack.Outputs {
		switch awssdk.StringValue(o.OutputKey) {
		case "DNSTarget", "NFSEndpoint", "InternalPKITarget", "InternalDockerHost":
			(*outputs)[awssdk.StringValue(o.OutputKey)] = awssdk.StringValue(o.OutputValue)
		}
	}

	return nil
}

// S3CreateBucket creates an s3 bucket
func S3CreateBucket(name string, acl string, region string) (string, error) {
	switch acl {
	case s3.BucketCannedACLPrivate, s3.BucketCannedACLPublicRead, s3.BucketCannedACLPublicReadWrite, s3.BucketCannedACLAuthenticatedRead:
	default:
		return "", errors.New(ErrMsgInvalidACL)
	}

	sess := session.Must(session.NewSession(&awssdk.Config{
		Region: awssdk.String(region),
	}))
	svc := s3.New(sess)

	input := &s3.CreateBucketInput{
		ACL:    awssdk.String(acl),
		Bucket: awssdk.String(name),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: awssdk.String(region),
		},
	}
	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Infoln(aerr.Error())
			err = errors.New(aerr.Code())
		}
		return "", err
	}
	return *result.Location, nil
}

func S3ListBuckets(region string) (*[]string, error) {
	var buckets []string
	sess := session.Must(session.NewSession(&awssdk.Config{
		Region: awssdk.String(region),
	}))
	svc := s3.New(sess)

	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Infoln(aerr.Error())
			return nil, errors.New(aerr.Code())
		}
		return nil, err
	}
	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	return &buckets, nil
}

func S3GetBucketLocation(bucket string, region string) (string, error) {
	sess := session.Must(session.NewSession(&awssdk.Config{
		Region: awssdk.String(region),
	}))
	svc := s3.New(sess)
	l, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: awssdk.String(bucket)})
	if err != nil {
		return "", err
	}
	return *l.LocationConstraint, nil
}

func S3DeleteBucket(name string, region string, force bool) error {
	sess := session.Must(session.NewSession(&awssdk.Config{
		Region: awssdk.String(region),
	}))
	svc := s3.New(sess)

	if force {
		result, err := svc.ListObjects(&s3.ListObjectsInput{
			Bucket:  awssdk.String(name),
			MaxKeys: awssdk.Int64(1024),
		})
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				log.Infoln(aerr.Error())
				return errors.New(aerr.Code())
			} else {
				return err
			}
		}
		if len(result.Contents) == 0 {
			log.Infoln("the force deletion flag was set, but the bucket was empty, ignoring the flag")

		} else {
			objectsToDelete := []*s3.ObjectIdentifier{}
			for _, o := range result.Contents {
				objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{Key: o.Key})
			}
			_, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
				Bucket: awssdk.String(name),
				Delete: &s3.Delete{
					Objects: objectsToDelete,
				},
			})
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					log.Infoln(aerr.Error())
					return errors.New(aerr.Code())
				}
				return err
			}
			log.Infoln(len(result.Contents), "objects removed in bucket", name)
		}

	}
	_, err := svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: awssdk.String(name)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Infoln(aerr.Error())
			return errors.New(aerr.Code())
		}
		return err
	}
	return nil
}
