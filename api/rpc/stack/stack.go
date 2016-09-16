package stack

import (
	"path"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/pkg/stringid"
	"golang.org/x/net/context"
)

const stackRootKey = "/stacks"
const servicesRootKey = "/services"
const stackIDLabelName = "io.amp.stack.id"

// Server is used to implement stack.StackService
type Server struct {
	Store storage.Interface
}

// Up implements stack.ServerService Up
func (s *Server) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	stack, err := parseStackYaml(in.Stackfile)
	if err != nil {
		return nil, err
	}
	stackID := stringid.GenerateNonCryptoID()
	s.Store.Create(ctx, path.Join(stackRootKey, "/", stackID), stack, nil, 0)
	reply := UpReply{
		StackId: stackID,
	}
	s.Store.Delete(ctx, path.Join(stackRootKey, "/", stackID+servicesRootKey), true, nil)
	serviceIDList := make([]string, len(stack.Services), len(stack.Services))
	for i, service := range stack.Services {
		serviceID, err := s.processService(ctx, stackID, service)
		if err != nil {
			return nil, err
		}
		serviceIDList[i] = serviceID
	}
	// Save the service id list in ETCD
	val := &ServiceIdList{
		List: serviceIDList,
	}
	createErr := s.Store.Create(ctx, path.Join(stackRootKey, "/", stackID, servicesRootKey), val, nil, 0)
	if createErr != nil {
		return nil, createErr
	}
	return &reply, nil
}

// start one service and if ok store it in ETCD:
func (s *Server) processService(ctx context.Context, stackID string, serv *service.ServiceSpec) (string, error) {
	serv.Labels[stackIDLabelName] = stackID
	request := &service.ServiceCreateRequest{
		ServiceSpec: serv,
	}
	server := service.Service{}
	reply, err := server.Create(ctx, request)
	if err != nil {
		return "", err
	}
	createErr := s.Store.Create(ctx, path.Join(servicesRootKey, "/", reply.Id), serv, nil, 0)
	if createErr != nil {
		return "", createErr
	}
	return reply.Id, nil
}
