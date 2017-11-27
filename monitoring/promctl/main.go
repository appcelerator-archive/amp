package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"github.com/appcelerator/amp/pkg/docker"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const (
	defaultTemplate      = "/etc/prometheus/prometheus.tpl"
	defaultConfiguration = "/etc/prometheus/prometheus.yml"
	defaultHost          = docker.DefaultURL
	defaultPeriod        = 1
	dockerForMacIP       = "192.168.65.1"
	prometheusCmd        = "/bin/prometheus"
	monitoringNetwork    = "monit"
	stackName            = "amp"
	metricsPortLabel     = "io.amp.metrics.port"
	metricsPathLabel     = "io.amp.metrics.path"
	metricsModeLabel     = "io.amp.metrics.mode"
	metricsModeTasks     = "tasks"
	metricsModeExporter  = "exporter"
	externalURLEnv       = "PROMETHEUS_EXTERNAL_URL"
	externalURLOption    = "--web.external-url"
)

var prometheusArgs = []string{
	"--config.file=/etc/prometheus/prometheus.yml",
	"--storage.tsdb.path=/prometheus",
	"--storage.tsdb.retention=15d",
	"--web.console.libraries=/usr/share/prometheus/console_libraries",
	"--web.console.templates=/usr/share/prometheus/consoles",
	"--storage.tsdb.no-lockfile",
}

type Inventory struct {
	Jobs []Job
}

type Job struct {
	Name           string
	Mode           string
	StaticConfigs  []StaticConfig
	RelabelConfigs []RelabelConfig
	MetricsPath    string
}

// static config for a prometheus job
type StaticConfig struct {
	Target string
	Port   int
	Labels map[string]string
}

type RelabelConfig struct {
	SourceLabels []string
	Separator    string
	TargetLabel  string
	Replacement  string
}

type Target struct {
	Name        string
	Port        int
	MetricsPath string
}

