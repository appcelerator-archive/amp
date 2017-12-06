package plugin

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	StackModeCreate = "create"
	StackModeDelete = "delete"
	StackModeUpdate = "update"
)

// used by the create function to parse the events and send meaningful information to the CLI
func (p *AWSPlugin) ScanEvents(mode string) error {
	var emptyReason string
	significantEventTypes := map[string]bool{
		"AWS::CloudFormation::Stack":         true,
		"AWS::EC2::VPC":                      true,
		"AWS::CloudFormation::WaitCondition": true,
		"AWS::AutoScaling::AutoScalingGroup": true,
		"AWS::EFS::FileSystem":               true,
	}
	significantEventStatuses := map[string]bool{
		cloudformation.StackStatusCreateInProgress:   false,
		cloudformation.StackStatusCreateComplete:     true,
		cloudformation.StackStatusCreateFailed:       true,
		cloudformation.StackStatusRollbackInProgress: true,
		cloudformation.StackStatusRollbackComplete:   true,
		cloudformation.StackStatusDeleteInProgress:   true,
	}
	eventInput := &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(p.Config.StackName),
		NextToken: nil,
	}
	eventIds := map[string]bool{}
	for {
		resp, err := p.cf.DescribeStackEvents(eventInput)
		if err != nil {
			if strings.Contains(err.Error(), "Throttling: Rate exceeded") {
				// ignore it, and continue processing the events
				time.Sleep(time.Second)
				continue
			} else if mode == StackModeDelete && strings.Contains(err.Error(), fmt.Sprintf("Stack [%s] does not exist", p.Config.StackName)) {
				// stack does not exist, probably because we've just successfuly deleted it
				event := StackEvent{
					EventId:           "NoSync-999",
					LogicalResourceId: p.Config.StackName,
					ResourceType:      "AWS::CloudFormation:Stack",
					ResourceStatus:    cloudformation.StackStatusDeleteComplete,
					Timestamp:         time.Now().Format(time.UnixDate),
				}
				j, err := PluginOutputToJSON(&event, nil, nil)
				if err != nil {
					return err
				}
				fmt.Println(j)
				return nil
			}
			return err
		}
		events := resp.StackEvents
		sort.Slice(events, func(i, j int) bool { return events[i].Timestamp.Unix() < events[j].Timestamp.Unix() })
		for _, se := range events {
			if eventIds[*se.EventId] {
				// already processed, ignore it
				continue
			}
			eventIds[*se.EventId] = true

			// filtering: stack events and creation / deletion of a few types
			if *se.ResourceType == "AWS::CloudFormation::Stack" ||
				(significantEventTypes[*se.ResourceType] == true && significantEventStatuses[*se.ResourceStatus] == true) {
				if se.ResourceStatusReason == nil {
					// we'll have an indirection, so better secure it
					se.ResourceStatusReason = &emptyReason
				}
				eventOutput := StackEvent{
					EventId:              *se.EventId,
					LogicalResourceId:    *se.LogicalResourceId,
					ResourceStatus:       *se.ResourceStatus,
					ResourceStatusReason: *se.ResourceStatusReason,
					ResourceType:         *se.ResourceType,
					Timestamp:            se.Timestamp.Format(time.UnixDate),
				}
				j, err := PluginOutputToJSON(&eventOutput, nil, nil)
				if err != nil {
					return err
				}
				fmt.Println(j)
			}
			// check end condition, based on the action
			if *se.ResourceType != "AWS::CloudFormation::Stack" {
				continue
			}
			switch mode {
			case StackModeCreate:
				switch *se.ResourceStatus {
				case cloudformation.StackStatusCreateComplete, cloudformation.StackStatusDeleteComplete, cloudformation.StackStatusRollbackComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			case StackModeDelete:
				switch *se.ResourceStatus {
				case cloudformation.StackStatusDeleteFailed, cloudformation.StackStatusDeleteComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			case StackModeUpdate:
				switch *se.ResourceStatus {
				case cloudformation.ResourceStatusUpdateFailed, cloudformation.ResourceStatusUpdateComplete, cloudformation.StackStatusRollbackComplete:
					return nil
				default:
					// not an end condition, continue the loop
				}
			default:
				return fmt.Errorf("unknown mode: %s", mode)
			}
		}
		// waiting before reading next batch of events (if too short, the API throttles)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
