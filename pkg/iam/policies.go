package iam

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
