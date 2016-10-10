package stack

import (
	"errors"
	"fmt"
	"path"

	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const stackRootKey = "stacks"
const servicesRootKey = "services"
const networksRootKey = "networks"
const stackRootNameKey = "stacks/names"
const stackIDLabelName = "io.amp.stack.id"
const stackNameLabelName = "io.amp.stack.name"
const serviceRoleLabelName = "io.amp.role"

// Server is used to implement stack.StackService
type Server struct {
	Store  storage.Interface
	Docker *client.Client
}

// Up implements stack.ServerService Up
func (s *Server) Up(ctx context.Context, in *UpRequest) (*UpReply, error) {
	stackByName := s.getStackByName(ctx, in.StackName)
	if stackByName.Id != "" {
		return nil, fmt.Errorf("Stack %s already exists", in.StackName)
	}
	stack, err := newStackFromYaml(ctx, in.Stackfile)
	if err != nil {
		return nil, err
	}
	stack.Name = in.StackName
	stack.IsPublic = s.isPublic(stack)
	errCreate := s.Store.Create(ctx, path.Join(stackRootKey, stack.Id), stack, nil, 0)
	if errCreate != nil {
		return nil, errCreate
	}
	stackID := StackID{Id: stack.Id}
	s.Store.Create(ctx, path.Join(stackRootNameKey, stack.Name), &stackID, nil, 0)
	startRequest := StackRequest{
		StackIdent: stack.Id,
	}
	_, errStart := s.Start(ctx, &startRequest)
	if errStart != nil {
		fmt.Printf("Error found during stack up: %v \n", errStart)
		s.rollbackETCDStack(ctx, stack)
		return nil, errStart
	}
	fmt.Printf("Stack is up: %s\n", stack.Id)
	reply := UpReply{
		StackId: stack.Id,
	}
	return &reply, nil
}

// determine if stack have at least one service having one public name
func (s *Server) isPublic(stack *Stack) bool {
	isPublic := false
	for _, service := range stack.Services {
		if service.PublishSpecs != nil {
			for _, public := range service.PublishSpecs {
				if public.Name != "" {
					isPublic = true
				}
			}
		}
	}
	return isPublic
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
func (s *Server) rollbackStack(ctx context.Context, stackID string, serviceIDList []string) {
	fmt.Printf("removing created networks %s\n", stackID)
	s.removeStackNetworks(ctx, stackID, true)
	fmt.Printf("removing created services %s\n", stackID)
	server := service.Service{
		Docker: s.Docker,
	}
	for _, ID := range serviceIDList {
		if ID != "" {
			server.Remove(ctx, &service.RemoveRequest{
				Ident: ID,
			})
			s.Store.Delete(ctx, path.Join(servicesRootKey, ID), true, nil)
		}
	}
	fmt.Printf("Services removed %s\n", stackID)
}

// clean up if error happended during stack creation, delete all created services and all etcd data
func (s *Server) rollbackETCDStack(ctx context.Context, stack *Stack) {
	fmt.Printf("Cleanning up ETCD storage %s\n", stack.Id)
	s.Store.Delete(ctx, path.Join(stackRootKey, stack.Id), true, nil)
	s.Store.Delete(ctx, path.Join(stackRootNameKey, stack.Name), true, nil)
	fmt.Printf("ETCD cleaned %s\n", stack.Id)
}

// start one service and if ok store it in ETCD:
func (s *Server) processService(ctx context.Context, stack *Stack, serv *service.ServiceSpec) (string, error) {
	//Add common labels to services
	if serv.Labels == nil {
		serv.Labels = make(map[string]string)
	}
	serv.Labels[stackIDLabelName] = stack.Id
	serv.Labels[stackNameLabelName] = stack.Name
	if serv.ContainerLabels == nil {
		serv.ContainerLabels = make(map[string]string)
	}
	serv.ContainerLabels[stackIDLabelName] = stack.Id
	serv.ContainerLabels[stackNameLabelName] = stack.Name
	// add default network
	if serv.Name != "haproxy" {
		if serv.Networks == nil {
			serv.Networks = []*service.NetworkAttachment{}
		}
		serv.Networks = append(serv.Networks, &service.NetworkAttachment{
			Target:  stack.Name,
			Aliases: []string{serv.Name},
		})
		isPublic := false
		if serv.PublishSpecs != nil {
			for _, public := range serv.PublishSpecs {
				if public.Name != "" {
					isPublic = true
				}
			}
		}
		if isPublic {
			serv.Networks = append(serv.Networks, &service.NetworkAttachment{
				Target:  fmt.Sprintf("%s-public", stack.Name),
				Aliases: []string{serv.Name},
			})
		}
	}
	//update name
	serv.Name = fmt.Sprintf("%s-%s", stack.Name, serv.Name)
	//Create service
	request := &service.ServiceCreateRequest{
		ServiceSpec: serv,
	}
	server := service.Service{
		Docker: s.Docker,
	}
	reply, err := server.Create(ctx, request)
	if err != nil {
		return "", err
	}
	//Save service defintion in ETCD
	createErr := s.Store.Create(ctx, path.Join(servicesRootKey, reply.Id), serv, nil, 0)
	if createErr != nil {
		return "", createErr
	}
	return reply.Id, nil
}

// add HAProxy service dedicated to the stack reverse proxy
func (s *Server) addHAProxyService(ctx context.Context, stack *Stack) (string, error) {
	serv := service.ServiceSpec{
		Image: "appcelerator/haproxy:1.0.1",
		Name:  "haproxy",
		Env:   []string{"STACKNAME=" + stack.Name},
		Networks: []*service.NetworkAttachment{
			{
				Target:  "amp-public",
				Aliases: []string{fmt.Sprintf("%s-haproxy", stack.Name)},
			},
			{
				Target:  fmt.Sprintf("%s-public", stack.Name),
				Aliases: []string{fmt.Sprintf("%s-haproxy", stack.Name)},
			},
			{
				Target:  "amp-infra",
				Aliases: []string{fmt.Sprintf("%s-haproxy", stack.Name)},
			},
		},
	}
	//Verify if there is an HAProxy service in the stack definition to update publish port if exist
	var publishPort uint32 = 0
	for _, service := range stack.Services {
		if service.Name == "haproxy" {
			for _, public := range service.PublishSpecs {
				if public.InternalPort == 80 {
					publishPort = public.PublishPort
				}
			}
		}
	}
	if publishPort != 0 {
		serv.PublishSpecs = []*service.PublishSpec{
			{
				InternalPort: 80,
				PublishPort:  publishPort,
			},
		}
	}
	return s.processService(ctx, stack, &serv)
}

// create network
func (s *Server) createStackNetwork(ctx context.Context, name string) (string, error) {
	fmt.Printf("Create network %s\n", name)
	rep, err := s.Docker.NetworkCreate(ctx, name, types.NetworkCreate{
		CheckDuplicate: false,
		Driver:         "overlay",
		EnableIPv6:     false,
		Internal:       true,
		Labels:         map[string]string{serviceRoleLabelName: "user"},
		IPAM: &network.IPAM{
			Driver: "default",
		},
	})
	if err != nil {
		return "", err
	}
	return rep.ID, nil
}

func (s *Server) createStackNetworks(ctx context.Context, stack *Stack) error {
	networkList := []string{}
	id, err := s.createStackNetwork(ctx, stack.Name)
	if err != nil {
		return err
	}
	networkList = append(networkList, id)
	if stack.IsPublic {
		id, err := s.createStackNetwork(ctx, stack.Name+"-public")
		if err != nil {
			s.removeStackNetworksFromList(ctx, networkList)
			return err
		}
		networkList = append(networkList, id)
	}
	//TODO: add custom network
	list := IdList{
		List: networkList,
	}
	if uerr := s.Store.Update(ctx, path.Join(stackRootKey, stack.Id, networksRootKey), &list, 0); uerr != nil {
		if cerr := s.Store.Create(ctx, path.Join(stackRootKey, stack.Id, networksRootKey), &list, nil, 0); cerr != nil {
			s.removeStackNetworksFromList(ctx, networkList)
			return cerr
		}
	}
	return nil
}

// Start implements stack.ServerService Stop
func (s *Server) Start(ctx context.Context, in *StackRequest) (*StackReply, error) {
	stack, errIdent := s.getStack(ctx, in)
	if errIdent != nil {
		return nil, errIdent
	}
	if stack.Services == nil || len(stack.Services) == 0 {
		return nil, fmt.Errorf("No services found for the stack %s \n", in.StackIdent)
	}
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Starting)); err != nil {
		return nil, err
	}
	fmt.Printf("Starting stack %s\n", in.StackIdent)
	err := s.createStackNetworks(ctx, stack)
	if err != nil {
		return nil, err
	}

	serviceIDList := []string{}
	if stack.IsPublic {
		serviceID, err := s.addHAProxyService(ctx, stack)
		if err != nil {
			s.rollbackStack(ctx, stack.Id, serviceIDList)
			return nil, err
		}
		serviceIDList = append(serviceIDList, serviceID)
	}
	for _, service := range stack.Services {
		if service.Name != "haproxy" {
			serviceID, err := s.processService(ctx, stack, service)
			if err != nil {
				s.rollbackStack(ctx, stack.Id, serviceIDList)
				return nil, err
			}
			serviceIDList = append(serviceIDList, serviceID)
		}
	}
	// Save the service id list in ETCD
	val := &IdList{
		List: serviceIDList,
	}
	updateErr := s.Store.Update(ctx, path.Join(stackRootKey, stack.Id, servicesRootKey), val, 0)
	if updateErr != nil {
		createErr := s.Store.Create(ctx, path.Join(stackRootKey, stack.Id, servicesRootKey), val, nil, 0)
		if createErr != nil {
			s.rollbackStack(ctx, stack.Id, serviceIDList)
			return nil, createErr
		}
	}
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Running)); err != nil {
		return nil, err
	}
	reply := StackReply{
		StackId: stack.Id,
	}
	fmt.Printf("Stack is running %s\n", in.StackIdent)
	return &reply, nil
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
	if err := s.stopStackServices(ctx, stack.Id, false); err != nil {
		fmt.Printf("catch error during stop services: %v", err)
	}
	if err := s.removeStackNetworks(ctx, stack.Id, false); err != nil {
		fmt.Printf("catch error during remove networks: %v", err)
	}
	if err := stackStateMachine.TransitionTo(stack.Id, int32(StackState_Stopped)); err != nil {
		fmt.Printf("catch error during stack state transition: %v", err)
	}
	reply := StackReply{
		StackId: stack.Id,
	}
	empty := &IdList{
		List: []string{},
	}
	s.Store.Update(ctx, path.Join(stackRootKey, stack.Id, servicesRootKey), empty, 0)
	fmt.Printf("Stack stopped %s\n", in.StackIdent)
	return &reply, nil
}

