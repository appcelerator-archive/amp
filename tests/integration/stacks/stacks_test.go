package stack

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/account"
	. "github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/docker/docker/pkg/stringid"
	"github.com/appcelerator/amp/tests"
	"github.com/stretchr/testify/assert"
)

var (
	h *helpers.Helper
)

func setup() (err error) {
	h, err = helpers.New()
	if err != nil {
		return err
	}
	return nil
}

func tearDown() {
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestStackDeploy(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create stack
	compose, err := ioutil.ReadFile("pinger.yml")
	assert.NoError(t, err)

	drq := &DeployRequest{
		Name:    "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16],
		Compose: compose,
	}
	drp, err := h.Stacks().Deploy(userCtx, drq)
	assert.NoError(t, err)
	assert.NotEmpty(t, drp.FullName)
	assert.NotEmpty(t, drp.Answer)

	rrq := &RemoveRequest{
		Stack: drp.FullName,
	}
	_, err = h.Stacks().Remove(userCtx, rrq)
	assert.NoError(t, err)
}

//func TestStackDeployBetweenOrganizations(t *testing.T) {
//	// Create organization with a user
//	testUser := h.RandomUser()
//	testOrg := h.RandomOrg()
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)
//	orgCtx := h.Switch(ownerCtx, t, testOrg.Name)
//
//	// Create another organization with a user
//	anotherUser := h.RandomUser()
//	anotherOrg := h.RandomOrg()
//	anotherOwnerCtx := h.CreateOrganization(t, &anotherOrg, &anotherUser)
//	anotherOrgCtx := h.Switch(anotherOwnerCtx, t, anotherOrg.Name)
//
//	// Compose file
//	compose, err := ioutil.ReadFile("pinger.yml")
//	assert.NoError(t, err)
//
//	// Deploy stack as org
//	orgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
//	rq := &DeployRequest{
//		Name:    orgStack,
//		Compose: compose,
//	}
//	r, err := h.Stacks().Deploy(orgCtx, rq)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, r.FullName)
//	assert.NotEmpty(t, r.Answer)
//
//	// Deploy another stack as another org
//	anotherOrgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
//	rq = &DeployRequest{
//		Name:    anotherOrgStack,
//		Compose: compose,
//	}
//	r, err = h.Stacks().Deploy(anotherOrgCtx, rq)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, r.FullName)
//	assert.NotEmpty(t, r.Answer)
//
//	// Update another stack as org should fail
//	rq = &DeployRequest{
//		Name:    anotherOrgStack,
//		Compose: compose,
//	}
//	r, err = h.Stacks().Deploy(orgCtx, rq)
//	assert.Error(t, err)
//
//	// Update stack as another org should fail
//	rq = &DeployRequest{
//		Name:    orgStack,
//		Compose: compose,
//	}
//	r, err = h.Stacks().Deploy(anotherOrgCtx, rq)
//	assert.Error(t, err)
//
//	// Remove stacks
//	_, err = h.Stacks().Remove(orgCtx, &RemoveRequest{Stack: orgStack})
//	assert.NoError(t, err)
//	_, err = h.Stacks().Remove(anotherOrgCtx, &RemoveRequest{Stack: anotherOrgStack})
//	assert.NoError(t, err)
//}

