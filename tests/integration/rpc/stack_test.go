package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/stretchr/testify/assert"
)

const (
	example1 = `
version: "3"
services:
    pinger1:
        image: appcelerator/pinger
        deploy:
            mode: replicated
            replicas: 2

    pinger2:
        image: appcelerator/pinger
        deploy:
            mode: replicated
            replicas: 2
            labels:
                io.amp.mapping: "www:3000"
`
)

//Test two stacks life cycle in the same time
func TestStackShouldManageStackLifeCycleSuccessfully(t *testing.T) {
	//Start stack essai1
	name := fmt.Sprintf("test-%d", time.Now().Unix())
	//Start stack test
	s1 := &stack.Stack{
		Name:     name,
		FileData: example1,
	}
	rUp, errUp := stackClient.Deploy(ctx, &stack.StackDeployRequest{Stack: s1})
	if errUp != nil {
		t.Fatal(errUp)
	}
	assert.NotEmpty(t, rUp.Answer, "Stack test1 reply should not be empty")
	time.Sleep(3 * time.Second)
	//verifyusing ls
	t.Log("perform stack ls")
	listRequest := stack.StackRequest{}
	_, errls := stackClient.ListStacks(ctx, &listRequest)
	if errls != nil {
		t.Fatal(errls)
	}
	//Prepare requests
	stackRequest := stack.StackRequest{
		Name: name,
	}
	//Stop stack test
	t.Log("stop stack " + name)
	rStop, errStop := stackClient.Remove(ctx, &stackRequest)
	if errStop != nil {
		t.Fatal(errStop)
	}
	assert.NotEmpty(t, rStop.Answer, "Stack test StackId should not be empty")
}
