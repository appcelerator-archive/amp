package core

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"docker.io/go-docker/api/types"
	"github.com/appcelerator/amp/api/rpc/logs"
	"github.com/appcelerator/amp/cmd/amplifier/server/configuration"
	"github.com/appcelerator/amp/data/storage"
	"github.com/appcelerator/amp/data/storage/etcd"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/appcelerator/amp/pkg/nats-streaming"
	"golang.org/x/net/context"
)

const (
	containersDataDir = "/containers"
)

// Agent data
type Agent struct {
	dock                *docker.Docker
	containers          map[string]*ContainerData
	eventStreamReading  bool
	logsSavedDatePeriod int
	natsStreaming       *ns.NatsStreaming
	nbMetrics           int
	nbMetricsComputed   int
	nbLogs              int
	logsBuffer          *logs.GetReply
	logsBufferMutex     *sync.Mutex
	first10Min          int //send metrics every 10 sec, the first 10 min and then use the setting value
	storage             storage.Interface
}

// AgentInit Connect to docker engine, get initial containers list and start the agent
func AgentInit(version, build string) error {
	agent := Agent{
		logsSavedDatePeriod: 60,
		logsBuffer:          &logs.GetReply{},
		logsBufferMutex:     &sync.Mutex{},
	}
	agent.trapSignal()
	conf.init(version, build)

	// containers dir creation
	if err := os.MkdirAll(containersDataDir, 0666); err != nil {
		return fmt.Errorf("Unable to create container data directory: %s", err)
	}

	// NATS Connect
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("Unable to get hostname: %s", err)
	}
	agent.natsStreaming = ns.NewClient(ns.DefaultURL, ns.ClusterID, os.Args[0]+"-"+hostname, time.Minute)
	if err = agent.natsStreaming.Connect(); err != nil {
		return err
	}

	// Connection to storage
	etcdEndpoints := []string{etcd.DefaultEndpoint}
	agent.storage = etcd.New(etcdEndpoints, "amp", configuration.DefaultTimeout)
	if err := agent.storage.Connect(); err != nil {
		return err
	}

	// Connection to Docker
	agent.dock = docker.NewClient(conf.dockerEngine, docker.DefaultVersion)
	if err = agent.dock.Connect(); err != nil {
		_ = agent.natsStreaming.Close()
		return err
	}
	log.Infoln("Connected to Docker-engine")

	agent.containers = make(map[string]*ContainerData)
	ContainerListOptions := types.ContainerListOptions{All: true}
	containers, err := agent.dock.GetClient().ContainerList(context.Background(), ContainerListOptions)
	if err != nil {
		_ = agent.natsStreaming.Close()
		return err
	}
	for _, cont := range containers {
		agent.addContainer(cont.ID)
	}
	log.Infoln("Successfully retrieved containers.")
	agent.start()
	return nil
}

// Main agent loop, starts buffers thread if any and looks for new containers added or removed (according to docker events)
func (a *Agent) start() {
	a.initAPI()
	nb := 0
	//start a thread looking for the Metrics Buffer to send it if full or period time reached
	a.startMetricsSender()
	//start a thread looking for the Logs Buffer to send it if full or period time reached
	a.startLogsBufferSender()
	for {
		//looks for new containers to add or to remove, for each added, opem a stream for logs and a stream for metrics feeding the buffers
		a.updateStreams()
		nb++
		if nb == 10 {
			log.Infof("Sent %d logs and %d metrics (%d computed) on the last %d seconds\n", a.nbLogs, a.nbMetrics, a.nbMetricsComputed, nb*conf.period)
			nb = 0
			a.nbLogs = 0
			a.nbMetrics = 0
			a.nbMetricsComputed = 0
		}
		time.Sleep(time.Duration(conf.period) * time.Second)
	}
}

//start a thread looking for the Metrics Buffer to send it if full or period time reached, if no buffer is set, then do nothing than initialize the buffer to one element
func (a *Agent) startMetricsSender() {
	a.first10Min = 60
	go func() {
		for {
			if conf.metricsPeriod > 10 && a.first10Min > 0 {
				a.first10Min--
				time.Sleep(10 * time.Second)
			} else {
				time.Sleep(time.Second * time.Duration(conf.metricsPeriod))
			}
			a.sendSquashedMetricsMessages()
		}
	}()
}

//start a thread looking for the Logs Buffer to send it if full or period time reached, if no buffer is set, then do nothing than initialize the buffer to one element
func (a *Agent) startLogsBufferSender() {
	if conf.logsBufferPeriod == 0 || conf.logsBufferSize == 0 {
		a.logsBuffer.Entries = make([]*logs.LogEntry, 1)
		return
	}
	go func() {
		time.Sleep(time.Second * time.Duration(conf.logsBufferPeriod))
		if len(a.logsBuffer.Entries) > 0 {
			a.logsBufferMutex.Lock()
			a.sendLogsBuffer()
			a.logsBufferMutex.Unlock()
		}
	}()
}

// Starts logs and metrics stream of eech new started container
func (a *Agent) updateStreams() {
	a.updateLogsStream()
	a.updateMetricsStreams()
	a.updateEventsStream()
}

// Close AgentInit resources
func (a *Agent) stop(status int) {
	a.closeLogsStreams()
	a.closeMetricsStreams()
	if err := a.dock.GetClient().Close(); err != nil {
		log.Errorf("Docker api close error: %v\n", err)
	}

	if err := a.natsStreaming.Close(); err != nil {
		log.Errorf("Nats close error: %v\n", err)
	}
	os.Exit(status)
}

// Launch a routine to catch SIGTERM Signal
func (a *Agent) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("\nagent received SIGTERM signal")
		a.stop(1)
	}()
}
