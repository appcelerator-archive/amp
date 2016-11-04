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
services:
  pinger:
    image: appcelerator/pinger
    replicas: 2
  pingerExt1:
    image: appcelerator/pinger
    replicas: 2
    public:
      - name: www1
        protocol: tcp
        internal_port: 3000
  pingerExt2:
    image: appcelerator/pinger
    replicas: 2
    public:
      - name: www2
        protocol: tcp
        publish_port: 3001
        internal_port: 3000`

	example2 = `
services:
  pinger:
    image: appcelerator/pinger
    replicas: 2
  pingerExt1:
    image: appcelerator/pinger
    replicas: 2
    public:
      - name: www1
        protocol: tcp
        internal_port: 3000
  pingerExt2:
    image: appcelerator/pinger
    replicas: 2
    public:
      - name: www2
        protocol: tcp
        publish_port: 3002
        internal_port: 3000`
)

//Test two stacks life cycle in the same time
func TestStackShouldManageStackLifeCycleSuccessfully(t *testing.T) {
	//Start stack essai1
	name := fmt.Sprintf("test-%d", time.Now().Unix())
	//Start stack test
	t.Log("start stack " + name)
	rUp, errUp := stackClient.Up(ctx, &stack.StackFileRequest{StackName: name, Stackfile: example1})
	if errUp != nil {
		t.Fatal(errUp)
	}
	assert.NotEmpty(t, rUp.StackId, "Stack test1 StackId should not be empty")
	time.Sleep(3 * time.Second)
	//verifyusing ls
	t.Log("perform stack ls")
	listRequest := stack.ListRequest{}
	_, errls := stackClient.List(ctx, &listRequest)
	if errls != nil {
		t.Fatal(errls)
	}
	//Prepare requests
	stackRequest := stack.StackRequest{
		StackIdent: rUp.StackId,
	}
	//Stop stack test
	t.Log("stop stack " + name)
	rStop, errStop := stackClient.Stop(ctx, &stackRequest)
	if errStop != nil {
		t.Fatal(errStop)
	}
	assert.NotEmpty(t, rStop.StackId, "Stack test StackId should not be empty")
	//Restart stack test
	t.Log("restart stack " + name)
	rRestart, errRestart := stackClient.Start(ctx, &stackRequest)
	if errRestart != nil {
		t.Fatal(errRestart)
	}
	assert.NotEmpty(t, rRestart.StackId, "Stack test StackId should not be empty")
	time.Sleep(3 * time.Second)
	t.Log("remove stack " + name)
	//Remove stack test
	removeRequest := stack.RemoveRequest{
		StackIdent: rUp.StackId,
		Force:      true,
	}
	rRemove, errRemove := stackClient.Remove(ctx, &removeRequest)
	if errRemove != nil {
		t.Fatal(errRemove)
	}
	assert.NotEmpty(t, rRemove.StackId, "Stack test StackId should not be empty")
}
