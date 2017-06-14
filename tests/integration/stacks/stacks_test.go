package stack

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	. "github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/tests"
	"github.com/docker/docker/pkg/stringid"
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

func TestStackDeployBetweenOrganizations(t *testing.T) {
	// Create organization with a user
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	ownerCtx := h.CreateOrganization(t, &testOrg, &testUser)
	orgCtx := h.Switch(ownerCtx, t, testOrg.Name)

	// Create another organization with a user
	anotherUser := h.RandomUser()
	anotherOrg := h.RandomOrg()
	anotherOwnerCtx := h.CreateOrganization(t, &anotherOrg, &anotherUser)
	anotherOrgCtx := h.Switch(anotherOwnerCtx, t, anotherOrg.Name)

	// Compose file
	compose, err := ioutil.ReadFile("pinger.yml")
	assert.NoError(t, err)

	// Deploy stack as org
	orgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
	rq := &DeployRequest{
		Name:    orgStack,
		Compose: compose,
	}
	r, err := h.Stacks().Deploy(orgCtx, rq)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.FullName)
	assert.NotEmpty(t, r.Answer)

	// Deploy another stack as another org
	anotherOrgStack := "my-awesome-stack" + stringid.GenerateNonCryptoID()[:16]
	rq = &DeployRequest{
		Name:    anotherOrgStack,
		Compose: compose,
	}
	r, err = h.Stacks().Deploy(anotherOrgCtx, rq)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.FullName)
	assert.NotEmpty(t, r.Answer)

	// Update another stack as org should fail
	rq = &DeployRequest{
		Name:    anotherOrgStack,
		Compose: compose,
	}
	r, err = h.Stacks().Deploy(orgCtx, rq)
	assert.Error(t, err)

	// Update stack as another org should fail
	rq = &DeployRequest{
		Name:    orgStack,
		Compose: compose,
	}
	r, err = h.Stacks().Deploy(anotherOrgCtx, rq)
	assert.Error(t, err)

	// Remove stacks
	_, err = h.Stacks().Remove(orgCtx, &RemoveRequest{Stack: orgStack})
	assert.NoError(t, err)
	_, err = h.Stacks().Remove(anotherOrgCtx, &RemoveRequest{Stack: anotherOrgStack})
	assert.NoError(t, err)
}
