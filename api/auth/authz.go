package auth

import (
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory-am/ladon"
)

const (
	// Resources
	AmpResource          = "amprn"
	AccountResource      = AmpResource + ":account"
	OrganizationResource = AmpResource + ":organization"
	BillingResource      = AmpResource + ":billing"
	TeamResource         = AmpResource + ":team"
	RepositoryResource   = AmpResource + ":repository"
	NodeResource         = AmpResource + ":node"
	StackResource        = AmpResource + ":stack"
	BuildResource        = AmpResource + ":build"

	// Actions
	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	AdminAction  = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction
)

var (
	// Stack owners are able to administrate their own stacks
	stackOwnerAdminOwnStacks = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"}, // This will match any subject (user name), we should consider using []string{"<.+>"} to have at least one character
		Description: "Stack owners are able to administrate their own stacks",
		Resources:   []string{StackResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &ladon.EqualsSubjectCondition{},
		},
	}

	// Repository owners are able to administrate their own repositories
	repositoryOwnerAdminOwnRepository = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"}, // This will match any subject (user name), we should consider using []string{"<.+>"} to have at least one character
		Description: "Repository owners are able to administrate their own repositories",
		Resources:   []string{RepositoryResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &ladon.EqualsSubjectCondition{},
		},
	}

	// Team owners are able to administrate their own teams
	teamOwnerAdminOwnRepository = &ladon.DefaultPolicy{
		ID:          stringid.GenerateNonCryptoID(),
		Subjects:    []string{"<.*>"}, // This will match any subject (user name), we should consider using []string{"<.+>"} to have at least one character
		Description: "Team owners are able to administrate their own teams",
		Resources:   []string{TeamResource},
		Actions:     []string{"<" + AdminAction + ">"},
		Effect:      ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"owner": &ladon.EqualsSubjectCondition{},
		},
	}

	// Policies represent access control policies for amp
	Policies = []ladon.Policy{
		repositoryOwnerAdminOwnRepository,
		stackOwnerAdminOwnStacks,
		teamOwnerAdminOwnRepository,
	}
)

//package iam
//
//import (
//"github.com/ory-am/ladon"
//"github.com/stretchr/testify/assert"
//"log"
//"os"
//"testing"
//)
//
//var warden *ladon.Ladon
//
//func setup() {
//	warden = &ladon.Ladon{
//		Manager: ladon.NewMemoryManager(),
//	}
//
//	// Register all policies
//	for _, policy := range Policies {
//		if err := warden.Manager.Create(policy); err != nil {
//			log.Fatal("Unable to create policy:", err)
//		}
//	}
//}
//
//func TestMain(m *testing.M) {
//	setup()
//	retCode := m.Run()
//	os.Exit(retCode)
//}
//
//func TestStackOwnerHasAdminRights(t *testing.T) {
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   CreateAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   ReadAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   UpdateAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   DeleteAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
//
//func TestStackNotOwnerHasNoRightsByDefault(t *testing.T) {
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   CreateAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   ReadAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   UpdateAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   DeleteAction,
//		Resource: StackResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
//
//func TestRepositoryOwnerHasAdminRights(t *testing.T) {
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   CreateAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   ReadAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   UpdateAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   DeleteAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
//
//func TestRepositoryNotOwnerHasNoRightsByDefault(t *testing.T) {
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   CreateAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   ReadAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   UpdateAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   DeleteAction,
//		Resource: RepositoryResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
//
//func TestTeamOwnerHasAdminRights(t *testing.T) {
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   CreateAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   ReadAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   UpdateAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.NoError(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "john",
//		Action:   DeleteAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
//
//func TestTeamNotOwnerHasNoRightsByDefault(t *testing.T) {
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   CreateAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   ReadAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   UpdateAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//
//	assert.Error(t, warden.IsAllowed(&ladon.Request{
//		Subject:  "alice",
//		Action:   DeleteAction,
//		Resource: TeamResource,
//		Context: ladon.Context{
//			"owner": "john",
//		},
//	}))
//}
