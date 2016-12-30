package servercore

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"github.com/docker/docker/api/types"
	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

//ClusterServer struct
type ClusterServer struct {
	dockerClient *client.Client
	agentMap     map[string]*Agent
	clientMap    map[string]*AmpClient
}

//Agent struct
type Agent struct {
	agentID          string
	nodeName         string
	nodeID           string
	address          string
	token            string
	role             string
	availability     string
	hostname         string
	hostArchitecture string
	hostOs           string
	dockerVersion    string
	status           string
	leader           bool
	lastBeat         time.Time
	client           agentgrpc.ClusterAgentServiceClient
	conn             *grpc.ClientConn
}

//AmpClient struct
type AmpClient struct {
	stream servergrpc.ClusterServerService_GetClientStreamServer
}

//Init Connect to docker engine, get initial containers list and start the agent
func (s *ClusterServer) Init(version string, build string) error {
	s.agentMap = make(map[string]*Agent)
	s.clientMap = make(map[string]*AmpClient)
	s.trapSignal()
	conf.init(version, build)

	// Connection to Docker
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(conf.dockerEngine, "v1.24", nil, defaultHeaders)
	if err != nil {
		return err
	}
	s.dockerClient = cli
	logf.info("Connected to Docker-engine\n")

	// Start server
	s.startGRPCServer()
	logf.info("GRPC server started\n")
	for {
		time.Sleep(1 * time.Second)
	}
}

func (s *ClusterServer) startGRPCServer() {
	serv := grpc.NewServer()
	servergrpc.RegisterClusterServerServiceServer(serv, s)
	go func() {
		logf.info("Starting GRPC server\n")
		lis, err := net.Listen("tcp", ":"+conf.grpcPort)
		if err != nil {
			logf.error("adm-server is unable to listen on: %s\n%v", ":"+conf.grpcPort, err)
			return
		}
		logf.info("adm-server is listening on port %s\n", conf.grpcPort)
		if err := serv.Serve(lis); err != nil {
			logf.error("Problem in adm-server: %s\n", err)
		}
	}()
}

func (s *ClusterServer) connectBackAgent(agent *Agent) error {
	logf.debug("ConnectBackAgent: %s\n", agent.address)
	conn, err := grpc.Dial(agent.address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*20))
	if err != nil {
		return err
	}
	agent.conn = conn
	agent.client = agentgrpc.NewClusterAgentServiceClient(conn)
	return nil
}

// Launch a routine to catch SIGTERM Signal
func (s *ClusterServer) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("adm-server received SIGTERM signal")
		os.Exit(1)
	}()
}
