package resources

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/dashboard"
	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/stack"
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

func TestResourcesListUserShouldOnlyGetHisOwnResources(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()

	// Create user and org
	userCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Deploy stack as user
	userID := stringid.GenerateNonCryptoID()[:32]
	err := h.DeployStack(userCtx, userID, "pinger.yml")
	assert.NoError(t, err)

	// Create a dashboard as user
	_, err = h.Dashboards().Create(userCtx, &dashboard.CreateRequest{Name: userID, Data: "data"})
	assert.NoError(t, err)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	orgID := stringid.GenerateNonCryptoID()[:32]
	err = h.DeployStack(orgCtx, orgID, "pinger.yml")
	assert.NoError(t, err)

	// Create a dashboard as organization
	_, err = h.Dashboards().Create(orgCtx, &dashboard.CreateRequest{Name: orgID, Data: "data"})
	assert.NoError(t, err)

	// Make sure we only get only our user resources
	reply, err := h.Resources().ListResources(userCtx, &resource.ListResourcesRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 2)

	// Make sure we only get only our organization resources
	reply, err = h.Resources().ListResources(orgCtx, &resource.ListResourcesRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 2)

	_, err = h.Stacks().Remove(userCtx, &stack.RemoveRequest{Stack: userID})
	assert.NoError(t, err)

	_, err = h.Stacks().Remove(orgCtx, &stack.RemoveRequest{Stack: orgID})
	assert.NoError(t, err)
}
