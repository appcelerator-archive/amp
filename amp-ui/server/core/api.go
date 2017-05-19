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

type aOrganization struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type aOrgName struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

type aTeamName struct {
	Organization string `json:"organization"`
	Team         string `json:"team"`
	Name         string `json:"name"`
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
	r.HandleFunc("/api/v1/organization/create", s.organizationCreate).Methods("POST")
	r.HandleFunc("/api/v1/organization/remove", s.organizationRemove).Methods("POST")
	r.HandleFunc("/api/v1/team/create", s.teamCreate).Methods("POST")
	r.HandleFunc("/api/v1/team/remove", s.teamRemove).Methods("POST")
	r.HandleFunc("/api/v1/organization/user/remove", s.removeUserFromOrganization).Methods("POST")
	r.HandleFunc("/api/v1/organization/user/add", s.addUserToOrganization).Methods("POST")
	r.HandleFunc("/api/v1/user/organizations", s.userOrganizations).Methods("POST")
	r.HandleFunc("/api/v1/team/user/remove", s.removeUserFromTeam).Methods("POST")
	r.HandleFunc("/api/v1/team/user/add", s.addUserToTeam).Methods("POST")
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

func (s *serverAPI) organizationCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("execute organizationCreate")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aOrganization
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.CreateOrganizationRequest{Name: t.Name, Email: t.Email}
	client := account.NewAccountClient(s.conn)
	reply, err := client.CreateOrganization(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.CreateOrganization server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) organizationRemove(w http.ResponseWriter, r *http.Request) {
	log.Println("execute organizationRemove")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aSingleString
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.DeleteOrganizationRequest{Name: t.Data}
	client := account.NewAccountClient(s.conn)
	reply, err := client.DeleteOrganization(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.DeleteOrganization server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) userOrganizations(w http.ResponseWriter, r *http.Request) {
	log.Println("execute userOrganizations")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aSingleString
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.GetUserOrganizationsRequest{Name: t.Data}
	client := account.NewAccountClient(s.conn)
	reply, err := client.GetUserOrganizations(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.GetUerOrganizations server error", http.StatusInternalServerError)
		return
	}
	log.Printf("return: %+v\n", reply)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) addUserToOrganization(w http.ResponseWriter, r *http.Request) {
	log.Println("execute addUserToOrganization")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aOrgName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.AddUserToOrganizationRequest{OrganizationName: t.Organization, UserName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.AddUserToOrganization(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.addUserToOrganization server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) removeUserFromOrganization(w http.ResponseWriter, r *http.Request) {
	log.Println("execute removeUserFromOrganization")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aOrgName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.RemoveUserFromOrganizationRequest{OrganizationName: t.Organization, UserName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.RemoveUserFromOrganization(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.RemoveUserFromOrganization server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) teamCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("execute teamCreate")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aOrgName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.CreateTeamRequest{OrganizationName: t.Organization, TeamName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.CreateTeam(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.CreateTeam server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) teamRemove(w http.ResponseWriter, r *http.Request) {
	log.Println("execute teamRemove")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aOrgName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.DeleteTeamRequest{OrganizationName: t.Organization, TeamName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.DeleteTeam(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.DeleteTeam server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) addUserToTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("execute addUserToTeam")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aTeamName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.AddUserToTeamRequest{OrganizationName: t.Organization, TeamName: t.Team, UserName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.AddUserToTeam(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.addUserToTeam server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func (s *serverAPI) removeUserFromTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("execute removeUserFromTeam")

	//parse request data
	var err error
	decoder := json.NewDecoder(r.Body)
	var t aTeamName
	defer r.Body.Close()
	err = decoder.Decode(&t)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("received: %v\n", t)

	req := &account.RemoveUserFromTeamRequest{OrganizationName: t.Organization, TeamName: t.Team, UserName: t.Name}
	client := account.NewAccountClient(s.conn)
	reply, err := client.RemoveUserFromTeam(s.setToken(r), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "api.account.RemoveUserFromTeam server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}
