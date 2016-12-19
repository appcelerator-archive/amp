package agentcore

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	//"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClusterAgent struct
type ClusterAgent struct {
	id           string
	realID       string
	nodeID       string
	hostname     string
	address      string
	dockerClient *client.Client
	client       servergrpc.ClusterServerServiceClient
	conn         *grpc.ClientConn
	token        string
	ctx          context.Context
	healthy      bool
}

//Init Connect to docker engine, get initial containers list and start the agent
func (g *ClusterAgent) Init(version string, build string) error {
	g.setToken()
	g.ctx = context.Background()
	g.trapSignal()
	conf.init(version, build)
	g.id = os.Getenv("HOSTNAME")

	// Connection to Docker
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(conf.dockerEngine, "v1.24", nil, defaultHeaders)
	if err != nil {
		return err
	}
	g.dockerClient = cli
	fmt.Println("Connected to Docker-engine")
	g.inspectContainer()
	//start GRPC
	g.startGRPCServer()
	// Connection to server
	if err := g.connectServer(); err != nil {
		return err
	}
	fmt.Println("Connected to adm-server")
	g.startHeartBeat()
	return nil
}

func (g *ClusterAgent) connectServer() error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", conf.serverAddr, conf.serverPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*20))
	if err != nil {
		return err
	}
	g.conn = conn
	g.client = servergrpc.NewClusterServerServiceClient(conn)
	logf.info("Connected to server\n")
	ctx := metadata.NewContext(context.Background(), metadata.Pairs("token", g.token))
	_, errReg := g.client.RegisterAgent(ctx, &servergrpc.RegisterRequest{
		Id:       g.id,
		NodeId:   g.nodeID,
		Hostname: g.hostname,
		Address:  g.address,
	})
	if errReg != nil {
		return fmt.Errorf("Register on server error: %v", errReg)
	}
	logf.info("Agent registered\n")
	g.healthy = true
	return nil
}

func (g *ClusterAgent) startGRPCServer() {
	serv := grpc.NewServer()
	agentgrpc.RegisterClusterAgentServiceServer(serv, g)
	logf.info("Starting GRPC server\n")
	lis, err := net.Listen("tcp", ":"+conf.grpcPort)
	if err != nil {
		logf.error("adm-agent is unable to listen on: %s\n%v", ":"+conf.grpcPort, err)
		return
	}
	go func() {
		logf.info("adm-agent is listening on port %s\n", conf.grpcPort)
		if err := serv.Serve(lis); err != nil {
			logf.error("Problem in adm-agent: %s\n", err)
			return
		}
	}()
	time.Sleep(2 * time.Second)
}

func (g *ClusterAgent) inspectContainer() {
	inspect, err := g.dockerClient.ContainerInspect(g.ctx, g.id)
	if err != nil {
		logf.error("Error inspecting container: %v\n", err)
		return
	}
	g.realID = inspect.ID
	g.nodeID = inspect.Config.Labels["com.docker.swarm.node.id"]
	g.address = fmt.Sprintf("%s:%s", inspect.NetworkSettings.Networks["amp-infra"].IPAddress, conf.grpcPort)
	logf.debug("Agent: id=%s, address=%s nodeId=%s\n", g.id, g.address, g.nodeID)
}

func (g *ClusterAgent) setToken() {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 64)
	rand.Read(b)
	g.token = fmt.Sprintf("%x", b)
}

func (g *ClusterAgent) startHeartBeat() {
	for {
		if _, err := g.client.AgentHealth(g.ctx, &servergrpc.AgentHealthRequest{Id: g.nodeID}); err != nil {
			logf.error("Server is not reachable anymore. Exit\n")
			os.Exit(1)
		}
		time.Sleep(10 * time.Second)
	}
}

// Launch a routine to catch SIGTERM Signal
func (g *ClusterAgent) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("adm-agent received SIGTERM signal")
		g.conn.Close()
		os.Exit(1)
	}()
}
