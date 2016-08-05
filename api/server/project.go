package server

import (
	"github.com/appcelerator/amp/api/rpc/project"

	"encoding/json"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

const (
	keySpace = "/amp/project"
)

// projectService is used to implement project.ProjectServer
type projectService struct {
}

// CreateProject implements project.ProjectServer
func (s *projectService) Create(ctx context.Context, in *project.CreateRequest) (*project.CreateReply, error) {
	// Generate an id
	id := uuid.NewV4().String()
	in.Project.Id = id

	// Store the project
	_, err := etc.Put(keySpace, id, in.Project)
	if err != nil {
		return nil, err
	}

	return &project.CreateReply{Created: in.Project}, nil
}

func (s *projectService) List(ctx context.Context, in *project.ListRequest) (*project.ListReply, error) {
	// Get all the projects
	all, err := etc.List(keySpace)
	if err != nil {
		return nil, err
	}

	// Create the reply
	projects := make([]*project.Project, len(all))
	for i, node := range all {
		var project project.Project
		err := json.Unmarshal([]byte(node.Value), &project)
		if err != nil {
			return nil, err
		}
		projects[i] = &project
	}

	return &project.ListReply{Projects: projects}, nil
}
