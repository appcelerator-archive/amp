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
