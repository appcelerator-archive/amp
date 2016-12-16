package storage

import (
	"fmt"
	"path"

	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const storageRootKey = "storage"

// Server is used to implement storage service
type Server struct {
	Store storage.Interface
}

// Put implements Server API for Storage service
func (s *Server) Put(ctx context.Context, in *PutStorage) (*StorageResponse, error) {
	response := &StorageResponse{
		Key: in.Key,
		Val: in.Val,
	}
	//save storage data in ETCD
	s.Store.Put(ctx, path.Join(storageRootKey, in.Key), response, 0)

	return response, nil
}

// Get implements Server API for Storage service
func (s *Server) Get(ctx context.Context, in *GetStorage) (*StorageResponse, error) {
	var response *StorageResponse
	response = s.getStorageByKey(ctx, in.Key)
	if response.Val == "" {
		return nil, grpc.Errorf(codes.NotFound, "storage key %s does not exist", in.Key)
	}
	return response, nil
}

// Retrieve storage key-value pair by key
func (s *Server) getStorageByKey(ctx context.Context, key string) *StorageResponse {
	response := &StorageResponse{}
	//get storage data from ETCD
	s.Store.Get(ctx, path.Join(storageRootKey, key), response, true)
	return response
}

// Delete implements Server API for Storage service
func (s *Server) Delete(ctx context.Context, in *DeleteStorage) (*StorageResponse, error) {
	var response *StorageResponse
	response = s.getStorageByKey(ctx, in.Key)
	if response.Val == "" {
		return nil, grpc.Errorf(codes.NotFound, "storage key %s does not exist", in.Key)
	}
	//delete storage data in ETCD
	s.Store.Delete(ctx, path.Join(storageRootKey, in.Key), true, nil)
	return response, nil
}

// List implements Server API for Storage service
func (s *Server) List(ctx context.Context, in *ListStorage) (*ListResponse, error) {
	var idList []proto.Message
	err := s.Store.List(ctx, storageRootKey, storage.Everything, &StorageKey{}, &idList)
	if err != nil {
		return nil, err
	}
	listInfo := []*StorageInfo{}
	for _, k := range idList {
		obj, _ := k.(*StorageKey)
		info := s.getStorageInfo(ctx, obj.Key)
		fmt.Println("info ::", info)
		listInfo = append(listInfo, s.getStorageInfo(ctx, obj.Key))
	}
	//set value for ListResponse
	response := &ListResponse{
		List: listInfo,
	}
	return response, nil
}

// return information to be displayed in storage ls
func (s *Server) getStorageInfo(ctx context.Context, key string) *StorageInfo {
	info := StorageInfo{}
	storage := StorageResponse{}
	err := s.Store.Get(ctx, path.Join(storageRootKey, key), &storage, true)
	if err == nil {
		info.Key = key
		info.Val = storage.Val
	}
	return &info
}
