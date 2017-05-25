package resources

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/tests"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var (
	ctx context.Context
	h   *helpers.Helper
)

func setup() (err error) {
	h, err = helpers.New()
	if err != nil {
		return err
	}
	ctx = context.Background()
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
	userStack := stringid.GenerateNonCryptoID()[:32]
	err := h.DeployStack(userCtx, userStack, "pinger.yml")
	assert.NoError(t, err)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	orgStack := stringid.GenerateNonCryptoID()[:32]
	err = h.DeployStack(orgCtx, orgStack, "pinger.yml")
	assert.NoError(t, err)

	// Make sure we only get only our stack as user
	reply, err := h.Resources().ListResources(userCtx, &resource.ListResourcesRequest{})
	assert.Len(t, reply.Resources, 1)

	// Make sure we only get only our stack as organization
	reply, err = h.Resources().ListResources(orgCtx, &resource.ListResourcesRequest{})
	assert.Len(t, reply.Resources, 1)

	_, err = h.Stacks().Remove(userCtx, &stack.RemoveRequest{Stack: userStack})
	assert.NoError(t, err)

	_, err = h.Stacks().Remove(orgCtx, &stack.RemoveRequest{Stack: orgStack})
	assert.NoError(t, err)
}
