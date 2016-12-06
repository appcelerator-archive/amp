package tests

import (
	. "github.com/appcelerator/amp/api/rpc/function"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestFunctionShouldCreateAndDeleteAFunction(t *testing.T) {
	created, err := functionClient.Create(ctx, &CreateRequest{
		Function: &FunctionEntry{
			Name:  "test-function",
			Image: "test-image",
		},
	})
	assert.NoError(t, err)

	_, err = functionClient.Delete(ctx, &DeleteRequest{
		Id: created.Function.Id,
	})
	assert.NoError(t, err)
}

func TestFunctionShouldFailWhenCreatingAnAlreadyExistingFunction(t *testing.T) {
	created, err := functionClient.Create(ctx, &CreateRequest{
		Function: &FunctionEntry{
			Name:  "test-function",
			Image: "test-image",
		},
	})
	assert.NoError(t, err)

	_, err = functionClient.Create(ctx, &CreateRequest{
		Function: &FunctionEntry{
			Name:  "test-function",
			Image: "test-image",
		},
	})
	assert.Error(t, err)

	_, err = functionClient.Delete(ctx, &DeleteRequest{
		Id: created.Function.Id,
	})
	assert.NoError(t, err)
}

func TestFunctionShouldListCreatedFunctions(t *testing.T) {
	r1, err := functionClient.Create(ctx, &CreateRequest{Function: &FunctionEntry{Name: "test-function-1", Image: "test-image-1"}})
	assert.NoError(t, err)
	r2, err := functionClient.Create(ctx, &CreateRequest{Function: &FunctionEntry{Name: "test-function-2", Image: "test-image-2"}})
	assert.NoError(t, err)
	r3, err := functionClient.Create(ctx, &CreateRequest{Function: &FunctionEntry{Name: "test-function-3", Image: "test-image-3"}})
	assert.NoError(t, err)

	reply, err := functionClient.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.Contains(t, reply.Functions, r1.Function)
	assert.Contains(t, reply.Functions, r2.Function)
	assert.Contains(t, reply.Functions, r3.Function)

	_, err = functionClient.Delete(ctx, &DeleteRequest{Id: r1.Function.Id})
	assert.NoError(t, err)
	_, err = functionClient.Delete(ctx, &DeleteRequest{Id: r2.Function.Id})
	assert.NoError(t, err)
	_, err = functionClient.Delete(ctx, &DeleteRequest{Id: r3.Function.Id})
	assert.NoError(t, err)

	reply, err = functionClient.List(ctx, &ListRequest{})
	assert.NoError(t, err)
	assert.NotContains(t, reply.Functions, r1.Function)
	assert.NotContains(t, reply.Functions, r2.Function)
	assert.NotContains(t, reply.Functions, r3.Function)
}

func TestFunctionShouldCreateInvokeAndDeleteAFunction(t *testing.T) {
	created, err := functionClient.Create(ctx, &CreateRequest{
		Function: &FunctionEntry{
			Name:  "test-" + stringid.GenerateNonCryptoID(),
			Image: "appcelerator/amp-demo-function",
		},
	})
	assert.NoError(t, err)

	testInput := "This is a test input that is supposed to be capitalized"
	resp, err := http.Post("http://amp-function-listener/"+created.Function.Id, "text/plain;charset=utf-8", strings.NewReader(testInput))
	assert.NoError(t, err)

	output, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, strings.Title(testInput), string(output))

	_, err = functionClient.Delete(ctx, &DeleteRequest{
		Id: created.Function.Id,
	})
	assert.NoError(t, err)
}
