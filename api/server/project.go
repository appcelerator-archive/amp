package server

// import (
// 	"github.com/appcelerator/amp/api/rpc/project"
//
// 	"encoding/json"
// 	"golang.org/x/net/context"
// 	"log"
// )
//
// const (
// 	keySpace = "/amp/project"
// )
//
// // projectService is used to implement project.ProjectServer
// type projectService struct {
// }
//
// // CreateProject implements project.ProjectServer
// func (s *projectService) Create(ctx context.Context, in *project.CreateRequest) (*project.CreateReply, error) {
// 	// Storing the project
// 	etc.Put(keySpace, in)
//
// 	// Iterate through all entries
// 	all, err := etc.All(keySpace)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, node := range all {
// 		// Deserialize node value into a CreateRequest (could also be just a map[string]interface{}).
// 		var cr project.CreateRequest
// 		err := json.Unmarshal([]byte(node.Value), &cr)
// 		cr.Id = node.Key
// 		if err != nil {
// 			return nil, err
// 		}
// 		log.Printf("%+v\n", cr)
// 	}
//
// 	return &project.CreateReply{Message: "Hello " + in.Name}, nil
// }
