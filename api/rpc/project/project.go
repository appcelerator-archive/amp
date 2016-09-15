package project

import (
	"context"
	"fmt"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
)

const (
	prefix = "project"
)

// Proj structure to implement ProjServer interface
type Proj struct {
	Store storage.Interface
}

// Create adds a new entry to the k,v data store
func (p *Proj) Create(ctx context.Context, req *CreateRequest) (*CreateReply, error) {
	key := fmt.Sprintf("%s/%v", prefix, req.Project.Id)
	ttl := int64(0)
	reply := CreateReply{Project: &ProjectEntry{}}
	err := p.Store.Create(ctx, key, req.Project, reply.Project, ttl)
	return &reply, err
}

//Update removes the entry for the specified Key
func (p *Proj) Update(ctx context.Context, req *UpdateRequest) (*UpdateReply, error) {
	// Build The Key for the Project
	key := fmt.Sprintf("%s/%v", prefix, req.Project.Id)
	ttl := int64(0)
	reply := UpdateReply{Project: &ProjectEntry{}}
	// Update the Value for the given Key
	err := p.Store.Update(ctx, key, req.Project, ttl)
	return &reply, err
}

//Get fetches the entry for the specified Key
func (p *Proj) Get(ctx context.Context, req *GetRequest) (*GetReply, error) {
	// Build The Key for the Project
	key := fmt.Sprintf("%s/%v", prefix, req.Id)
	reply := GetReply{Project: &ProjectEntry{}}
	// Retrieve the Value for the Given Key
	err := p.Store.Get(ctx, key, reply.Project, true)
	return &reply, err
}

//Delete removes the entry for the specified Key
func (p *Proj) Delete(ctx context.Context, req *DeleteRequest) (*DeleteReply, error) {
	// Build The Key for the Project
	key := fmt.Sprintf("%s/%v", prefix, req.Id)
	reply := DeleteReply{Project: &ProjectEntry{}}
	// Delete the Key and Value from the Store
	err := p.Store.Delete(ctx, key, reply.Project)
	return &reply, err
}

//List retrieves all of the values for a specified Key range
func (p *Proj) List(ctx context.Context, req *ListRequest) (*ListReply, error) {
	// Build the Key based on the prefix only
	key := fmt.Sprintf("%s", prefix)
	obj := ProjectEntry{}
	var out []proto.Message
	reply := ListReply{}
	// Return all of the entries for the given Key Pattern
	err := p.Store.List(ctx, key, storage.Everything, &obj, &out)
	if err != nil {
		return &reply, err
	}
	reply.Projects = make([]*ProjectEntry, len(out))
	for i, project := range out {
		reply.Projects[i] = project.(*ProjectEntry)
	}

	//// TODO Could not avoid the wire type error without explicit type conversion
	//// 	reply.Projects = out.([]*ProjectEntry) <-- Illegal Syntax
	//
	//for i := range out {
	//	p, ok := out[i].(*ProjectEntry)
	//	if !ok {
	//		return reply, fmt.Errorf("Invalid Type Conversion")
	//	}
	//	reply.Projects = append(reply.Projects, p)
	//}
	return &reply, err
}
