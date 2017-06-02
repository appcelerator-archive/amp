package dashboard

import (
	"log"
	"os"
	"testing"

	"github.com/appcelerator/amp/api/rpc/dashboard"
	"github.com/appcelerator/amp/data/accounts"
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

func TestDashboardCreate(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + stringid.GenerateNonCryptoID(),
		Data: "my awesome data",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)
	assert.NotEmpty(t, r)
	assert.NotEmpty(t, r.Dashboard)
	assert.NotEmpty(t, r.Dashboard.Id)
	assert.NotEmpty(t, r.Dashboard.CreateDt)
	assert.NotEmpty(t, r.Dashboard.Owner)
	assert.Equal(t, accounts.AccountType_USER, r.Dashboard.Owner.Type)
	assert.Equal(t, testUser.Name, r.Dashboard.Owner.Name)
	assert.Equal(t, rq.Name, r.Dashboard.Name)
	assert.Equal(t, rq.Data, r.Dashboard.Data)

}

func TestDashboardCreateNameAlreadyExistsshouldFail(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "my awesome data",
	}
	_, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// Create same dashboard
	_, err = h.Dashboards().Create(userCtx, rq)
	assert.Error(t, err)
}

func TestDashboardGet(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "my awesome data",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// Get dashboard
	rqGet := &dashboard.GetRequest{
		Id: r.Dashboard.Id,
	}
	rGet, err := h.Dashboards().Get(userCtx, rqGet)
	assert.NoError(t, err)
	assert.Equal(t, r.Dashboard, rGet.Dashboard)
}

func TestDashboardList(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// List dashboard
	rqList := &dashboard.ListRequest{}
	rList, err := h.Dashboards().List(userCtx, rqList)
	assert.NoError(t, err)
	assert.NotEmpty(t, rList)
	assert.NotEmpty(t, rList.Dashboards)
	assert.Contains(t, rList.Dashboards, r.Dashboard)
}

func TestDashboardUpdateName(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "my awesome data",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// Update dashboard name
	rqUpdateName := &dashboard.UpdateNameRequest{
		Id:   r.Dashboard.Id,
		Name: "updated name" + id,
	}
	_, err = h.Dashboards().UpdateName(userCtx, rqUpdateName)
	assert.NoError(t, err)

	// Get dashboard
	rqGet := &dashboard.GetRequest{
		Id: r.Dashboard.Id,
	}
	rGet, err := h.Dashboards().Get(userCtx, rqGet)
	assert.NoError(t, err)
	assert.Equal(t, rGet.Dashboard.Name, rqUpdateName.Name)
}

func TestDashboardUpdateData(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "my awesome data",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// Update dashboard name
	rqUpdateData := &dashboard.UpdateDataRequest{

		Id:   r.Dashboard.Id,
		Data: "updated data",
	}
	_, err = h.Dashboards().UpdateData(userCtx, rqUpdateData)
	assert.NoError(t, err)

	// Get dashboard
	rqGet := &dashboard.GetRequest{
		Id: r.Dashboard.Id,
	}
	rGet, err := h.Dashboards().Get(userCtx, rqGet)
	assert.NoError(t, err)
	assert.Equal(t, rGet.Dashboard.Data, rqUpdateData.Data)
}

func TestDashboardRemove(t *testing.T) {
	testUser := h.RandomUser()

	// Create user
	userCtx := h.CreateUser(t, &testUser)

	// Create dashboard
	id := stringid.GenerateNonCryptoID()
	rq := &dashboard.CreateRequest{
		Name: "my awesome dashboard" + id,
		Data: "my awesome data",
	}
	r, err := h.Dashboards().Create(userCtx, rq)
	assert.NoError(t, err)

	// Remove dashboard
	rqRemove := &dashboard.RemoveRequest{
		Id: r.Dashboard.Id,
	}
	_, err = h.Dashboards().Remove(userCtx, rqRemove)
	assert.NoError(t, err)

	// List dashboard
	rqList := &dashboard.ListRequest{}
	rList, err := h.Dashboards().List(userCtx, rqList)
	assert.NoError(t, err)
	assert.NotContains(t, rList.Dashboards, r.Dashboard)
}
