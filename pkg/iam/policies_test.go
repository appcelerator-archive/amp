package iam

import (
	"github.com/ory-am/ladon"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var warden *ladon.Ladon

func setup() {
	warden = &ladon.Ladon{
		Manager: ladon.NewMemoryManager(),
	}

	// Register all policies
	for _, policy := range Policies {
		if err := warden.Manager.Create(policy); err != nil {
			log.Fatal("Unable to create policy:", err)
		}
	}
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestStackOwnerHasAdminRights(t *testing.T) {
	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   CreateAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   ReadAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   UpdateAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   DeleteAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}

func TestStackNotOwnerHasNoRightsByDefault(t *testing.T) {
	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   CreateAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   ReadAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   UpdateAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   DeleteAction,
		Resource: StackResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}

func TestRepositoryOwnerHasAdminRights(t *testing.T) {
	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   CreateAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   ReadAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   UpdateAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   DeleteAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}

func TestRepositoryNotOwnerHasNoRightsByDefault(t *testing.T) {
	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   CreateAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   ReadAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   UpdateAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   DeleteAction,
		Resource: RepositoryResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}

func TestTeamOwnerHasAdminRights(t *testing.T) {
	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   CreateAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   ReadAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   UpdateAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.NoError(t, warden.IsAllowed(&ladon.Request{
		Subject:  "john",
		Action:   DeleteAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}

func TestTeamNotOwnerHasNoRightsByDefault(t *testing.T) {
	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   CreateAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   ReadAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   UpdateAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))

	assert.Error(t, warden.IsAllowed(&ladon.Request{
		Subject:  "alice",
		Action:   DeleteAction,
		Resource: TeamResource,
		Context: ladon.Context{
			"owner": "john",
		},
	}))
}
