package resources

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/account"
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

func TestListOrganizationShouldOnlyGetItsOwnResources(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	anotherUser := h.RandomUser()
	anotherOrg := h.RandomOrg()

	// Create user and org
	userCtx := h.CreateOrganization(t, &testOrg, &testUser)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// Create a dashboard as organization
	_, err = h.Dashboards().Create(orgCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "data"})
	assert.NoError(t, err)

	// Create another user and another org
	anotherUserCtx := h.CreateOrganization(t, &anotherOrg, &anotherUser)

	// Switch to another organization account
	anotherOrgCtx := h.Switch(anotherUserCtx, t, anotherOrg.Name)

	// Deploy stack as another organization
	anotherStackID, err := h.DeployStack(anotherOrgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// Create a dashboard as another organization
	_, err = h.Dashboards().Create(anotherOrgCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "another data"})
	assert.NoError(t, err)

	// Make sure we only get only our organization resources
	reply, err := h.Resources().List(orgCtx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 2)

	// Make sure we only get only our another organization resources
	reply, err = h.Resources().List(anotherOrgCtx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 2)

	_, err = h.Stacks().Remove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
	_, err = h.Stacks().Remove(anotherOrgCtx, &stack.RemoveRequest{Stack: anotherStackID})
	assert.NoError(t, err)
}

func TestAddSameResourceTwiceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user, org and team
	userCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// AddToTeam
	_, err = h.Resources().AddToTeam(orgCtx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// AddToTeam again
	_, err = h.Resources().AddToTeam(orgCtx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.Error(t, err)

	_, err = h.Stacks().Remove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestRemoveNonExistingResourceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user, org and team
	userCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// RemoveFromTeam
	_, err = h.Resources().RemoveFromTeam(orgCtx, &resource.RemoveFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.Error(t, err)

	_, err = h.Stacks().Remove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestAuthorizations(t *testing.T) {
	testUser := h.RandomUser()
	testMember := h.RandomUser()
	testOtherMember := h.RandomUser()
	testOrg := h.RandomOrg()
	testTeam := h.RandomTeam(testOrg.Name)
	su, err := h.SuperLogin()
	assert.NoError(t, err)

	/// Create team
	ownerCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Create member in org
	memberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testMember)
	otherMemberCtx := h.CreateAndAddUserToOrganization(ownerCtx, t, &testOrg, &testOtherMember)

	// AddUserToTeam
	_, err = h.Accounts().AddUserToTeam(ownerCtx, &account.AddUserToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		UserName:         testMember.Name,
	})
	assert.NoError(t, err)

	// User vs themselves

	// user can read themselves
	reply, err := h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// user can update themselves
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// user can delete themselves
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// User vs others

	// user cannot read others
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// user cannot update others
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// user cannot delete others
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read others
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update others
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete others
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testUser.Name,
			Type:   resource.ResourceType_RESOURCE_USER,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Organizations owners

	// owner can read organization
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update organization
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete organization
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can create team
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_TEAM,
			Action: resource.Action_ACTION_CREATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Organizations members

	// member cannot read organization
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// member cannot update organization
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// member cannot delete organization
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// member cannot create team
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_TEAM,
			Action: resource.Action_ACTION_CREATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read organization
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update organization
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete organization
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_ORGANIZATION,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can create team
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     testOrg.Name,
			Type:   resource.ResourceType_RESOURCE_TEAM,
			Action: resource.Action_ACTION_CREATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Stacks

	// Deploy stack as user
	userStackID, err := h.DeployStack(ownerCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// Owners

	// owner can read stack
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update stack
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete stack
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Create a dashboard as user
	r, err := h.Dashboards().Create(ownerCtx, &dashboard.CreateRequest{Name: userStackID, Data: "data"})
	assert.NoError(t, err)
	userDashboardId := r.Dashboard.Id

	// Owners

	// owner can read dashboard
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update dashboard
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete dashboard
	reply, err = h.Resources().Authorizations(ownerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     userDashboardId,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	_, err = h.Stacks().Remove(ownerCtx, &stack.RemoveRequest{Stack: userStackID})
	assert.NoError(t, err)

	// Deploy stack as organization owner
	orgOwnerCtx := h.Switch(ownerCtx, t, testOrg.Name)
	orgOwnerStackID, err := h.DeployStack(orgOwnerCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// Owners

	// owner can read stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete stack
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	_, err = h.Stacks().Remove(orgOwnerCtx, &stack.RemoveRequest{Stack: orgOwnerStackID})
	assert.NoError(t, err)

	// Deploy stack as organization owner
	orgMemberCtx := h.Switch(memberCtx, t, testOrg.Name)
	orgMemberStackID, err := h.DeployStack(orgMemberCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// Owners

	// owner can read stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete stack
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read stack
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update stack
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete stack
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete stack
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberStackID,
			Type:   resource.ResourceType_RESOURCE_STACK,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	_, err = h.Stacks().Remove(su, &stack.RemoveRequest{Stack: orgMemberStackID})
	assert.NoError(t, err)

	// Deploy dashboard as organization owner
	r, err = h.Dashboards().Create(orgOwnerCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "data"})
	assert.NoError(t, err)
	orgOwnerDashboardID := r.Dashboard.Id

	// Owners

	// owner can read dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete dashboard
	reply, err = h.Resources().Authorizations(memberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgOwnerDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Deploy dashboard as organization owner
	r, err = h.Dashboards().Create(orgMemberCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "data"})
	assert.NoError(t, err)
	orgMemberDashboardID := r.Dashboard.Id

	// Owners

	// owner can read dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can update dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// owner can delete dashboard
	reply, err = h.Resources().Authorizations(orgOwnerCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// Others

	// others cannot read dashboard
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot update dashboard
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// others cannot delete dashboard
	reply, err = h.Resources().Authorizations(otherMemberCtx, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.False(t, reply.Replies[0].Authorized)

	// SuperUser

	// su can read dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_READ,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can update dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_UPDATE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)

	// su can delete dashboard
	reply, err = h.Resources().Authorizations(su, &resource.AuthorizationsRequest{
		Requests: []*resource.IsAuthorizedRequest{{
			Id:     orgMemberDashboardID,
			Type:   resource.ResourceType_RESOURCE_DASHBOARD,
			Action: resource.Action_ACTION_DELETE,
		}},
	})
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.Len(t, reply.Replies, 1)
	assert.True(t, reply.Replies[0].Authorized)
}
