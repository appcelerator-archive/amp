package accounts

import (
	"log"

	"github.com/appcelerator/amp/api/auth"
	"github.com/docker/docker/pkg/stringid"
	"github.com/ory/ladon"
	"golang.org/x/net/context"
)

// Resources and actions
const (
	AmpResourceName = "amprn"
	UserRN          = AmpResourceName + ":user"
	StackRN         = AmpResourceName + ":stack"
	DashboardRN     = AmpResourceName + ":dashboard"

	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
	AdminAction  = CreateAction + "|" + ReadAction + "|" + UpdateAction + "|" + DeleteAction
)

var (
	usersCanDeleteThemselves = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{UserRN},
		Actions:   []string{"<" + DeleteAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	stacksAdminByUserOwner = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{StackRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	usersAdminByThemselves = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{UserRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	dashboardsAdminByUserOwner = &ladon.DefaultPolicy{
		ID:        stringid.GenerateNonCryptoID(),
		Subjects:  []string{"<.*>"},
		Resources: []string{DashboardRN},
		Actions:   []string{"<" + AdminAction + ">"},
		Effect:    ladon.AllowAccess,
		Conditions: ladon.Conditions{
			"user": &ladon.EqualsSubjectCondition{},
		},
	}

	// Policies represent access control policies for amp
	policies = []ladon.Policy{
		usersCanDeleteThemselves,
		stacksAdminByUserOwner,
		dashboardsAdminByUserOwner,
		usersAdminByThemselves,
	}
)

// Authorization

// GetRequesterAccount gets the requester account from the given context, i.e. the user or organization performing the request
func GetRequesterAccount(ctx context.Context) *Account {
	return &Account{
		Type: AccountType_USER,
		Name: auth.GetUser(ctx),
	}
}

// IsAuthorized returns whether the requesting user is authorized to perform the given action on given resource
func (s *Store) IsAuthorized(ctx context.Context, owner *Account, action string, resource string, resourceID string) bool {
	subject := auth.GetUser(ctx)
	log.Printf("IsAuthorized: ctx(subject): %s, owner: %v, action: %s, resource: %s, resourceID: %s\n", subject, owner, action, resource, resourceID)
	if owner == nil {
		return false
	}
	err := s.warden.IsAllowed(&ladon.Request{
		Subject:  subject,
		Action:   action,
		Resource: resource,
		Context: ladon.Context{
			"user": owner.Name,
		},
	})
	return err == nil
}

// SuperUserCondition is a condition which is fulfilled if the request's subject is the super user
type SuperUserCondition struct {
}

// Fulfills returns true if subject is granted resource access
func (c *SuperUserCondition) Fulfills(value interface{}, r *ladon.Request) bool {
	return r.Subject == superUser
}

// GetName returns the condition's name.
func (c *SuperUserCondition) GetName() string {
	return "SuperUserCondition"
}
