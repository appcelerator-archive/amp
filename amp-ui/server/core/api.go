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
	"github.com/appcelerator/amp/api/rpc/stack"
)

type serverAPI struct {
	conf *ServerConfig
	conn *grpc.ClientConn
}

type aSingleString struct {
	Data string `json:"data"`
}

type aLogin struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type aEndpoint struct {
	Host  string `json:"host"`
	Local bool   `json:"local"`
}

func (s *serverAPI) handleAPIFunctions(r *mux.Router) {
	r.HandleFunc("/api/v1/login", s.login).Methods("POST")
	r.HandleFunc("/api/v1/users", s.users).Methods("GET")
	r.HandleFunc("/api/v1/stacks", s.stacks).Methods("GET")
}

func (s *serverAPI) setToken(r *http.Request) context.Context {
	md := metadata.Pairs(auth.AuthorizationHeader, r.Header.Get("AuthorizationHeader"))
	return metadata.NewContext(context.Background(), md)
}

//API functions

func (s *serverAPI) login(w http.ResponseWriter, r *http.Request) {

	log.Println("execute login")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aLogin
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Execute amp login
	client := account.NewAccountClient(s.conn)
	request := &account.LogInRequest{
		Name:     t.Name,
		Password: t.Pwd,
	}
	header := metadata.MD{}
	_, err = client.Login(context.Background(), request, grpc.Header(&header))
	if err != nil {
		log.Println(err)
		http.Error(w, "login error", http.StatusInternalServerError)
		return
	}

	//return token
	tokens := header[auth.AuthorizationHeader]
	if len(tokens) == 0 {
		log.Println("invalid token.")
		http.Error(w, "invalid token.", http.StatusInternalServerError)
		return
	}
	token := tokens[0]
	if token == "" {
		log.Println("invalid token.")
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	ret, _ := json.Marshal(aSingleString{Data: token})
	w.Write([]byte(ret))
}

func (s *serverAPI) users(w http.ResponseWriter, r *http.Request) {
	log.Println("execute users")

	req := &account.ListUsersRequest{}
	client := account.NewAccountClient(s.conn)
	reply, err := client.ListUsers(context.Background(), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.user.ListUsers server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) stacks(w http.ResponseWriter, r *http.Request) {
	log.Println("execute stacks")

	req := &stack.ListRequest{}
	client := stack.NewStackClient(s.conn)
	reply, err := client.List(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.stack.List server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}
