package stack

import (
	"errors"
	"fmt"
	"path"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const stackRootKey = "stacks"
const servicesRootKey = "services"
const stackRootNameKey = "stacks/names"
const stackIDLabelName = "io.amp.stack.id"

// Server is used to implement stack.StackService
type Server struct {
	Store storage.Interface
}

// Up implements stack.ServerService Up
func (s *Server) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	stackByName := s.getStackByName(ctx, in.StackName)
	if stackByName.Id != "" {
		return nil, fmt.Errorf("Stack %s already exists", in.StackName)
	}
	stack, err := NewStackFromYaml(ctx, in.Stackfile)
	if err != nil {
		return nil, err
	}
	stack.Name = in.StackName
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Starting)); err != nil {
		return nil, err
	}
	err2 := s.Store.Create(ctx, path.Join(stackRootKey, stack.Id), stack, nil, 0)
	if err2 != nil {
		fmt.Println("error ", err2)
	}
	stackID := StackID{Id: stack.Id}
	s.Store.Create(ctx, path.Join(stackRootNameKey, stack.Name), &stackID, nil, 0)
	serviceIDList := make([]string, len(stack.Services), len(stack.Services))
	for i, service := range stack.Services {
		serviceID, err := s.processService(ctx, stack, service)
		if err != nil {
			s.rollbackStack(ctx, stack.Id, serviceIDList, err)
			return nil, err
		}
		serviceIDList[i] = serviceID
	}
	// Save the service id list in ETCD
	val := &ServiceIdList{
		List: serviceIDList,
	}
	createErr := s.Store.Create(ctx, path.Join(stackRootKey, stack.Id, servicesRootKey), val, nil, 0)
	if createErr != nil {
		s.rollbackStack(ctx, stack.Id, serviceIDList, err)
		return nil, createErr
	}
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Running)); err != nil {
		return nil, err
	}
	reply := UpReply{
		StackId: stack.Id,
	}
	fmt.Printf("Stack is running: %s\n", stack.Id)
	return &reply, nil
}

func (s *Server) getStackByName(ctx context.Context, name string) *Stack {
	stackID := &StackID{}
	s.Store.Get(ctx, path.Join(stackRootNameKey, name), stackID, true)
	stack := Stack{}
	if stackID.Id != "" {
		s.Store.Get(ctx, path.Join(stackRootKey, stackID.Id), &stack, true)
	}
	return &stack
}

func (s *Server) getStackByID(ctx context.Context, ID string) *Stack {
	stack := &Stack{}
	s.Store.Get(ctx, path.Join(stackRootKey, ID), stack, true)
	return stack
}

func (s *Server) getStack(ctx context.Context, in *StackRequest) (*Stack, error) {
	var stack *Stack
	stack = s.getStackByName(ctx, in.StackIdent)
	if stack.Id == "" {
		stack = s.getStackByID(ctx, in.StackIdent)
	}
	if stack.Id == "" {
		return nil, fmt.Errorf("The stack %s doesn't exist", in.StackIdent)
	}
	return stack, nil
}

// clean up if error happended during stack creation, delete all created services and all etcd data
func (s *Server) rollbackStack(ctx context.Context, stackID string, serviceIDList []string, err error) {
	fmt.Printf("Error found: %v \n", err)
	fmt.Printf("Cancel stack Up, cleanning up stack %s\n", stackID)
	server := service.Service{}
	for _, ID := range serviceIDList {
		if ID != "" {
			server.Remove(ctx, ID)
		}
	}
	s.Store.Delete(ctx, path.Join(stackRootKey, stackID), true, nil)
	fmt.Printf("Stack cleaned %s\n", stackID)
}

