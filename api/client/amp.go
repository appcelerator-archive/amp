package client

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/api/rpc/stat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strconv"
)

const (
	defaultPort   = ":50101"
	serverAddress = "localhost" + defaultPort
)

// Configuration is for all configurable client settings
type Configuration struct {
	Verbose bool
	Github  string
	Target  string
	Images  []string
}

// AMP holds the state for the current envirionment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Configuration *Configuration
	Conn          *grpc.ClientConn
}

func (a *AMP) connect() {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	a.Conn = conn
}

func (a *AMP) disconnect() {
	err := a.Conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (a *AMP) getAuthorizedContext() (ctx context.Context, err error) {
	if a.Configuration.Github == "" {
		return nil, fmt.Errorf("Requires login")
	}
	md := metadata.Pairs("sessionkey", a.Configuration.Github)
	ctx = metadata.NewContext(context.Background(), md)
	return
}

func (a *AMP) verbose() bool {
	return a.Configuration.Verbose
}

// NewAMP creates a new AMP instance
func NewAMP(c *Configuration) *AMP {
	return &AMP{Configuration: c}
}

// Create a new swarm
func (a *AMP) Create() {
	if a.verbose() {
		fmt.Println("Create")
	}
}

// Start the swarm
func (a *AMP) Start() {
	if a.verbose() {
		fmt.Println("Start")
	}
}

// Update the swarm
func (a *AMP) Update() {
	if a.verbose() {
		fmt.Println("Update")
	}
}

// Stop the swarm
func (a *AMP) Stop() {
	if a.verbose() {
		fmt.Println("Stop")
	}
}

// Status returns the current status
func (a *AMP) Status() {
	if a.verbose() {
		fmt.Println("Status")
	}
}

// Logs fetches the logs
func (a *AMP) Logs(cmd *cobra.Command) error {
	ctx, err := a.getAuthorizedContext()
	if err != nil {
		return err
	}
	if a.verbose() {
		fmt.Println("Logs")
		fmt.Printf("service_id: %v\n", cmd.Flag("service_id").Value)
		fmt.Printf("service_name: %v\n", cmd.Flag("service_name").Value)
		fmt.Printf("message: %v\n", cmd.Flag("message").Value)
		fmt.Printf("container_id: %v\n", cmd.Flag("container_id").Value)
		fmt.Printf("node_id: %v\n", cmd.Flag("node_id").Value)
		fmt.Printf("from: %v\n", cmd.Flag("from").Value)
		fmt.Printf("n: %v\n", cmd.Flag("n").Value)
		fmt.Printf("short: %v\n", cmd.Flag("short").Value)
	}
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := logs.GetRequest{}
	request.ServiceId = cmd.Flag("service_id").Value.String()
	request.ServiceName = cmd.Flag("service_name").Value.String()
	request.Message = cmd.Flag("message").Value.String()
	request.ContainerId = cmd.Flag("container_id").Value.String()
	request.NodeId = cmd.Flag("node_id").Value.String()
	if request.From, err = strconv.ParseInt(cmd.Flag("from").Value.String(), 10, 64); err != nil {
		log.Panicf("Unable to convert from parameter: %v\n", cmd.Flag("from").Value.String())
	}
	if request.Size, err = strconv.ParseInt(cmd.Flag("n").Value.String(), 10, 64); err != nil {
		log.Panicf("Unable to convert n parameter: %v\n", cmd.Flag("n").Value.String())
	}

	c := logs.NewLogsClient(conn)
	r, err := c.Get(ctx, &request)
	if err != nil {
		return err
	}
	for _, entry := range r.Entries {
		var short bool
		if short, err = strconv.ParseBool(cmd.Flag("short").Value.String()); err != nil {
			log.Panicf("Unable to convert short parameter: %v\n", cmd.Flag("short").Value.String())
		}
		if short {
			fmt.Printf("%s\n", entry.Message)
		} else {
			fmt.Printf("%+v\n", entry)
		}
	}
	conn.Close()
	return nil
}


//CPU Display CPU Stats
func (a *AMP) CPU(cmd *cobra.Command) error {
	ctx, err := a.getAuthorizedContext()
	if err != nil {
		return err
	}
	if a.verbose() {
		fmt.Println("Cpu")
		fmt.Printf("Ressource: %v\n", cmd.Flag("ressourceName").Value)
	}
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return err
	}

	request := stat.CPURequest{}
	request.RessourceName = cmd.Flag("RessourceName").Value.String()

	config := stat.Config{
		Connstr: "http://influxdb:8086",
		Dbname:  "telegraf",
		U:       "",
		P:       "",
	}
	c := logs.NewsStatClient(conn)
	r, err := c.CPUQuery(ctx, &request)
	if err != nil {
		return err
	}
	for _, entry := range r.Entries {
		fmt.Printf("%s %s\t\t%d", entry.ID, entry.Name, (entry.UsageUser*100)/entry.UsageTotal)
		//TODO format and
		//add (entry.UsageKernel*100)/entry.UsageTotal
		//add (entry.UsageSystem*100)/entry.UsageTotal

	}
	conn.Close()
	return nil
}
