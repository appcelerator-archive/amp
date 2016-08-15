package cli

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

const (
	defaultPort   = ":50101"
	serverAddress = "localhost" + defaultPort
)

// AMP holds the state for the current envirionment
type AMP struct {
	// Config contains all the configuration settings that were loaded
	Config *Config
}

func (a *AMP) verbose() bool {
	return a.Config.Verbose
}

// NewAMP creates a new AMP instance
func NewAMP(c *Config) *AMP {
	return &AMP{Config: c}
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
	r, err := c.Get(context.Background(), &request)
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
