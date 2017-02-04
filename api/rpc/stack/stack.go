package stack

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/appcelerator/amp/api/rpc/stack/docker/stack"
	"github.com/appcelerator/amp/data/storage"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	//Name of the amp core stack, the minimal needed set of services
	ampCoreStackName = "ampcore"
	//Name of the label to create mapping between a free label (part of host) and a service port
	haproxyMappingLabelName = "io.amp.mapping"
	//Name of the root key in etcd under which stack information are stored
	stackRootKey = "stacks"
	//Path where the stack infrastructure compose files are stored in the amplifier container
	stackFilePath = "/var/lib/amp"
	//StackFileVarName Name of the file used to store the image names and the images version used in infratructure compose ymy files
	StackFileVarName = "amp.var"
)

// InfraStackList list of all the infrastructure amp stacks
var InfraStackList = []string{ampCoreStackName, "ampfunction", "amplog", "ampmonitor", "ampadmin"}

// InfraShortStackList list of the infrastructure amp stacks wich are started using 'all' keyword
var InfraShortStackList = []string{"ampfunction", "amplog", "ampmonitor", "ampadmin"}

// Server is used to implement stack.StackService
type Server struct {
	Store      storage.Interface
	Docker     *client.Client
	serviceMap map[string]*ampService
}

// NewServer instantiates a server
func NewServer(store storage.Interface, docker *client.Client) *Server {
	return &Server{
		Store:      store,
		Docker:     docker,
		serviceMap: make(map[string]*ampService),
	}
}

// Deploy a stack infrastructure or user
func (s *Server) Deploy(ctx context.Context, in *StackDeployRequest) (*StackReply, error) {
	if err := parseStack(in.Stack); err != nil {
		return nil, err
	}
	fileName := fmt.Sprintf("/tmp/%d-%s.yml", time.Now().UnixNano, in.Stack.Name)
	if err := ioutil.WriteFile(fileName, []byte(in.Stack.FileData), 0777); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, w)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	deployOpt := stack.NewDeployOptions(in.Stack.Name, fileName, in.RegistryAuth)
	if err := stack.RunDeploy(dockerCli, deployOpt); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &StackReply{
		Answer: string(outs),
	}
	if err := s.storeHAProxyMappings(ctx, in.Stack); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return ans, nil
}

