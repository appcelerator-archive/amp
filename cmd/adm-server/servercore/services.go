package servercore

import (
	"fmt"
	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

// GetClientStream connect a bidirectionnal-stream
func (s *ClusterServer) GetClientStream(stream servergrpc.ClusterServerService_GetClientStreamServer) error {
	logf.debug("Receive client connection\n")
	clientID := fmt.Sprintf("Client-%d", len(s.clientMap)+1)
	s.clientMap[clientID] = &AmpClient{
		stream: stream,
	}
	err := stream.Send(&servergrpc.ClientMes{
		Function: "ClientAck",
		ClientId: clientID,
	})
	if err != nil {
		logf.error("Error sending back ack to client %v\n", err)
		return grpc.Errorf(codes.Internal, "%v\n", err)
	}
	logf.debug("Client id=%s\n", clientID)
	for {
		mes, err := stream.Recv()
		if err != nil {
			delete(s.clientMap, clientID)
			return grpc.Errorf(codes.Internal, "Stream Server-client ended: %v\n", err)
		}
		logf.debug("received client message: %v\n", mes)
	}
}

// RegisterAgent register an agent
func (s *ClusterServer) RegisterAgent(ctx context.Context, req *servergrpc.RegisterRequest) (*servergrpc.ServerRet, error) {
	agent := &Agent{
		agentID: req.Id,
		nodeID:  req.NodeId,
		address: req.Address,
	}
	logf.debug("Received agent register request: %v\n", req)
	md, exist := metadata.FromContext(ctx)
	if !exist {
		logf.error("RegisterAgent %s no token found", req.Id)
		return nil, grpc.Errorf(codes.NotFound, "Token not found")
	}
	agent.token = md["token"][0]
	if err := s.connectBackAgent(agent); err != nil {
		logf.error("RegisterAgent %s connect back error: %v", req.Id, err)
		return nil, grpc.Errorf(codes.Internal, "%v\n", err)
	}
	s.updateAgentInfo(ctx, agent)
	s.agentMap[req.NodeId] = agent
	return &servergrpc.ServerRet{AgentId: agent.nodeName, Error: ""}, nil
}

// AgentHealth verifie if the server i is still available
func (s *ClusterServer) AgentHealth(ctx context.Context, req *servergrpc.AgentHealthRequest) (*servergrpc.ServerRet, error) {
	agent, ok := s.agentMap[req.Id]
	if !ok {
		logf.warn("Received heartbeat from a not registered agent: %s\n", req.Id)
		return nil, grpc.Errorf(codes.FailedPrecondition, "Agent not registered")
	}
	agent.lastBeat = time.Now()
	return &servergrpc.ServerRet{}, nil
}

func (s *ClusterServer) updateAgentInfo(ctx context.Context, agent *Agent) error {
	node, _, err := s.dockerClient.NodeInspectWithRaw(ctx, agent.nodeID)
	if err != nil {
		return grpc.Errorf(codes.Internal, "%v", err)
	}
	agent.nodeID = node.ID
	agent.role = string(node.Spec.Role)
	agent.availability = string(node.Spec.Availability)
	agent.hostname = node.Description.Hostname
	agent.hostArchitecture = node.Description.Platform.Architecture
	agent.hostOs = node.Description.Platform.OS
	agent.dockerVersion = node.Description.Engine.EngineVersion
	agent.status = string(node.Status.State)
	if agent.role == "manager" {
		agent.leader = node.ManagerStatus.Leader
	}
	agent.nodeName = fmt.Sprintf("%s (%s)", agent.hostname, agent.nodeID)
	return nil
}

// GetNodesInfo get node info
func (s *ClusterServer) GetNodesInfo(ctx context.Context, req *servergrpc.GetNodesInfoRequest) (*servergrpc.NodesInfo, error) {
	ret := &servergrpc.NodesInfo{}
	if req.Node != "" {
		agent, ok := s.agentMap[req.Node]
		if !ok {
			return nil, grpc.Errorf(codes.NotFound, "Node %s doesn't exist or is not registered", agent.nodeName)
		}
		inf := s.getNodeInfo(ctx, req, agent)
		ret.Nodes = append(ret.Nodes, inf)
	} else {
		for _, agent := range s.agentMap {
			inf := s.getNodeInfo(ctx, req, agent)
			ret.Nodes = append(ret.Nodes, inf)

		}
	}
	return ret, nil
}

func (s *ClusterServer) getNodeInfo(ctx context.Context, req *servergrpc.GetNodesInfoRequest, agent *Agent) *servergrpc.NodeInfo {
	areq := &agentgrpc.GetNodeInfoRequest{}
	inf := &servergrpc.NodeInfo{}
	inf.Address = agent.address
	ad := strings.Split(agent.address, ":")
	if len(ad) > 1 {
		inf.Address = ad[0]
	}
	if err := s.updateAgentInfo(ctx, agent); err != nil {
		inf.Error = err.Error()
	} else {
		inf.NodeName = agent.nodeName
		inf.Id = agent.nodeID
		inf.AgentId = agent.agentID
		inf.Role = agent.role
		inf.Availability = agent.availability
		inf.Hostname = agent.hostname
		inf.HostArchitecture = agent.hostArchitecture
		inf.HostOs = agent.hostOs
		inf.DockerVersion = agent.dockerVersion
		inf.Status = agent.status
	}
	if info, err := agent.client.GetNodeInfo(ctx, areq); err != nil {
		if inf.Error == "" {
			inf.Error = err.Error()
		} else {
			inf.Error = fmt.Sprintf("%s\n%s", inf.Error, err.Error())
		}
	} else {
		inf.Cpu = info.Cpu
		inf.Memory = info.Memory
		inf.NbContainers = info.NbContainers
		inf.NbContainersRunning = info.NbContainersRunning
		inf.NbContainersPaused = info.NbContainersPaused
		inf.NbContainersStopped = info.NbContainersStopped
		inf.Images = info.Images
	}
	return inf
}

// PurgeNodes purge nodes
func (s *ClusterServer) PurgeNodes(ctx context.Context, req *servergrpc.PurgeNodesRequest) (*servergrpc.PurgeNodesAnswers, error) {
	sret := &servergrpc.PurgeNodesAnswers{}
	if req.Node != "" {
		agent, ok := s.agentMap[req.Node]
		if !ok {
			return nil, grpc.Errorf(codes.NotFound, "Node %s doesn't exist or is not registered", agent.nodeName)
		}
		aret := s.purgeNode(ctx, req, agent)
		sret.Agents = append(sret.Agents, aret)
	} else {
		for _, agent := range s.agentMap {
			aret := s.purgeNode(ctx, req, agent)
			sret.Agents = append(sret.Agents, aret)
		}
	}
	return sret, nil
}

func (s *ClusterServer) purgeNode(ctx context.Context, req *servergrpc.PurgeNodesRequest, agent *Agent) *servergrpc.PurgeNodeAnswer {
	ret, err := agent.client.PurgeNode(ctx, &agentgrpc.PurgeNodeRequest{
		Node:      req.Node,
		Container: req.Container,
		Volume:    req.Volume,
		Image:     req.Image,
	})
	aret := &servergrpc.PurgeNodeAnswer{
		AgentId: agent.nodeName,
		Error:   "ok",
	}
	if err != nil {
		aret.Error = err.Error()
	} else {
		aret.NbContainers = ret.NbContainers
		aret.NbVolumes = ret.NbVolumes
		aret.NbImages = ret.NbImages
	}
	return aret
}

// AmpMonitor amp monitor
func (s *ClusterServer) AmpMonitor(tx context.Context, req *servergrpc.AmpRequest) (*servergrpc.AmpMonitorAnswers, error) {
	manager := &AMPInfraManager{}
	if err := manager.Init(s, req.ClientId, ""); err != nil {
		return nil, err
	}
	lineList, err := manager.Monitor(getAMPInfrastructureStack(manager))
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "Monitor error: %v", err)
	}
	return &servergrpc.AmpMonitorAnswers{Outputs: *lineList}, nil
}

