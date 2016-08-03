package server

import (
	"fmt"

	"encoding/json"

	"github.com/appcelerator/amp/api/rpc/project"
	"github.com/appcelerator/amp/data/elasticsearch"

	"golang.org/x/net/context"
)

const (
	esIndex = "amp-project"
	esType  = "project"
)

var (
	es elasticsearch.ElasticSearch
)

// projectService is used to implement project.ProjectServer
type projectService struct {
}

// CreateProject implements project.ProjectServer
func (s *projectService) Create(ctx context.Context, in *project.CreateRequest) (*project.CreateReply, error) {
	// Indexing the project
	es.Index(esIndex, esType, in)

	// Get all the projects
	all := es.All(esIndex)

	// Iterate through results
	for _, hit := range all {
		// Deserialize hit.Source into a CreateRequest (could also be just a map[string]interface{}).
		var cr project.CreateRequest
		err := json.Unmarshal(*hit.Source, &cr)
		cr.Id = hit.Id
		if err != nil {
			// Deserialization failed
		}
		fmt.Printf("Project: %+v\n", cr)
	}

	return &project.CreateReply{Message: "Hello " + in.Name}, nil
}