// start one service and if ok store it in ETCD:
func (s *Server) processService(ctx context.Context, stack *Stack, serv *service.ServiceSpec) (string, error) {
	serv.Labels[stackIDLabelName] = stack.Id
	serv.Name = stack.Name + "-" + serv.Name
	request := &service.ServiceCreateRequest{
		ServiceSpec: serv,
	}
	server := service.Service{}
	reply, err := server.Create(ctx, request)
	if err != nil {
		return "", err
	}
	createErr := s.Store.Create(ctx, path.Join(servicesRootKey, reply.Id), serv, nil, 0)
	if createErr != nil {
		return "", createErr
	}
	return reply.Id, nil
}

// Stop implements stack.ServerService Stop
func (s *Server) Stop(ctx context.Context, in *StackRequest) (*StackReply, error) {
	stack, errIdent := s.getStack(ctx, in)
	if errIdent != nil {
		return nil, errIdent
	}
	if running, err := stackStateMachine.Is(stack.Id, int32(StackState_Running)); err != nil {
		return nil, err
	} else if !running {
		return nil, errors.New("Stack is not running")
	}
	fmt.Printf("Stopping stack %s\n", in.StackIdent)
	server := service.Service{}
	listKeys := &ServiceIdList{}
	err := s.Store.Get(ctx, path.Join(stackRootKey, stack.Id, servicesRootKey), listKeys, true)
	if err != nil {
		return nil, err
	}
	var removeErr error
	for _, key := range listKeys.List {
		err := server.Remove(ctx, key)
		if err != nil {
			removeErr = err
		}

	}
	if removeErr != nil {
		return nil, removeErr
	}
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Stopped)); err != nil {
		return nil, err
	}
	reply := StackReply{
		StackId: stack.Id,
	}
	fmt.Printf("Stack stopped %s\n", in.StackIdent)
	return &reply, nil
}

// Remove implements stack.ServerService Remove
func (s *Server) Remove(ctx context.Context, in *StackRequest) (*StackReply, error) {
	stack, errIdent := s.getStack(ctx, in)
	if errIdent != nil {
		return nil, errIdent
	}
	if stopped, err := stackStateMachine.Is(stack.Id, int32(StackState_Stopped)); err != nil {
		return nil, err
	} else if !stopped {
		return nil, errors.New("The stack is not stopped")
	}
	fmt.Printf("Removing stack %s\n", in.StackIdent)
	s.Store.Delete(ctx, path.Join(stackRootKey, stack.Id), true, nil)
	s.Store.Delete(ctx, path.Join(stackRootNameKey, stack.Name), true, nil)
	stackStateMachine.DeleteState(stack.Id)
	reply := StackReply{
		StackId: stack.Id,
	}
	fmt.Printf("Stack removed %s\n", in.StackIdent)
	return &reply, nil
}

// List list all available stack with there status
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	var idList []proto.Message
	err := s.Store.List(ctx, stackRootNameKey, storage.Everything, &StackID{}, &idList)
	if err != nil {
		return nil, err
	}
	listInfo := make([]*StackInfo, len(idList), len(idList))
	for i, ID := range idList {
		obj, _ := ID.(*StackID)
		listInfo[i] = s.getStackInfo(ctx, obj.Id)
	}
	reply := ListReply{
		List: listInfo,
	}
	return &reply, nil
}

func (s *Server) getStackInfo(ctx context.Context, ID string) *StackInfo {
	info := StackInfo{}
	stack := Stack{}
	err := s.Store.Get(ctx, path.Join(stackRootKey, ID), &stack, true)
	if err == nil {
		info.Name = stack.Name
		info.Id = stack.Id
	}
	/* Waiting for Bertrand's getState
	state := &State{}
	errGet := s.Store.Get(ctx, path.Join(stackRootKey, ID, stackStackKey), state, true)
	info.State = "nc"
	if errGet == nil {
		switch state.Value {
		case 0:
			info.State = "Stopped"
		case 1:
			info.State = "Starting"
		case 2:
			info.State = "Running"
		case 3:
			info.State = "Redeploying"
		}
	}
	*/
	return &info
}
