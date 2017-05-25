package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/docker/docker/pkg/stringid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AmplifierConnection returns a grpc connection to amplifier
func AmplifierConnection() (*grpc.ClientConn, error) {
	// Connect to amplifier
	amplifierEndpoint := "amplifier" + configuration.DefaultPort
	log.Println("Connecting to amplifier")
	conn, err := grpc.Dial(amplifierEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(60*time.Second))
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to amplifier on: %store\n%v", amplifierEndpoint, err)
	}
	log.Println("Connected to amplifier")
	return conn, nil
}

// NewAccountsStore returns a new account store instance
func NewAccountsStore() accounts.Interface {
	store := etcd.New([]string{etcd.DefaultEndpoint}, "amp", 5*time.Second)
	if err := store.Connect(); err != nil {
		log.Panicf("Unable to connect to store on: %store\n%v", etcd.DefaultEndpoint, err)
	}
	return accounts.NewStore(store, configuration.RegistrationNone)
}

// ResetStorage resets the storage
func ResetStorage() {
	NewAccountsStore().Reset(context.Background())
}

// Login sign-up, verify and log-in a random user and returns the associated credentials
func Login() (metadata.MD, error) {
	id := stringid.GenerateNonCryptoID()
	randomUser := &account.SignUpRequest{
		Name:     id,
		Password: "password",
		Email:    id + "@user.amp",
	}

	conn, err := AmplifierConnection()
	if err != nil {
		return nil, err
	}

	client := account.NewAccountClient(conn)

	ctx := context.Background()

	// SignUp
	_, err = client.SignUp(ctx, randomUser)
	if err != nil {
		return nil, fmt.Errorf("SignUp error: %v", err)
	}

	// Login
	header := metadata.MD{}
	_, err = client.Login(ctx, &account.LogInRequest{Name: randomUser.Name, Password: randomUser.Password}, grpc.Header(&header))
	if err != nil {
		return nil, fmt.Errorf("Login error: %v", err)
	}

	// Extract token from header
	tokens := header[auth.TokenKey]
	if len(tokens) == 0 {
		return nil, errors.New("No token in response header")
	}
	token := tokens[0]
	return metadata.Pairs(auth.AuthorizationHeader, auth.ForgeAuthorizationHeader(token)), nil
}
