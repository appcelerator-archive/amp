package stack

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/pkg/docker/docker/stack"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"golang.org/x/net/context"
)

// Server is used to implement stack.StackServer
type Server struct {
	Accounts accounts.Interface
	Stacks   stacks.Interface
}

// Deploy implements stack.Server
func (s *Server) Deploy(ctx context.Context, in *DeployRequest) (*DeployReply, error) {
	log.Println("[stack] Deploy", in.String())

	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, w)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	fileName := fmt.Sprintf("/tmp/%d-%s.yml", time.Now().UnixNano(), in.Name)
	if err := ioutil.WriteFile(fileName, []byte(in.Compose), 0666); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}

	fullName, err := s.getStackInst(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	deployOpt := stack.NewDeployOptions(fullName, fileName, true)
	if err := stack.RunDeploy(dockerCli, deployOpt); err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "%v", err)
	}
	w.Close()
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &DeployReply{
		FullName: fullName,
		Answer:   string(outs),
	}
	return ans, nil
}

// verify if stack id alreaday exist, if yes, it's an update, if not create a new stack data entry
func (s *Server) getStackInst(ctx context.Context, name string) (string, error) {
	if stackInst, err := s.Stacks.GetStack(ctx, name); err == nil && stackInst != nil {
		return fmt.Sprintf("%s-%s", stackInst.Name, stackInst.Id), nil
	}
	ids := strings.Split(name, "-")
	id := ids[len(ids)-1]
	if stackInst, err := s.Stacks.GetStack(ctx, id); err == nil && stackInst != nil {
		return name, nil
	}
	stackInst, err := s.Stacks.CreateStack(ctx, name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", stackInst.Name, stackInst.Id), nil
}

// List implements stack.Server
func (s *Server) List(ctx context.Context, in *ListRequest) (*ListReply, error) {
	log.Println("[stack] List", in.String())

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
	cols := strings.Split(outs, "\n")
	ans := &ListReply{
		List: []*StackReply{},
	}
	for _, col := range cols[1:] {
		stack := s.getOneStackListLine(ctx, col)
		if stack.Id != "" {
			ans.List = append(ans.List, stack)
		}
	}
	return ans, nil
}

func (s *Server) getOneStackListLine(ctx context.Context, line string) *StackReply {
	cols := strings.Split(line, " ")
	name := cols[0]
	id := ""
	ll := strings.LastIndex(cols[0], "-")
	if ll >= 0 {
		id = name[ll+1:]
		name = name[0:ll]
	}
	ret := &StackReply{
		Id:   strings.Trim(id, " "),
		Name: strings.Trim(name, " "),
	}
	if stackInst, err := s.Stacks.GetStack(ctx, id); err == nil && stackInst != nil {
		ret.Owner = stackInst.Owner.Name
	}
	for _, col := range cols[1:] {
		if col != "" {
			ret.Service = col
			return ret
		}
	}
	return ret
}

// Remove implements stack.Server
func (s *Server) Remove(ctx context.Context, in *RemoveRequest) (*RemoveReply, error) {
	log.Println("[stack] Remove", in.String())

	r, w, _ := os.Pipe()
	dockerCli := command.NewDockerCli(os.Stdin, w, w)
	opts := cliflags.NewClientOptions()
	if err := dockerCli.Initialize(opts); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("error in cli initialize: %v", err))
	}
	name := in.Id
	stackInst, err := s.Stacks.GetStack(ctx, in.Id)
	if err == nil && stackInst != nil {
		name = fmt.Sprintf("%s-%s", stackInst.Name, stackInst.Id)
	} else {
		return nil, fmt.Errorf("Stack %s is not an amp stack", in.Id)
	}
	if !s.Accounts.IsAuthorized(ctx, stackInst.Owner, accounts.DeleteAction, accounts.StackRN, stackInst.Id) {
		return nil, grpc.Errorf(codes.PermissionDenied, "user not authorized")
	}
	rmOpt := stack.NewRemoveOptions([]string{name})
	if err := stack.RunRemove(dockerCli, rmOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	lid := strings.Split(name, "-")
	if len(lid) >= 2 {
		id := lid[len(lid)-1]
		if err := s.Stacks.DeleteStack(ctx, id); err != nil {
			return nil, grpc.Errorf(codes.Internal, "%v", err)
		}
	}
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &RemoveReply{
		Answer: string(outs),
	}
	log.Printf("Stack %s removed", in.Id)
	return ans, nil
}
