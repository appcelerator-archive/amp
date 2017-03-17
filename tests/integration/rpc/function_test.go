package tests

import (
	"github.com/appcelerator/amp/api/rpc/account"
	. "github.com/appcelerator/amp/api/rpc/function"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestFunctionShouldCreateAndDeleteAFunction(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())
	functionStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	created, err := functionClient.Create(ownerCtx, &CreateRequest{
		Name:  "test-function",
		Image: "test-image",
	})
	assert.NoError(t, err)

	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{
		Id: created.Function.Id,
	})
	assert.NoError(t, err)
}

func TestFunctionShouldFailWhenCreatingAnAlreadyExistingFunction(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())
	functionStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	_, err := functionClient.Create(ownerCtx, &CreateRequest{
		Name:  "test-function",
		Image: "test-image",
	})
	assert.NoError(t, err)

	_, err = functionClient.Create(ownerCtx, &CreateRequest{
		Name:  "test-function",
		Image: "test-image",
	})
	assert.Error(t, err)
}

func TestFunctionShouldListCreatedFunctions(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())
	functionStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	r1, err := functionClient.Create(ownerCtx, &CreateRequest{Name: "test-function-1", Image: "test-image-1"})
	assert.NoError(t, err)
	r2, err := functionClient.Create(ownerCtx, &CreateRequest{Name: "test-function-2", Image: "test-image-2"})
	assert.NoError(t, err)
	r3, err := functionClient.Create(ownerCtx, &CreateRequest{Name: "test-function-3", Image: "test-image-3"})
	assert.NoError(t, err)

	reply, err := functionClient.List(ownerCtx, &ListRequest{})
	assert.NoError(t, err)
	assert.Contains(t, reply.Functions, r1.Function)
	assert.Contains(t, reply.Functions, r2.Function)
	assert.Contains(t, reply.Functions, r3.Function)

	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{Id: r1.Function.Id})
	assert.NoError(t, err)
	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{Id: r2.Function.Id})
	assert.NoError(t, err)
	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{Id: r3.Function.Id})
	assert.NoError(t, err)

	reply, err = functionClient.List(ownerCtx, &ListRequest{})
	assert.NoError(t, err)
	assert.NotContains(t, reply.Functions, r1.Function)
	assert.NotContains(t, reply.Functions, r2.Function)
	assert.NotContains(t, reply.Functions, r3.Function)
}

func TestFunctionShouldCreateInvokeAndDeleteAFunction(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())
	functionStore.Reset(context.Background())

	// Create a user
	ownerCtx := createUser(t, &testUser)

	created, err := functionClient.Create(ownerCtx, &CreateRequest{
		Name:  "test-" + stringid.GenerateNonCryptoID(),
		Image: "appcelerator/amp-demo-function",
	})
	assert.NoError(t, err)

	testInput := "This is a test input that is supposed to be capitalized"
	resp, err := http.Post("http://amp-function-listener/"+created.Function.Id, "text/plain;charset=utf-8", strings.NewReader(testInput))
	assert.NoError(t, err)

	output, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, strings.Title(testInput), string(output))

	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{
		Id: created.Function.Id,
	})
	assert.NoError(t, err)
}

func TestFunctionPermissions(t *testing.T) {
	// Reset the storage
	accountStore.Reset(context.Background())
	functionStore.Reset(context.Background())

	// Create an organization and its owner
	ownerCtx := createOrganization(t, &testOrg, &testUser)

	// Create and add a member to the organization, then promote to org owner
	memberCtx := createAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
	changeOrganizationMemberRole(ownerCtx, t, &testOrg, &testMember, accounts.OrganizationRole_ORGANIZATION_OWNER)

	// Create a user function
	userFunction, err := functionClient.Create(ownerCtx, &CreateRequest{Name: "user-function", Image: "hello-world"})
	assert.NoError(t, err)

	// Create an org function
	orgContext := switchAccount(ownerCtx, t, testOrg.Name)
	orgFunction, err := functionClient.Create(orgContext, &CreateRequest{Name: "org-function", Image: "hello-world"})
	assert.NoError(t, err)

	// Remove user from organization
	_, err = accountClient.RemoveUserFromOrganization(memberCtx, &account.RemoveUserFromOrganizationRequest{OrganizationName: testOrg.Name, UserName: testUser.Name})
	assert.NoError(t, err)

	// Member cannot delete user function
	_, err = functionClient.Delete(memberCtx, &DeleteRequest{Id: userFunction.Function.Id})
	assert.Error(t, err)

	// User cannot delete org function
	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{Id: orgFunction.Function.Id})
	assert.Error(t, err)

	// Member can delete his own function
	_, err = functionClient.Delete(memberCtx, &DeleteRequest{Id: orgFunction.Function.Id})
	assert.NoError(t, err)

	// User can delete his own function
	_, err = functionClient.Delete(ownerCtx, &DeleteRequest{Id: userFunction.Function.Id})
	assert.NoError(t, err)
}
