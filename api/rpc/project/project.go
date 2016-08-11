package project

// import (
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
// // Service is used to implement ProjectServer
// type Service struct {
// }
//
// // CreateProject implements ProjectServer
// func (s *Service) Create(ctx context.Context, in *CreateRequest) (*CreateReply, error) {
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
// 	return &CreateReply{Message: "Hello " + in.Name}, nil
// }