//func TestStackListBetweenOrganizations(t *testing.T) {
//	// Create organization with a user
//	testUser := h.RandomUser()
//	testOrg := h.RandomOrg()
//	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)
//	orgCtx := h.Switch(ownerCtx, t, testOrg.Name)
//
//	// Create another organization with a user
//	anotherUser := h.RandomUser()
//	anotherOrg := h.RandomOrg()
//	anotherOwnerCtx := h.CreateOrganization(t, &anotherOrg, &anotherUser)
//	anotherOrgCtx := h.Switch(anotherOwnerCtx, t, anotherOrg.Name)
//
//	// Compose file
//	compose, err := ioutil.ReadFile("pinger.yml")
//	assert.NoError(t, err)
//
//	// Deploy stack as org
//	orgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
//	rq := &DeployRequest{
//		Name:    orgStack,
//		Compose: compose,
//	}
//	r, err := h.Stacks().Deploy(orgCtx, rq)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, r.FullName)
//	assert.NotEmpty(t, r.Answer)
//
//	// Deploy another stack as another org
//	anotherOrgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
//	rq = &DeployRequest{
//		Name:    anotherOrgStack,
//		Compose: compose,
//	}
//	r, err = h.Stacks().Deploy(anotherOrgCtx, rq)
//	assert.NoError(t, err)
//	assert.NotEmpty(t, r.FullName)
//	assert.NotEmpty(t, r.Answer)
//
//	// Listing stacks as org
//	stacks, err := h.Stacks().List(orgCtx, &ListRequest{})
//	assert.NoError(t, err)
//	assert.Len(t, stacks.Entries, 1)
//	for _, stack := range stacks.Entries {
//		assert.NotEqual(t, stack.Stack.Name, anotherOrgStack)
//	}
//
//	// Listing stacks as another org
//	stacks, err = h.Stacks().List(anotherOrgCtx, &ListRequest{})
//	assert.NoError(t, err)
//	assert.Len(t, stacks.Entries, 1)
//	for _, stack := range stacks.Entries {
//		assert.NotEqual(t, stack.Stack.Name, orgStack)
//	}
//
//	// Remove stacks
//	_, err = h.Stacks().Remove(orgCtx, &RemoveRequest{Stack: orgStack})
//	assert.NoError(t, err)
//	_, err = h.Stacks().Remove(anotherOrgCtx, &RemoveRequest{Stack: anotherOrgStack})
//	assert.NoError(t, err)
//}

//func TestDeleteAnOrganizationOwningStacksShouldFail(t *testing.T) {
//	// Create organization with a user
//	testUser := h.RandomUser()
//	testOrg := h.DefaultOrg()
//	ownerCtx := h.CreateUser(t, &testUser)
//	orgCtx := h.Switch(ownerCtx, t, testOrg.Name)
//
//	// Compose file
//	compose, err := ioutil.ReadFile("pinger.yml")
//	assert.NoError(t, err)
//
//	// Deploy stack as org
//	orgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
//	rq := &DeployRequest{
//		Name:    orgStack,
//		Compose: compose,
//	}
//	_, err = h.Stacks().Deploy(orgCtx, rq)
//	assert.NoError(t, err)
//
//	// Deleting the organization should fail
//	_, err = h.Accounts().DeleteOrganization(orgCtx, &account.DeleteOrganizationRequest{Name: testOrg.Name})
//	assert.Error(t, err)
//
//	// Remove stack
//	_, err = h.Stacks().Remove(orgCtx, &RemoveRequest{Stack: orgStack})
//	assert.NoError(t, err)
//
//	// Deleting the organization should succeed
//	_, err = h.Accounts().DeleteOrganization(orgCtx, &account.DeleteOrganizationRequest{Name: testOrg.Name})
//	assert.NoError(t, err)
//}

func TestDeleteUserOwningStacksShouldFail(t *testing.T) {
	// Create a user
	testUser := h.RandomUser()
	userCtx := h.CreateUser(t, &testUser)

	// Compose file
	compose, err := ioutil.ReadFile("pinger.yml")
	assert.NoError(t, err)

	// Deploy stack as user
	stack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
	rq := &DeployRequest{
		Name:    stack,
		Compose: compose,
	}
	_, err = h.Stacks().Deploy(userCtx, rq)
	assert.NoError(t, err)

	// Deleting the user should fail
	_, err = h.Accounts().DeleteUser(userCtx, &account.DeleteUserRequest{Name: testUser.Name})
	assert.Error(t, err)

	// Remove stack
	_, err = h.Stacks().Remove(userCtx, &RemoveRequest{Stack: stack})
	assert.NoError(t, err)

	// Deleting the user should succeed
	_, err = h.Accounts().DeleteUser(userCtx, &account.DeleteUserRequest{Name: testUser.Name})
	assert.NoError(t, err)
}
