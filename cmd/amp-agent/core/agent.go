package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"github.com/appcelerator/amp/pkg/config"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	defaultTimeOut    = 30 * time.Second
	natsClientID      = "amp-agent"
	containersDateDir = "/containers"
)

//Agent data
type Agent struct {
	dockerClient       *client.Client
	containers         map[string]*ContainerData
	eventStreamReading bool
	//lastUpdate           time.Time
	logsSavedDatePeriod int
	natsStreaming       ns.NatsStreaming
	nbLogs              int
	nbMetrics           int
}

//AgentInit Connect to docker engine, get initial containers list and start the agent
func AgentInit(version string, build string) error {
	agent := Agent{}
	agent.trapSignal()
	conf.init(version)

	//containers dir creqtion
	os.MkdirAll(containersDateDir, 0666)

	// NATS Connect
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln("Unable to get hostname:", err)
	}
	if agent.natsStreaming.Connect(amp.NatsDefaultURL, amp.NatsClusterID, os.Args[0]+"-"+hostname, amp.DefaultTimeout) != nil {
		return err
	}

	// Connection to Docker
	defaultHeaders := map[string]string{"User-Agent": "amp-agent"}
	cli, err := client.NewClient(conf.dockerEngine, amp.DockerDefaultVersion, nil, defaultHeaders)
	if err != nil {
		agent.natsStreaming.Close()
		return err
	}
	agent.dockerClient = cli
	fmt.Println("Connected to Docker-engine")

	fmt.Println("Extracting containers list...")
	agent.containers = make(map[string]*ContainerData)
	ContainerListOptions := types.ContainerListOptions{All: true}
	containers, err := agent.dockerClient.ContainerList(context.Background(), ContainerListOptions)
	if err != nil {
		agent.natsStreaming.Close()
		return err
	}
	for _, cont := range containers {
		agent.addContainer(cont.ID)
	}
	fmt.Println("done")
	agent.start()
	return nil
}

//Main agent loop, verify if events and logs stream are started if not start them
func (a *Agent) start() {
	a.initAPI()
	nb := 0
	for {
		a.uddateStreams()
		nb++
		if nb == 10 {
			fmt.Printf("Sent %d logs and %d metrics on the last %d seconds\n", a.nbLogs, a.nbMetrics, nb*conf.period)
			nb = 0
			a.nbLogs = 0
			a.nbMetrics = 0
		}
		time.Sleep(time.Duration(conf.period) * time.Second)
	}
}

//starts logs and metrics stream of eech new started container
func (a *Agent) uddateStreams() {
	a.updateLogsStream()
	a.updateMetricsStreams()
	a.updateEventsStream()
}

// Close AgentInit ressources
func (a *Agent) stop() {
	a.closeLogsStreams()
	a.closeMetricsStreams()
	a.dockerClient.Close()
	a.natsStreaming.Close()
}

// Launch a routine to catch SIGTERM Signal
func (a *Agent) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-agent received SIGTERM signal")
		a.stop()
		os.Exit(1)
	}()
}