//store in etcd all haproxy mapping of a stack
func (s *Server) storeHAProxyMappings(ctx context.Context, stack *Stack) error {
	if stack.Name == ampCoreStackName {
		return nil
	}
	for _, serv := range stack.Services {
		labels := serv.Labels
		for name, value := range labels {
			if name == haproxyMappingLabelName {
				if _, err := EvalMappingString(value); err != nil {
					return fmt.Errorf("%v\nstack %s is started but this haproxy mapping is not done", err, stack.Name)
				}
				if err := s.storeHAProxyMapping(ctx, stack.Name, serv.Name, value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//store in etcd one haproxy mapping
func (s *Server) storeHAProxyMapping(ctx context.Context, stackName string, serviceName string, mapping string) error {
	mappingStore := &StackMapping{
		StackName:   stackName,
		ServiceName: serviceName,
		Mapping:     mapping,
	}
	mapPath := path.Join(stackRootKey, stackName, serviceName)
	if erru := s.Store.Update(ctx, mapPath, mappingStore, 0); erru != nil {
		if err := s.Store.Create(ctx, mapPath, mappingStore, nil, 0); err != nil {
			return err
		}
	}
	log.Printf("Added HAProxy mappings stack=%s service=%s mapping=%s\n", stackName, serviceName, mapping)
	return nil
}

// UpdateFile Update one yml file into amplifier: copy it to .sav before remplacing it. TODO: make it working with several amplifier instances
func (s *Server) UpdateFile(ctx context.Context, in *StackUpdateRequest) (*StackReply, error) {
	found := false
	if in.FileName == "variable" {
		in.FileName = StackFileVarName
		found = true
	} else {
		for _, name := range InfraStackList {
			if name == in.FileName {
				found = true
			}
		}
	}
	if !found {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", fmt.Errorf("%s is not an infrastructure name nor 'variable' for the variables file", in.FileName))
	}
	targetName := path.Join(stackFilePath, in.FileName)
	targetNameSaved := path.Join(stackFilePath, in.FileName+".sav")
	os.Remove(targetNameSaved)
	os.Rename(targetName, targetNameSaved)
	err := ioutil.WriteFile(targetName, []byte(in.FileData), 0666)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	log.Printf("File %s updated\n", in.FileName)
	ans := &StackReply{Answer: "done"}
	return ans, nil
}

// RestoreFile Restore one yml file into amplifier. if .sav exist, copy it in place of the regular one. TODO: make it working with several amplifier instances
func (s *Server) RestoreFile(ctx context.Context, in *StackRestoreRequest) (*StackReply, error) {
	found := false
	if in.FileName == "variable" {
		in.FileName = StackFileVarName
		found = true
	} else {
		for _, name := range InfraStackList {
			if name == in.FileName {
				found = true
			}
		}
	}
	if !found {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", fmt.Errorf("%s is not an infrastructure name nor 'variable' for the variables file", in.FileName))
	}
	targetName := path.Join(stackFilePath, in.FileName)
	targetNameSaved := path.Join(stackFilePath, in.FileName+".sav")

	if _, err := os.Stat(targetNameSaved); err != nil {
		return nil, grpc.Errorf(codes.NotFound, "%v", fmt.Errorf("There is no previous file saved for this stack: %s", in.FileName))
	}
	os.Remove(targetName)
	os.Rename(targetNameSaved, targetName)
	log.Printf("File %s restored\n", in.FileName)
	ans := &StackReply{Answer: "done"}
	return ans, nil
}

//ListStacks Return the list of all stacks up including infrastructure ones
func (s *Server) ListStacks(ctx context.Context, in *StackRequest) (*StackReply, error) {
	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, os.Stderr)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	listOpt := stack.NewListOptions()
	if err := stack.RunList(dockerCli, listOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &StackReply{
		Answer: string(outs),
	}
	return ans, nil
}

//ListTasks List the tasks of a stack, infrastrucuture or user
func (s *Server) ListTasks(ctx context.Context, in *StackPsRequest) (*StackReply, error) {
	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, os.Stderr)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	psOpt := stack.NewPsOptions(in.Name, in.NoTrunc, in.NoResolve, in.Filter)
	if err := stack.RunPS(dockerCli, psOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &StackReply{
		Answer: string(outs),
	}
	return ans, nil
}

//Remove a stack, infrastrucuture or user
func (s *Server) Remove(ctx context.Context, in *StackRequest) (*StackReply, error) {
	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, w)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	rmOpt := stack.NewRemoveOptions(in.Name)
	if err := stack.RunRemove(dockerCli, rmOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &StackReply{
		Answer: string(outs),
	}
	s.removeHAProxyMappings(ctx, in.Name)
	return ans, nil
}

// Remove all haproxy mapping of a stack
func (s *Server) removeHAProxyMappings(ctx context.Context, stackName string) {
	if stackName == ampCoreStackName {
		return
	}
	s.Store.Delete(ctx, path.Join(stackRootKey, stackName), true, nil)
}

// Services list the services of a stack, infrastrucuture or user
func (s *Server) Services(ctx context.Context, in *StackServicesRequest) (*StackReply, error) {
	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, os.Stderr)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, fmt.Errorf("rror in cli initialize: %v", err)
	}
	servicesOpt := stack.NewServicesOptions(in.Name, in.Quiet, in.Filter)
	if err := stack.RunServices(dockerCli, servicesOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &StackReply{
		Answer: string(outs),
	}
	return ans, nil
}

// GetStackStatus list the status of all stacks up
func (s *Server) GetStackStatus(ctx context.Context, in *StackRequest) (*StackReply, error) {
	stack := &Stack{Name: in.Name}
	if err := parseStack(stack); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	ready := 0
	starting := 0
	failing := 0
	for _, serv := range stack.Services {
		id, exist := s.DoesServiceExist(ctx, fmt.Sprintf("%s_%s", in.Name, serv.Name))
		if exist {
			sReady, sStarting, sFailing := s.getServiceStatus(ctx, id)
			if sReady > 0 {
				ready++
			} else if sFailing > 0 {
				failing++
			} else if sStarting > 0 {
				starting++
			}
		}
	}
	status := s.getServiceStatusString(ready, starting, failing)
	ans := &StackReply{
		Answer: status,
	}
	return ans, nil
}

// Monitor return a list of all services of all stacks up
func (s *Server) Monitor(ctx context.Context, in *StackRequest) (*StackMonitorReply, error) {
	ret, err := s.GetMonitorLines(ctx)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &StackMonitorReply{Lines: ret}, nil
}

//GetImages return the list of all needed infrastructure images
func (s *Server) GetImages(ctx context.Context, in *StackRequest) (*StackImagesReply, error) {
	local := false
	if in.Name == "local" {
		local = true
	}
	list, err := s.getImages(local)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &StackImagesReply{Images: list}, nil
}

// PullImage pull one image
func (s *Server) PullImage(ctx context.Context, in *StackRequest) (*StackReply, error) {
	if err := s.pullImage(ctx, in.Name); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	return &StackReply{}, nil
}