// AmpPull amp pull
func (s *ClusterServer) AmpPull(tx context.Context, req *servergrpc.AmpRequest) (*servergrpc.AmpRet, error) {
	manager := &AMPInfraManager{}
	if err := manager.Init(s, req.ClientId, "Pulling AMP images"); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	manager.Verbose = req.Verbose
	manager.Silence = req.Silence
	manager.Force = req.Force
	manager.Local = req.Local
	manager.pullCount = 0
	max := len(s.agentMap)
	if req.Node != "" {
		max = 1
		agent, ok := s.agentMap[req.Node]
		if !ok {
			return nil, fmt.Errorf("Node %s doesn't exist or is not registered", agent.nodeName)
		}
		s.pullOnAgent(manager, agent)
	} else {
		for _, agent := range s.agentMap {
			s.pullOnAgent(manager, agent)
		}
	}
	t0 := time.Now()
	for manager.pullCount < max {
		time.Sleep(1 * time.Second)
		if time.Now().Sub(t0).Seconds() > 120*5 {
			return nil, fmt.Errorf("Pull timeout")
		}
	}
	return &servergrpc.AmpRet{}, nil
}

func (s *ClusterServer) pullOnAgent(manager *AMPInfraManager, agent *Agent) {
	go func() {
		ok, ko := manager.Pull(agent, getAMPInfrastructureStack(manager))
		if ok == 0 && ko == 0 {
			manager.printf(colWarn, "Agent %s: No image have been pulled at all", agent.nodeName)
		} else if ok == 0 && ko > 0 {
			manager.printf(colError, "Agent %s: Pull error: %d, no image have been pulled", agent.nodeName, ko)
		} else if ok > 0 && ko == 0 {
			manager.printf(colSuccess, "Agent %s: Images pulled: %d", agent.nodeName, ok)
		} else {
			manager.printf(colSuccess, "Agent %s: Images pulled: %d, Images pull error: %d", agent.nodeName, ok, ko)
		}
		manager.pullCount++
	}()
}