func (s *Server) removeStackNetworks(ctx context.Context, ID string, force bool) error {
	networkList := &IdList{}
	err := s.Store.Get(ctx, path.Join(stackRootKey, ID, networksRootKey), networkList, true)
	if err != nil && !force {
		return err
	}
	return s.removeStackNetworksFromList(ctx, networkList.List)
}

func (s *Server) removeStackNetworksFromList(ctx context.Context, networkList []string) error {
	var removeErr error
	for _, key := range networkList {
		err := s.Docker.NetworkRemove(ctx, key)
		if err != nil {
			removeErr = err
		}
		fmt.Printf("Remove network %s\n", key)
	}
	if removeErr != nil {
		return removeErr
	}
	return nil
}

func (s *Server) stopStackServices(ctx context.Context, ID string, force bool) error {
	listKeys := &IdList{}
	err := s.Store.Get(ctx, path.Join(stackRootKey, ID, servicesRootKey), listKeys, true)
	if err != nil && !force {
		return err
	}
	server := service.Service{
		Docker: s.Docker,
	}
	var removeErr error
	for _, key := range listKeys.List {
		_, err := server.Remove(ctx, &service.RemoveRequest{
			Ident: key,
		})
		if err != nil {
			removeErr = err
		}
		s.Store.Delete(ctx, path.Join(servicesRootKey, key), false, nil)

	}
	if removeErr != nil {
		return removeErr
	}
	return nil
}

