package core

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/cli"
)

type serverAPI struct {
	endpointConn *grpc.ClientConn
}

type aSingleString struct {
	Data string `json:"data"`
}

type aLogin struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

func (s *serverAPI) handleAPIFunctions(r *mux.Router) {
	r.HandleFunc("/api/v1/connect", s.connectEndpoint).Methods("POST")
	r.HandleFunc("/api/v1/login", s.login).Methods("POST")
	r.HandleFunc("/api/v1/users", s.users).Methods("GET")
}

func (s *serverAPI) users(w http.ResponseWriter, r *http.Request) {
	request := &account.ListUsersRequest{}
	client := account.NewAccountClient(s.endpointConn)
	reply, err := client.ListUsers(context.Background(), request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) connectEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("execute connectEndpoint")

	//pase request data
	decoder := json.NewDecoder(r.Body)
	var t aSingleString
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Printf("received: %+v\n", t)
	defer r.Body.Close()

	//connect to endpoint
	conn, err := cli.NewClientConn(t.data, cli.GetToken())
	if err != nil {
		w.WriteHeader(400)
		return
	}
	s.endpointConn = conn
}

func (s *serverAPI) login(w http.ResponseWriter, r *http.Request) {

	log.Println("execute login")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aLogin
	err = decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Printf("received: %+v\n", t)

	//Execute amp login
	client := account.NewAccountClient(s.endpointConn)
	request := &account.LogInRequest{
		Name:     t.Name,
		Password: t.Pwd,
	}
	header := metadata.MD{}
	_, err = client.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	//return token
	tokens := header[auth.TokenKey]
	if len(tokens) == 0 {
		log.Println("invalid token.")
		w.WriteHeader(400)
		return
	}
	token := tokens[0]
	if token == "" {
		log.Println("invalid token.")
		w.WriteHeader(400)
		return
	}
	w.Write([]byte(token))
}
