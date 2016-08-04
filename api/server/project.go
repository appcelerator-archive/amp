package server

import (
	"github.com/appcelerator/amp/api/rpc/project"

	"encoding/json"
	"golang.org/x/net/context"
	"log"
)

const (
	keySpace = "/amp/project"
)

// projectService is used to implement project.ProjectServer
type projectService struct {
}

// CreateProject implements project.ProjectServer
func (s *projectService) Create(ctx context.Context, in *project.CreateRequest) (*project.CreateReply, error) {
	// Storing the project
	etc.NewKey(keySpace, in)

	// Iterate through all entries
	for _, node := range etc.All(keySpace) {
		// Deserialize node value into a CreateRequest (could also be just a map[string]interface{}).
		var cr project.CreateRequest
		err := json.Unmarshal([]byte(node.Value), &cr)
		cr.Id = node.Key
		if err != nil {
			// Deserialization failed
		}
		log.Printf("%+v\n", cr)
	}

	return &project.CreateReply{Message: "Hello " + in.Name}, nil
}