// discovers services with the io.amp.metrics.port label
// get the name and host IP of the tasks of the services
func prepareJobs(client *docker.Docker, networkResource types.NetworkResource) ([]Job, error) {
	var jobs []Job
	filter := filters.NewArgs()
	filter.Add("label", metricsPortLabel)
	// only available on manager nodes
	services, err := client.GetClient().ServiceList(context.Background(), types.ServiceListOptions{Filters: filter})
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		fmt.Println("Warning: no service discovered for monitoring")
	}
	for _, service := range services {
		name := service.Spec.Annotations.Name
		strMetricsPort, ok := service.Spec.Annotations.Labels[metricsPortLabel]
		if !ok {
			fmt.Printf("Warning: unable to get metrics port label for service %s, ignoring it\n", name)
			continue
		}
		metricsPort, err := strconv.Atoi(strMetricsPort)
		if err != nil {
			log.Printf("Warning: non numerical port for service %s: %s\n", name, strMetricsPort)
			continue
		}
		metricsPath, ok := service.Spec.Annotations.Labels[metricsPathLabel]
		if !ok {
			metricsPath = "/metrics"
		}
		// mode can be tasks or exporter
		metricsMode, ok := service.Spec.Annotations.Labels[metricsModeLabel]
		if !ok {
			metricsMode = metricsModeTasks
		}
		fmt.Printf("discovered service %s on port %d and path %s, mode %s\n", name, metricsPort, metricsPath, metricsMode)
		s, ok := networkResource.Services[name]
		if !ok {
			fmt.Printf("Warning: service %s not found in network %s, ignoring it\n", name, monitoringNetwork)
			continue
		}
		switch metricsMode {
		default:
			fmt.Printf("Warning: wrong metrics mode (%s) for service %s, force it to %s\n", metricsMode, name, metricsModeTasks)
			metricsMode = metricsModeTasks
			fallthrough
		case metricsModeTasks:
			if len(s.Tasks) == 0 {
				continue
			}
			job := Job{Name: strings.TrimPrefix(name, fmt.Sprintf("%s_", stackName)), Mode: metricsModeTasks, MetricsPath: metricsPath}
			for _, task := range s.Tasks {
				job.StaticConfigs = append(job.StaticConfigs, StaticConfig{
					Target: task.EndpointIP,
					Port:   metricsPort,
					Labels: map[string]string{
						"hostip":   task.Info["Host IP"],
						"taskname": task.Name,
					},
				})
			}
			// all "tasks" jobs have the same relabel config
			job.RelabelConfigs = append(job.RelabelConfigs,
				RelabelConfig{SourceLabels: []string{"hostip"}, Separator: "@", TargetLabel: "instance"})
			jobs = append(jobs, job)
		case metricsModeExporter:
			shortName := strings.TrimSuffix(strings.TrimPrefix(name, fmt.Sprintf("%s_", stackName)), "_exporter")
			job := Job{Name: shortName, Mode: metricsModeExporter, MetricsPath: metricsPath}
			job.StaticConfigs = []StaticConfig{
				{
					Target: name,
					Port:   metricsPort,
				},
			}
			job.RelabelConfigs = append(job.RelabelConfigs,
				RelabelConfig{Replacement: shortName, TargetLabel: "instance"})
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func update(client *docker.Docker, configurationTemplate string, configuration string) error {
	var configurationFile *os.File
	// connect to the engine API
	if err := client.Connect(); err != nil {
		return err
	}
	filter := filters.NewArgs()
	filter.Add("name", monitoringNetwork)
	networkResources, err := client.GetClient().NetworkList(context.Background(), types.NetworkListOptions{Filters: filter})
	if err != nil {
		return err
	}
	if len(networkResources) != 1 {
		return errors.New("network lookup failed")
	}
	networkId := networkResources[0].ID
	networkResource, err := client.GetClient().NetworkInspect(context.Background(), networkId, types.NetworkInspectOptions{Verbose: true})

	jobs, err := prepareJobs(client, networkResource)
	if err != nil {
		return err
	}

	inventory := &Inventory{Jobs: jobs}
	// prepare the configuration
	t := template.Must(template.New("prometheus.tpl").Funcs(template.FuncMap{"StringsJoin": strings.Join}).ParseFiles(configurationTemplate))
	configurationFile, err = os.Create(configuration)
	if err != nil {
		return err
	}
	err = t.Execute(configurationFile, inventory)
	if err != nil {
		return err
	}
	configurationFile.Close()

	// reload prometheus
	cmd := exec.Command("/usr/bin/killall", "-HUP", "prometheus")
	err = cmd.Run()
	if err != nil {
		log.Println("Prometheus reload failed, error message follows")
		return err
	}
	return nil
}

// Docker client init
func dockerClientInit(host string) (*docker.Docker, error) {
	hasAScheme, err := regexp.MatchString(".*://.*", host)
	if err != nil {
		return nil, err
	}
	if !hasAScheme {
		host = "tcp://" + host
	}
	hasAPort, err := regexp.MatchString(".*(:[0-9]+|sock)", host)
	if err != nil {
		return nil, err
	}
	if !hasAPort {
		host = host + ":2375"
	}
	client := docker.NewClient(host, docker.DefaultVersion)
	return client, nil
}

// am I a manager?
func isAManager(client *docker.Docker) (bool, error) {
	if err := client.Connect(); err != nil {
		return false, err
	}
	info, err := client.GetClient().Info(context.Background())
	if err != nil {
		return false, err
	}
	nodeId := info.Swarm.NodeID
	for _, peer := range info.Swarm.RemoteManagers {
		if peer.NodeID == nodeId {
			return true, nil
		}
	}
	return false, nil
}

func main() {
	var configuration string
	var configurationTemplate string
	var host string
	var period int32

	var RootCmd = &cobra.Command{
		Use:   "promctl",
		Short: "Prometheus controller",
		Long:  `Keep the Prometheus configuration up to date with swarm discovery`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := dockerClientInit(host)
			if err != nil {
				return err
			}
			manager, err := isAManager(client)
			if !manager || err != nil {
				return fmt.Errorf("service discovery requires a connection to a manager engine socket")
			}
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
			tick := time.Tick(time.Duration(period) * time.Minute)
			time.Sleep(5 * time.Second)
			if err := update(client, configurationTemplate, configuration); err != nil {
				return err
			}

		loop:
			for {
				select {
				case <-tick:
					if err := update(client, configurationTemplate, configuration); err != nil {
						fmt.Println(err.Error())
					}
				case sig := <-stop:
					log.Printf("%v signal trapped\n", sig)
					break loop
				}
			}
			log.Println("Stopping Prometheus")
			stopCmd := exec.Command("/usr/bin/killall", "prometheus")
			if err := stopCmd.Run(); err != nil {
				log.Println("unable to stop Prometheus, error message follows")
				return err
			}
			return nil
		},
	}

	RootCmd.PersistentFlags().StringVarP(&configuration, "config", "c", defaultConfiguration, "config file")
	RootCmd.PersistentFlags().StringVarP(&configurationTemplate, "template", "t", defaultTemplate, "template file")
	RootCmd.PersistentFlags().StringVar(&host, "host", defaultHost, "host")
	RootCmd.PersistentFlags().Int32VarP(&period, "period", "p", defaultPeriod, "reload period in minute")

	// Set Prometheus external URL if provided
	if url := os.Getenv(externalURLEnv); url != "" {
		prometheusArgs = append(prometheusArgs, fmt.Sprintf("%s=%s", externalURLOption, url))
	}
	// start Prometheus
	proc := exec.Command(prometheusCmd, prometheusArgs...)
	stdout, err := proc.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}
	outscanner := bufio.NewScanner(stdout)
	errscanner := bufio.NewScanner(stderr)
	go func() {
		for outscanner.Scan() {
			fmt.Println(outscanner.Text())
		}
	}()
	go func() {
		for errscanner.Scan() {
			fmt.Println(errscanner.Text())
		}
	}()

	go func() {
		err := proc.Start()
		if err != nil {
			log.Fatalln(err)
		}
		err = proc.Wait()
		log.Println("Prometheus has exited")
		if err != nil {
			log.Fatalln(err)
		}
		// TODO: process is terminated without error, we should probably also stop the root cmd
	}()
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
