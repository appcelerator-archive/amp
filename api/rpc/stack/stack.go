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

	"github.com/appcelerator/amp/data/stacks"
	"github.com/appcelerator/amp/pkg/docker/docker/stack"
	"github.com/docker/docker/cli/command"
	cliflags "github.com/docker/docker/cli/flags"
	"golang.org/x/net/context"
)

// Server is used to implement stack.StackServer
type Server struct {
	Stacks stacks.Interface
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
	stackInst, err := s.Stacks.CreateStack(ctx, in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	fullName := fmt.Sprintf("%s-%s", stackInst.Name, stackInst.Id)
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
	ans := &ListReply{
		Answer: string(outs),
	}
	return ans, nil
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
	rmOpt := stack.NewRemoveOptions([]string{in.Id})
	if err := stack.RunRemove(dockerCli, rmOpt); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	w.Close()
	lid := strings.Split(in.Id, "-")
	if len(lid) < 2 {
		return nil, grpc.Errorf(codes.Internal, "%v", fmt.Errorf("Bad stack name format"))
	}
	id := lid[len(lid)-1]
	if err := s.Stacks.DeleteStack(ctx, id); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	out, _ := ioutil.ReadAll(r)
	outs := strings.Replace(string(out), "docker", "amp", -1)
	ans := &RemoveReply{
		Answer: string(outs),
	}
	log.Printf("Stack %s removed", in.Id)
	return ans, nil
}