// GetAmpStatus get amp status
func (s *ClusterServer) GetAmpStatus(tx context.Context, req *servergrpc.AmpRequest) (*servergrpc.AmpStatusAnswer, error) {
	manager := &AMPInfraManager{}
	manager.Init(s, req.ClientId, "")
	manager.Verbose = req.Verbose
	manager.Silence = req.Silence
	manager.Force = req.Force
	manager.ComputeStatus(getAMPInfrastructureStack(manager))
	manager.Pregularf("status: %s\n", manager.Status)
	return &servergrpc.AmpStatusAnswer{Status: manager.Status}, nil
}

// AmpStart amp start
func (s *ClusterServer) AmpStart(tx context.Context, req *servergrpc.AmpRequest) (*servergrpc.AmpRet, error) {
	manager := &AMPInfraManager{}
	if err := manager.Init(s, req.ClientId, "Starting AMP platform"); err != nil {
		return nil, grpc.Errorf(codes.Internal, "%v", err)
	}
	manager.Verbose = req.Verbose
	manager.Silence = req.Silence
	manager.Force = req.Force
	manager.Local = req.Local
	stack := getAMPInfrastructureStack(manager)
	manager.ComputeStatus(stack)
	if manager.Status == "running" {
		if !manager.Force {
			manager.printf(colWarn, "AMP platform already started (-f to force a re-start)")
			return &servergrpc.AmpRet{}, nil
		}
		if err := manager.Stop(stack); err != nil {
			logf.error("Error stopting amp at (1): %v\n", err)
		}
	}
	if err := manager.systemPrerequisites(); err != nil {
		logf.printf("Prerequisite error: %v\n", err)
		return nil, grpc.Errorf(codes.Internal, "Prerequisite error: %v\n", err)
	}
	if err := manager.Start(stack); err != nil {
		logf.error("Error starting amp: %v\n", err)
		if err := manager.Stop(stack); err != nil {
			logf.error("Error stopting amp at (2): %v\n", err)
		}
		return nil, grpc.Errorf(codes.Internal, "Start error: %v\n", err)
	}
	return &servergrpc.AmpRet{}, nil
}

// AmpStop amp stop
func (s *ClusterServer) AmpStop(tx context.Context, req *servergrpc.AmpRequest) (*servergrpc.AmpRet, error) {
	manager := &AMPInfraManager{}
	if err := manager.Init(s, req.ClientId, "Stopping AMP platform"); err != nil {
		return nil, err
	}
	manager.Verbose = req.Verbose
	manager.Silence = req.Silence
	manager.Force = req.Force
	manager.Local = req.Local
	stack := getAMPInfrastructureStack(manager)
	manager.ComputeStatus(stack)
	if manager.Status == "stopped" {
		manager.printf(colWarn, "AMP platform already stopped")
		return &servergrpc.AmpRet{}, nil
	}
	if err := manager.Stop(stack); err != nil {
		return nil, grpc.Errorf(codes.Internal, "Stop error: %v\n", err)
	}
	return &servergrpc.AmpRet{}, nil
}
