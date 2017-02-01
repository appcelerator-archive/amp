package account

import (
	"log"
	"os"
	"testing"
	"time"

	//"github.com/appcelerator/amp/config"
	"strings"

	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/schema"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"golang.org/x/net/context"
)

const (
	defTimeout = 5 * time.Second
)

var (
	acct     account.Interface
	testAcct schema.Account
)

func newContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), defTimeout)
	return ctx
}

func initData() {
	testAcct = schema.Account{
		Id:         "",
		Name:       "axway",
		Type:       schema.AccountType_ORGANIZATION,
		Email:      "testowner@axway.com",
		IsVerified: false,
	}
	store.Delete(context.Background(), "/accounts", true, nil)

}

var store storage.Interface

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")

	etcdEndpoints := []string{amp.EtcdDefaultEndpoint}
	log.Printf("connecting to etcd at %s", strings.Join(etcdEndpoints, ","))
	store = etcd.New(etcdEndpoints, "amp-test")
	if err := store.Connect(defTimeout); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", etcdEndpoints, err)
	}
	initData()
	log.Printf("connected to etcd at %v", strings.Join(store.Endpoints(), ","))
	acct = account.NewStore(store)
	os.Exit(m.Run())
}

func TestAddAccount(t *testing.T) {
	s, err := acct.AddAccount(newContext(), &testAcct)
	if err != nil {
		t.Error(err)
	}
	if s != testAcct.Id {
		t.Errorf("expected %v, got %v", testAcct.Id, s)
	}
}
func TestListAccount(t *testing.T) {
	acct.AddAccount(newContext(), &testAcct)
	testAcct.Name = "axway2"
	testAcct.Id = ""
	acct.AddAccount(newContext(), &testAcct)
	testAcct.Id = ""
	testAcct.Name = "axway3"
	testAcct.Type = schema.AccountType_USER
	acct.AddAccount(newContext(), &testAcct)
	accList, err := acct.GetAccounts(newContext(), schema.AccountType_USER)
	if err != nil {
		log.Panicf("Unable to Fetch Account List: %v", err)
	}
	if len(accList) != 1 {
		t.Errorf("expected %v, got %v", 1, len(accList))
	}

}