// Remove implements stack.ServerService Remove
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*StackReply, error) {
	request := &StackRequest{StackIdent: in.StackIdent}
	stack, errIdent := s.getStack(ctx, request)
	if errIdent != nil {
		return nil, errIdent
	}
	if !in.Force {
		if stopped, err := stackStateMachine.Is(stack.Id, int32(StackState_Stopped)); err != nil {
			return nil, err
		} else if !stopped {
			return nil, errors.New("The stack is not stopped")
		}
	} else {
		fmt.Printf("Removing services stack %s\n", in.StackIdent)
		s.stopStackServices(ctx, stack.Id, true)
		fmt.Printf("Removing networks stack %s\n", in.StackIdent)
		s.removeStackNetworks(ctx, stack.Id, true)
	}
	fmt.Printf("Removing stack %s\n", in.StackIdent)
	s.Store.Delete(ctx, path.Join(stackRootKey, stack.Id), true, nil)
	s.Store.Delete(ctx, path.Join(stackRootNameKey, stack.Name), true, nil)
	err := stackStateMachine.DeleteState(stack.Id)
	if err != nil {
		fmt.Printf("catching error: %v\n", err)
	}
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
	state, errGet := stackStateMachine.GetState(stack.Id)
	info.State = "nc"
	if errGet == nil {
		switch state {
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
	return &info
}

// newStackFromYaml create a new stack from yaml
func newStackFromYaml(ctx context.Context, config string) (stack *Stack, err error) {
	stack, err = ParseStackfile(ctx, config)
	if err != nil {
		return
	}

	// Create stack state
	if err = stackStateMachine.CreateState(stack.Id, int32(StackState_Stopped)); err != nil {
		return
	}

	return
}
