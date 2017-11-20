package resources

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/resource"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/appcelerator/amp/data/accounts"
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

//func TestListResourcesShouldReturnOnlyActiveOrganizationResources(t *testing.T) {
//	testUser := h.RandomUser()
//	testOrg := h.RandomOrg()
//	anotherUser := h.RandomUser()
//	anotherOrg := h.RandomOrg()
//
//	// Create user and org
//	userCtx := h.CreateOrganization(t, &testOrg, &testUser)
//
//	// Switch to organization account
//	orgCtx := h.Switch(userCtx, t, testOrg.Name)
//
//	// Deploy stack as organization
//	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
//	assert.NoError(t, err)
//
//	// Create a dashboard as organization
//	_, err = h.Dashboards().Create(orgCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "data"})
//	assert.NoError(t, err)
//
//	// Create another user and another org
//	anotherUserCtx := h.CreateOrganization(t, &anotherOrg, &anotherUser)
//
//	// Switch to another organization account
//	anotherOrgCtx := h.Switch(anotherUserCtx, t, anotherOrg.Name)
//
//	// Deploy stack as another organization
//	anotherStackID, err := h.DeployStack(anotherOrgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
//	assert.NoError(t, err)
//
//	// Create a dashboard as another organization
//	_, err = h.Dashboards().Create(anotherOrgCtx, &dashboard.CreateRequest{Name: stringid.GenerateNonCryptoID()[:32], Data: "another data"})
//	assert.NoError(t, err)
//
//	// Make sure we only get only our organization resources
//	reply, err := h.Resources().ResourceList(orgCtx, &resource.ListRequest{})
//	assert.NoError(t, err)
//	assert.Len(t, reply.Resources, 2)
//
//	// Make sure we only get only our another organization resources
//	reply, err = h.Resources().ResourceList(anotherOrgCtx, &resource.ListRequest{})
//	assert.NoError(t, err)
//	assert.Len(t, reply.Resources, 2)
//
//	_, err = h.Stacks().StackRemove(orgCtx, &stack.RemoveRequest{Stack: stackID})
//	assert.NoError(t, err)
//	_, err = h.Stacks().StackRemove(anotherOrgCtx, &stack.RemoveRequest{Stack: anotherStackID})
//	assert.NoError(t, err)
//}

func TestAddSameResourceTwiceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user, org and team
	userCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// AddToTeam
	_, err = h.Resources().ResourceAddToTeam(orgCtx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// AddToTeam again
	_, err = h.Resources().ResourceAddToTeam(orgCtx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.Error(t, err)

	_, err = h.Stacks().StackRemove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestRemoveNonExistingResourceShouldFail(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user, org and team
	userCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// RemoveFromTeam
	_, err = h.Resources().ResourceRemoveFromTeam(orgCtx, &resource.RemoveFromTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.Error(t, err)

	_, err = h.Stacks().StackRemove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestDeletedStackShouldNotBelongToTheTeamAnymore(t *testing.T) {
	testUser := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user, org and team
	userCtx := h.CreateTeam(t, &testOrg, &testUser, &testTeam)

	// Switch to organization account
	orgCtx := h.Switch(userCtx, t, testOrg.Name)

	// Deploy stack as organization
	stackID, err := h.DeployStack(orgCtx, stringid.GenerateNonCryptoID()[:32], "pinger.yml")
	assert.NoError(t, err)

	// AddToTeam
	_, err = h.Resources().ResourceAddToTeam(orgCtx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// GetTeam
	reply, err := h.Accounts().GetTeam(orgCtx, &account.GetTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, reply.Team.Resources)

	// Delete the stack
	_, err = h.Stacks().StackRemove(orgCtx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)

	// GetTeam
	reply, err = h.Accounts().GetTeam(orgCtx, &account.GetTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
	})
	assert.NoError(t, err)
	assert.Empty(t, reply.Team.Resources)
}

func TestShareResourceReadPermission(t *testing.T) {
	testUser := h.RandomUser()
	testMember1 := h.RandomUser()
	testMember2 := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user and org
	h.CreateUser(t, &testUser)

	// Add members
	member1Ctx := h.CreateUser(t, &testMember1)
	member2Ctx := h.CreateUser(t, &testMember2)

	// Switch to organization account
	member1Ctx = h.Switch(member1Ctx, t, testOrg.Name)
	member2Ctx = h.Switch(member2Ctx, t, testOrg.Name)

	// Member 1 create a team
	_, err := h.Accounts().CreateTeam(member1Ctx, &testTeam)
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err := h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// Deploy stack as organization
	stackName := stringid.GenerateNonCryptoID()[:32]
	stackID, err := h.DeployStack(member1Ctx, stackName, "pinger.yml")
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// AddToTeam
	_, err = h.Resources().ResourceAddToTeam(member1Ctx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// Add member 2 to the team
	h.AddUserToTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)

	// Member 2 should not be able to update the stack
	_, err = h.DeployStack(member2Ctx, stackName, "pinger.yml")
	assert.Error(t, err)

	// Member 2 should not be able to remove the stack
	_, err = h.Stacks().StackRemove(member2Ctx, &stack.RemoveRequest{Stack: stackID})
	assert.Error(t, err)

	// Remove member 2 from the team
	h.RemoveUserFromTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// Remove stack
	_, err = h.Stacks().StackRemove(member1Ctx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestShareResourceWritePermission(t *testing.T) {
	testUser := h.RandomUser()
	testMember1 := h.RandomUser()
	testMember2 := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user and org
	h.CreateUser(t, &testUser)

	// Add members
	member1Ctx := h.CreateUser(t, &testMember1)
	member2Ctx := h.CreateUser(t, &testMember2)

	// Switch to organization account
	member1Ctx = h.Switch(member1Ctx, t, testOrg.Name)
	member2Ctx = h.Switch(member2Ctx, t, testOrg.Name)

	// Member 1 create a team
	_, err := h.Accounts().CreateTeam(member1Ctx, &testTeam)
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err := h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// Deploy stack as organization
	stackName := stringid.GenerateNonCryptoID()[:32]
	stackID, err := h.DeployStack(member1Ctx, stackName, "pinger.yml")
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// AddToTeam
	_, err = h.Resources().ResourceAddToTeam(member1Ctx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// ChangePermissionLevel
	_, err = h.Resources().ResourceChangePermissionLevel(member1Ctx, &resource.ChangePermissionLevelRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
		PermissionLevel:  accounts.TeamPermissionLevel_TEAM_WRITE,
	})
	assert.NoError(t, err)

	// Add member 2 to the team
	h.AddUserToTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)

	// Member 2 should be able to update the stack
	_, err = h.DeployStack(member2Ctx, stackName, "pinger.yml")
	assert.NoError(t, err)

	// Member 2 should not be able to remove the stack
	_, err = h.Stacks().StackRemove(member2Ctx, &stack.RemoveRequest{Stack: stackID})
	assert.Error(t, err)

	// Remove member 2 from the team
	h.RemoveUserFromTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// Remove stack
	_, err = h.Stacks().StackRemove(member1Ctx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)
}

func TestShareResourceAdminPermission(t *testing.T) {
	testUser := h.RandomUser()
	testMember1 := h.RandomUser()
	testMember2 := h.RandomUser()
	testOrg := h.DefaultOrg()
	testTeam := h.RandomTeam(testOrg.Name)

	// Create user and org
	h.CreateUser(t, &testUser)

	// Add members
	member1Ctx := h.CreateUser(t, &testMember1)
	member2Ctx := h.CreateUser(t, &testMember2)

	// Switch to organization account
	member1Ctx = h.Switch(member1Ctx, t, testOrg.Name)
	member2Ctx = h.Switch(member2Ctx, t, testOrg.Name)

	// Member 1 create a team
	_, err := h.Accounts().CreateTeam(member1Ctx, &testTeam)
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err := h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// Deploy stack as organization
	stackName := stringid.GenerateNonCryptoID()[:32]
	stackID, err := h.DeployStack(member1Ctx, stackName, "pinger.yml")
	assert.NoError(t, err)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)

	// AddToTeam
	_, err = h.Resources().ResourceAddToTeam(member1Ctx, &resource.AddToTeamRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
	})
	assert.NoError(t, err)

	// ChangePermissionLevel
	_, err = h.Resources().ResourceChangePermissionLevel(member1Ctx, &resource.ChangePermissionLevelRequest{
		OrganizationName: testTeam.OrganizationName,
		TeamName:         testTeam.TeamName,
		ResourceId:       stackID,
		PermissionLevel:  accounts.TeamPermissionLevel_TEAM_ADMIN,
	})
	assert.NoError(t, err)

	// Add member 2 to the team
	h.AddUserToTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Len(t, reply.Resources, 1)

	// Member 2 should be able to update the stack
	_, err = h.DeployStack(member2Ctx, stackName, "pinger.yml")
	assert.NoError(t, err)

	// Member 2 should be able to remove the stack
	_, err = h.Stacks().StackRemove(member2Ctx, &stack.RemoveRequest{Stack: stackID})
	assert.NoError(t, err)

	// Remove member 2 from the team
	h.RemoveUserFromTeam(member1Ctx, t, &testTeam, &testMember2)

	// Make sure we can list only our resources
	reply, err = h.Resources().ResourceList(member1Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)
	reply, err = h.Resources().ResourceList(member2Ctx, &resource.ListRequest{})
	assert.NoError(t, err)
	assert.Empty(t, reply.Resources)
}
