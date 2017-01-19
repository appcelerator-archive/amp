package account

import (
	"os"
	"testing"
	"time"
	"log"

	"github.com/appcelerator/amp/config"
	"github.com/appcelerator/amp/data/storage/etcd"
	"golang.org/x/net/context"
	"strings"
	"github.com/appcelerator/amp/data/account"
	"github.com/appcelerator/amp/data/schema"
)

const (
	defTimeout = 5 * time.Second

)


var (
	acct account.Interface
	testAcct schema.Account

)

func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defTimeout)
}

func initData() {
	testAcct = schema.Account {
		Id:         "1",
		Name:       "axway",
		Type:       schema.AccountType_ORGANIZATION,
		Email:      "testowner@axway.com",
		IsVerified: false,
	}

}
func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
	log.SetPrefix("test: ")

	etcdEndpoints := []string{amp.EtcdDefaultEndpoint}
	log.Printf("connecting to etcd at %s", strings.Join(etcdEndpoints, ","))
	store := etcd.New(etcdEndpoints, "amp")
	if err := store.Connect(defTimeout); err != nil {
		log.Panicf("Unable to connect to etcd on: %s\n%v", etcdEndpoints, err)
	}
	initData()
	log.Printf("connected to etcd at %v", strings.Join(store.Endpoints(), ","))
	acct = account.NewEtcd(store, context.Background())
	os.Exit(m.Run())
}

func TestAddAccount(t *testing.T) {
	s, err := acct.AddAccount(&testAcct)
	if err != nil {
		t.Error(err)
	}
	if (s != testAcct.Id) {
		t.Errorf("expected %v, got %v", testAcct.Id, s)
	}
}

