package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/appcelerator/amp/api/auth"
	"github.com/appcelerator/amp/api/rpc/account"
	"github.com/appcelerator/amp/api/rpc/stack"
)

type serverAPI struct {
	conf            *ServerConfig
	endpointConnMap map[string]*grpc.ClientConn
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
	r.HandleFunc("/api/v1/endpoints", s.endpoints).Methods("GET")
	r.HandleFunc("/api/v1/connect", s.connect).Methods("POST")
	r.HandleFunc("/api/v1/login", s.login).Methods("POST")
	r.HandleFunc("/api/v1/users", s.users).Methods("GET")
	r.HandleFunc("/api/v1/stacks", s.stacks).Methods("GET")
}

func (s *serverAPI) setToken(r *http.Request) context.Context {
	md := metadata.Pairs(auth.TokenKey, r.Header.Get("TokenKey"))
	return metadata.NewContext(context.Background(), md)
}

func (s *serverAPI) getEndpointConn(r *http.Request) (*grpc.ClientConn, error) {
	endpointName := r.Header.Get("Endpoint")
	conn, ok := s.endpointConnMap[endpointName]
	if !ok {
		return nil, fmt.Errorf("Unknown endpoint: '%s'", endpointName)
	}
	return conn, nil
}

//API functions

func (s *serverAPI) endpoints(w http.ResponseWriter, r *http.Request) {
	log.Println("execute endpoints")

	list := []string{}
	if s.conf.localEndpoint {
		list = append(list, "LocalEndpoint")
	}
	for _, ep := range s.conf.endpoints {
		list = append(list, ep)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
	//http.Error(w, "error server 1", http.StatusInternalServerError)
}

func (s *serverAPI) connect(w http.ResponseWriter, r *http.Request) {
	log.Println("execute connectEndpoint")

	//parse request data
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var t aEndpoint
	err := decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %+v\n", t)

	//connect to endpoint
	endpointName := t.Host
	if t.Local && t.Host != "localhost" {
		t.Host = "amp_amplifier"
	}
	conn, err := grpc.Dial(t.Host+":50101",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*3),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("connected to endpoint: %s host=%s\n", endpointName, t.Host)
	s.endpointConnMap[endpointName] = conn
	w.Write([]byte("done"))
}

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
	conn, errc := s.getEndpointConn(r)
	if errc != nil {
		http.Error(w, errc.Error(), http.StatusInternalServerError)
		return
	}
	client := account.NewAccountClient(conn)
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
	tokens := header[auth.TokenKey]
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

	conn, errc := s.getEndpointConn(r)
	if errc != nil {
		http.Error(w, errc.Error(), http.StatusInternalServerError)
		return
	}
	req := &account.ListUsersRequest{}
	client := account.NewAccountClient(conn)
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

	conn, errc := s.getEndpointConn(r)
	if errc != nil {
		http.Error(w, errc.Error(), http.StatusInternalServerError)
		return
	}
	req := &stack.ListRequest{}
	client := stack.NewStackClient(conn)
	reply, err := client.List(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.stack.List server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}
