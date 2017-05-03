package storage

import (
	"path"

	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const storageRootKey = "storage"

// Server is used to implement storage service
type Server struct {
	Store storage.Interface
}

func (s *Server) getStorageByKey(ctx context.Context, key string) *StorageEntry {
	entry := &StorageEntry{}
	s.Store.Get(ctx, path.Join(storageRootKey, key), entry, true)
	return entry
}

// Put implements Server API for Storage service
func (s *Server) Put(ctx context.Context, in *PutRequest) (*PutReply, error) {
	response := &PutReply{
		Entry: &StorageEntry{
			Key: in.Key,
			Val: in.Val,
		},
	}
	//save storage data in ETCD
	s.Store.Put(ctx, path.Join(storageRootKey, in.Key), response, 0)
	return response, nil
}

// Get implements Server API for Storage service
func (s *Server) Get(ctx context.Context, in *GetRequest) (*GetReply, error) {
	entry := s.getStorageByKey(ctx, in.Key)
	if entry.Val == "" {
		return nil, status.Errorf(codes.NotFound, "storage key %s does not exist", in.Key)
	}
	return &GetReply{Entry: entry}, nil
}

// Delete implements Server API for Storage service
func (s *Server) Delete(ctx context.Context, in *DeleteRequest) (*DeleteReply, error) {
	entry := s.getStorageByKey(ctx, in.Key)
	if entry.Val == "" {
		return nil, status.Errorf(codes.NotFound, "storage key %s does not exist", in.Key)
	}
	//delete storage data in ETCD
	s.Store.Delete(ctx, path.Join(storageRootKey, in.Key), true, nil)
	return &DeleteReply{Entry: entry}, nil
}

// List implements Server API for Storage service
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	protos := []proto.Message{}
	err := s.Store.List(ctx, storageRootKey, storage.Everything, &StorageEntry{}, &protos)
	if err != nil {
		return nil, err
	}
	entries := []*StorageEntry{}
	for _, proto := range protos {
		entries = append(entries, proto.(*StorageEntry))
	}
	return &ListReply{Entries: entries}, nil
}
