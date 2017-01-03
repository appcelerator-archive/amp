package servercore

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/appcelerator/amp/cmd/adm-agent/agentgrpc"
	"github.com/appcelerator/amp/cmd/adm-server/servergrpc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	//DockerURL docker url
	DockerURL = "unix:///var/run/docker.sock"
	//DockerVersion docker version
	DockerVersion = "1.24"
	//ServiceFailedTimeout max time allowed to start a service
	ServiceFailedTimeout = 30 //seconds
)

const (
	blank     = "                                                                                         "
	separator = "-----------------------------------------------------------------------------------------"
)

// AMPInfraManager infra stack management
type AMPInfraManager struct {
	docker       *client.Client
	ctx          context.Context
	stack        *ampStack
	Silence      bool
	Verbose      bool
	Force        bool
	Local        bool
	Status       string
	clientID     string
	clientStream servergrpc.ClusterServerService_GetClientStreamServer
	server       *ClusterServer
	pullCount    int
}

type ampStack struct {
	networks        []string
	serviceMap      map[string]*ampService
	imageMap        map[string]string
	volumesToRemove []string
}

type ampService struct {
	readyToStart    bool
	starting        bool
	ready           bool
	failed          bool
	failedTime      time.Time
	forced          bool
	id              string
	name            string
	image           string
	desiredReplicas int
	dependencies    []string
	labels          map[string]string
	spec            *swarm.ServiceSpec
	//monitor data
	status          string
	containerOk     int
	containerFailed int
	user            bool
}

var currentColorTheme = "default"
var (
	colRegular = 0
	colInfo    = 1
	colWarn    = 2
	colError   = 3
	colSuccess = 4
	colUser    = 5
)

// Init initialize the struct
func (s *AMPInfraManager) Init(server *ClusterServer, clientID string, firstMessage string) error {
	if err := s.setStream(server, clientID); err != nil {
		return err
	}
	s.ctx = context.Background()
	if firstMessage != "" {
		s.printf(colRegular, firstMessage+"")
	}
	if s.Force {
		s.printf(colWarn, "Force mode: on")
	}
	defaultHeaders := map[string]string{"User-Agent": "amplifier"}
	cli, err := client.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("impossible to connect to Docker on: %s\n%v", DockerURL, err)
	}
	s.docker = cli
	return nil
}

func (s *AMPInfraManager) setStream(server *ClusterServer, clientID string) error {
	s.server = server
	s.clientID = clientID
	logf.debug("Set stream client id: %s\n", clientID)
	cli, ok := server.clientMap[clientID]
	if !ok {
		return fmt.Errorf("Client %s is not register", clientID)
	}
	s.clientStream = cli.stream
	return nil
}

// system prerequisites
func (s *AMPInfraManager) systemPrerequisites() error {
	sysctl := false
	// checks if GOOS is set
	goos := os.Getenv("GOOS")
	if goos == "linux" {
		sysctl = true
	} else if goos == "" {
		// check if sysctl exists on the system
		if _, err := os.Stat("/etc/sysctl.conf"); err == nil {
			sysctl = true
		}
	}
	if sysctl {
		var out bytes.Buffer
		var stderr bytes.Buffer
		mmcmin := 262144
		cmd := exec.Command("sysctl", "-n", "vm.max_map_count")
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		mmc, err := strconv.Atoi(strings.TrimRight(out.String(), "\n"))
		if err != nil {
			return err
		}
		if mmc < mmcmin {
			// admin rights are needed
			u, err := user.Current()
			if err != nil {
				return err
			}
			uid, err := strconv.Atoi(u.Uid)
			if err != nil {
				return err
			}
			if uid != 0 {
				return fmt.Errorf("vm.max_map_count should be at least 262144, admin rights are needed to update it")
			}
			if s.Verbose {
				s.printf(colRegular, "setting max virtual memory areas\n")
			}
			cmd = exec.Command("sysctl", "-w", "vm.max_map_count=262144")
			err = cmd.Run()
		} else if s.Verbose {
			s.printf(colRegular, "max virtual memory areas is already at a safe value\n")
		}
	}
	return nil
}

// Start start amp
func (s *AMPInfraManager) Start(stack *ampStack) error {
	s.stack = stack
	defer s.close()
	for _, name := range stack.networks {
		if err := s.createNetwork(name); err != nil {
			return err
		}
	}
	started := time.Now()
	for _, name := range stack.volumesToRemove {
		s.removeVolume(name)
		if time.Now().Sub(started) > time.Second*20 {
			s.printf(colWarn, "Impossible to delete volume %s", name)
			break
		}
	}
	for {
		serviceReady := 0
		serviceTotal := 0
		for _, service := range stack.serviceMap {
			if !service.user {
				serviceTotal++
				if err := s.updateServiceStates(service); err != nil {
					return err
				}
				if err := s.updateServiceDependencies(service); err != nil {
					return err
				}
				if !service.ready && !service.starting && service.readyToStart {
					if err := s.createService(service); err != nil {
						return err
					}
				}
				if service.ready {
					serviceReady++
				}
			}
		}
		if serviceReady == serviceTotal {
			s.printf(colRegular, "AMP platform started")
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

// Stop stop amp
func (s *AMPInfraManager) Stop(stack *ampStack) error {
	s.stack = stack
	s.addUserServices(stack)
	s.updateServiceInformation()
	for _, service := range stack.serviceMap {
		if _, ok := service.labels["io.amp.role"]; ok {
			if _, ok := s.doesServiceExist(service.name); !ok {
				s.printf(colInfo, "Service %s already stopped", service.name)
			} else {
				if err := s.removeService(service); err != nil {
					s.printf(colError, "Error stopping service %s: %v", service.name, err)
				}
			}
		}
	}
	s.printf(colRegular, "AMP platform stopped")
	return nil
}

// Pull pull all amp images
func (s *AMPInfraManager) Pull(agent *Agent, stack *ampStack) (int, int) {
	ok := 0
	ko := 0
	for _, image := range stack.imageMap {
		if s.Local && image == "appcelerator/amp:local" {
			continue
		}
		s.printf(colInfo, "Node %s: Pulling image %s", agent.nodeName, image)
		_, err := agent.client.PullImage(s.ctx, &agentgrpc.PullImageRequest{
			Image: image,
		})
		if err != nil {
			s.printf(colError, "Node %s: Image %s pull error: %v", agent.nodeName, image, err)
			ko++
		} else {
			ok++
		}
		s.printf(colInfo, "Node %s: Image %s pulled", agent.nodeName, image)
	}
	return ok, ko
}

// ComputeStatus compute amp status
func (s *AMPInfraManager) ComputeStatus(stack *ampStack) {
	s.stack = stack
	nbNotStarted := 0
	nbStarted := 0
	for name := range stack.serviceMap {
		if _, exist := s.doesServiceExist(name); exist {
			nbStarted++
		} else {
			nbNotStarted++
		}
	}
	if nbNotStarted == 0 {
		s.Status = "running"
	} else if nbStarted == 0 {
		s.Status = "stopped"
	} else {
		s.Status = "partially running"
	}
}

func (s *AMPInfraManager) addUserServices(stack *ampStack) {
	list, err := s.docker.ServiceList(s.ctx, types.ServiceListOptions{})
	if err != nil || len(list) == 0 {
		return
	}
	for _, serv := range list {
		if service, ok := stack.serviceMap[serv.Spec.Annotations.Name]; !ok {
			stack.serviceMap[serv.Spec.Annotations.Name] = &ampService{
				name:            serv.Spec.Annotations.Name,
				image:           serv.Spec.TaskTemplate.ContainerSpec.Image,
				desiredReplicas: s.getReplicas(serv.Spec),
				user:            true,
			}
		} else {
			service.desiredReplicas = s.getReplicas(serv.Spec)
		}
	}
}

func (s *AMPInfraManager) getReplicas(spec swarm.ServiceSpec) int {
	mode := spec.Mode
	if mode.Replicated != nil {
		return int(*mode.Replicated.Replicas)
	}
	return 0
}

// Monitor monitor amp services
func (s *AMPInfraManager) Monitor(stack *ampStack) (*[]*servergrpc.TypedOutput, error) {
	s.stack = stack
	s.Verbose = true
	cols := []int{14, 20, 10, 12, 10, 12}
	infraListName := []string{}
	for name := range stack.serviceMap {
		infraListName = append(infraListName, name)
	}
	sort.Strings(infraListName)

	s.addUserServices(stack)
	s.updateServiceInformation()
	for name, serv := range stack.serviceMap {
		if len(name) > cols[1] {
			cols[1] = len(name) + 2
		}
		if len(serv.status) > cols[2] {
			cols[2] = len(serv.status) + 2
		}
	}
	userListName := []string{}
	for name, serv := range stack.serviceMap {
		if serv.user {
			if name != "adm-server" && name != "adm-agent" {
				userListName = append(userListName, name)
			}
		}
	}
	sort.Strings(userListName)
	listName := []string{}
	if _, ok := stack.serviceMap["adm-server"]; ok {
		listName = append(listName, "adm-server")
	}
	if _, ok := stack.serviceMap["adm-agent"]; ok {
		listName = append(listName, "adm-agent")
	}

	listName = append(listName, infraListName...)
	listName = append(listName, userListName...)
	//fmt.Println("list: %v", listName)

	output := &[]*servergrpc.TypedOutput{}
	s.addOutput(output, colRegular, fmt.Sprintf("%s%s%s%s%s%s", col("ID", cols[0]), col("SERVICE", cols[1]), col("STATUS", cols[2]), col("MODE", cols[3]), col("REPLICAS", cols[4]), col("TASK FAILED", cols[5])))
	s.addOutput(output, colRegular, fmt.Sprintf("%s%s%s%s%s%s", col("-", cols[0]), col("-", cols[1]), col("-", cols[2]), col("-", cols[3]), col("-", cols[4]), col("-", cols[5])))
	for _, name := range listName {
		serv := stack.serviceMap[name]
		if serv.status == "running" {
			if serv.user && serv.name != "adm-server" && serv.name != "adm-agent" { //TODO remove specific name usage need
				s.addOutput(output, colUser, s.displayService(serv, cols))
			} else {
				s.addOutput(output, colSuccess, s.displayService(serv, cols))
			}
		} else if serv.status == "failing" {
			s.addOutput(output, colError, s.displayService(serv, cols))
		} else if serv.status == "partially running" {
			s.addOutput(output, colWarn, s.displayService(serv, cols))
		} else if serv.status == "starting" {
			s.addOutput(output, colRegular, s.displayService(serv, cols))
		} else {
			s.addOutput(output, colInfo, s.displayService(serv, cols))
		}
	}
	return output, nil
}

func (s *AMPInfraManager) addOutput(list *[]*servergrpc.TypedOutput, outputType int, output string) {
	*list = append(*list, &servergrpc.TypedOutput{OutputType: int32(outputType), Output: output})
}

func (s *AMPInfraManager) displayService(service *ampService, cols []int) string {
	var dispID string
	if service.id == "" {
		dispID = col("", cols[0])
	} else {
		dispID = col(service.id[0:12], cols[0])
	}
	disp := fmt.Sprintf("%s%s%s", dispID, col(service.name, cols[1]), col(service.status, cols[2]))
	if service.desiredReplicas == 0 {
		disp += col("global", cols[3])
		disp += col(fmt.Sprintf("%d", service.containerOk), cols[4])
	} else {
		disp += col("replicated", cols[3])
		disp += col(fmt.Sprintf("%d/%d", service.containerOk, service.desiredReplicas), cols[4])
	}
	disp += col(fmt.Sprintf("%d", service.containerFailed), cols[5])
	return disp
}

func (s *AMPInfraManager) updateServiceInformation() {
	for name, service := range s.stack.serviceMap {
		service.status = "stopped"
		service.id = ""
		service.containerOk = 0
		service.containerFailed = 0
		if _, exist := s.doesServiceExist(name); exist {
			serv, err := s.inspectService(service)
			if err != nil {
				service.status = "inspect error"
			} else {
				service.id = serv.ID
				service.status = "starting"
				taskList, err := s.getServiceTasks(service)
				if err != nil {
					service.status = "get task error"
				} else {
					for _, task := range taskList {
						if task.DesiredState == swarm.TaskStateRunning && task.Status.State == swarm.TaskStateRunning {
							service.containerOk++
						}
						if task.DesiredState == swarm.TaskStateShutdown || task.DesiredState == swarm.TaskStateFailed || task.DesiredState == swarm.TaskStateRejected {
							service.containerFailed++
						}
					}
					if service.containerOk > 0 {
						if service.containerOk == service.desiredReplicas || service.desiredReplicas == 0 {
							service.status = "running"
						} else {
							service.status = "partially running"
						}
					} else if service.containerFailed > 0 {
						service.status = "failing"
					}
				}
			}
		}
	}
}

func (s *AMPInfraManager) updateServiceStates(service *ampService) error {
	if service.id == "" || service.forced {
		return nil
	}
	taskList, err := s.getServiceTasks(service)
	if err != nil {
		return err
	}
	//Verify fist that there's at least a failed container
	for _, task := range taskList {
		if task.DesiredState == swarm.TaskStateRunning && task.Status.State == swarm.TaskStateRunning {
			if !service.ready {
				s.printf(colSuccess, "Service %s is ready", service.name)
			}
			service.ready = true
			service.failed = false
			return nil
		}
	}
	service.ready = false
	//no container running, then verify that there's at least a failed container
	for _, task := range taskList {
		if task.DesiredState == swarm.TaskStateShutdown || task.DesiredState == swarm.TaskStateFailed || task.DesiredState == swarm.TaskStateRejected {
			if !service.failed {
				service.failedTime = time.Now()
				s.printf(colWarn, "Warning: Service %s is failing", service.name)
			} else {
				if time.Now().Sub(service.failedTime) > ServiceFailedTimeout*time.Second {
					if s.Force {
						s.forceService(service)
						return nil
					}
					return fmt.Errorf("Service %s startup timeout", service.name)
				}
			}
			service.ready = false
			service.failed = true
			return nil
		}
	}
	service.failed = false
	//No container running and failing then service is still starting without error
	return nil
}

func (s *AMPInfraManager) forceService(service *ampService) {
	s.printf(colWarn, "Force mode: service %s is considered as ready", service.name)
	service.forced = true
	service.ready = true
	service.failed = true
}

func (s *AMPInfraManager) updateServiceDependencies(service *ampService) error {
	service.readyToStart = true
	for _, depName := range service.dependencies {
		serv, ok := s.stack.serviceMap[depName]
		if !ok {
			return fmt.Errorf("The dependency %s doesn't exist for the service %s", depName, service.name)
		}
		if !serv.ready {
			service.readyToStart = false
			return nil
		}
	}
	return nil
}

func (s *AMPInfraManager) getServiceTasks(service *ampService) ([]swarm.Task, error) {
	//filter := filters.NewArgs()
	//filter.Add("name", service.name)
	list, err := s.docker.TaskList(s.ctx, types.TaskListOptions{
	//Filter: filter,
	})
	if err != nil {
		return nil, err
	}
	taskList := []swarm.Task{}
	//get only service id task, name is not enough to discriminate
	for _, task := range list {
		if task.ServiceID == service.id {
			taskList = append(taskList, task)
		}
	}
	return taskList, nil
}

func (s *AMPInfraManager) inspectService(service *ampService) (swarm.Service, error) {
	serv, _, err := s.docker.ServiceInspectWithRaw(s.ctx, service.name)
	return serv, err
}

func (s *AMPInfraManager) createService(service *ampService) error {
	service.id = ""
	service.starting = false
	options := types.ServiceCreateOptions{}
	if len(service.dependencies) == 0 {
		s.printf(colInfo, "Starting service %s dependency: none", service.name)
	} else {
		dep := ""
		for _, depName := range service.dependencies {
			dep += depName + " "
		}
		if len(service.dependencies) == 1 {
			s.printf(colInfo, "Starting service %s dependency: %s", service.name, dep)
		} else {
			s.printf(colInfo, "Starting service %s dependencies: %s", service.name, dep)
		}
	}
	if _, ok := s.doesImageExist(service.image); !ok {
		if s.Force {
			s.printf(colWarn, "Service %s image %s doesn't exist", service.name, service.image)
			s.forceService(service)
			return nil
		}
		s.printf(colError, "Service %s image %s doesn't exist", service.name, service.image)
		return fmt.Errorf("Service image error")
	}
	r, err := s.docker.ServiceCreate(s.ctx, *service.spec, options)
	if err != nil {
		return err
	}
	service.starting = true
	service.id = r.ID
	return nil
}

// verify if service exist
func (s *AMPInfraManager) doesServiceExist(name string) (string, bool) {
	//filter := filters.NewArgs() //remove filter for docker version compatibility to be re-added later
	//filter.Add("name", name)
	list, err := s.docker.ServiceList(s.ctx, types.ServiceListOptions{
	//Filter: filter,
	})
	if err != nil || len(list) == 0 {
		return "", false
	}
	for _, serv := range list {
		if serv.Spec.Annotations.Name == name {
			return serv.ID, true
		}
	}
	return "", false
}

// verify if image exist locally
func (s *AMPInfraManager) doesImageExist(image string) (string, bool) {
	filter := filters.NewArgs()
	filter.Add("reference", image)
	filter.Add("dangling", "false")
	options := types.ImageListOptions{All: false, Filters: filter}

	list, err := s.docker.ImageList(s.ctx, options)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	if len(list) == 0 {
		return "", false
	}
	//fmt.Printf("images: %+v\n", list)
	for _, ima := range list {
		if len(ima.RepoTags) > 0 && ima.RepoTags[0] == image {
			return list[0].ID, true
		}
	}
	return "", false
}

// verify if network already exist
func (s *AMPInfraManager) doestNetworkExist(name string) (string, bool) {
	//filter := filters.NewArgs() //remove filter for docker version compatibility to be re-added later
	//filter.Add("name", name)
	list, err := s.docker.NetworkList(s.ctx, types.NetworkListOptions{
	//Filters: filter,
	})
	if err != nil || len(list) == 0 {
		return "", false
	}
	for _, net := range list {
		if net.Name == name {
			return net.ID, true
		}
	}
	return "", false
}

// create network
func (s *AMPInfraManager) createNetwork(name string) error {
	if _, exist := s.doestNetworkExist(name); exist {
		s.printf(colInfo, "Network %s already exist", name)
		return nil
	}
	s.printf(colInfo, "Create network %s", name)
	IPAM := network.IPAM{
		Driver:  "default",
		Options: make(map[string]string),
	}
	networkCreate := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		IPAM:           &IPAM,
	}
	_, err := s.docker.NetworkCreate(s.ctx, name, networkCreate)
	if err != nil {
		return err
	}
	return nil
}

func (s *AMPInfraManager) removeService(service *ampService) error {
	err := s.docker.ServiceRemove(s.ctx, service.name)
	if err == nil {
		s.printf(colSuccess, "Service %s removed", service.name)
	}
	return err

}

func (s *AMPInfraManager) removeVolume(name string) error {
	//s.printf(colorInfo, "Remove volume %s", name)
	//docker bug on remove volume, even using force a volume is not removed if used somewhere
	time.Sleep(3 * time.Second)
	return s.docker.VolumeRemove(s.ctx, name, true)
}

func (s *AMPInfraManager) cleanVolume(name string) error {
	s.printf(colInfo, "Remove volume %s", name)
	started := time.Now()
	filter := filters.NewArgs()
	filter.Add("name", name)
	for {
		s.docker.VolumeRemove(s.ctx, name, true)
		ret, err := s.docker.VolumeList(s.ctx, filter)
		if err == nil {
			return fmt.Errorf("Failed to get volume list: %v", err)
		}
		list := ret.Volumes
		ok := true
		for _, vol := range list {
			if vol.Name == name {
				ok = false
			}
		}
		if ok {
			return nil
		}
		if time.Now().Sub(started) > 30*time.Second {
			return fmt.Errorf("Timeout waiting for all services removed")
		}
		time.Sleep(1 * time.Second)
	}
}

// Perrorf print error
func (s *AMPInfraManager) Perrorf(format string, args ...interface{}) {
	s.printf(colError, format, args...)
}

// Pwarnf print warn
func (s *AMPInfraManager) Pwarnf(format string, args ...interface{}) {
	s.printf(colWarn, format, args...)
}

// Pinfof print info
func (s *AMPInfraManager) Pinfof(format string, args ...interface{}) {
	s.printf(colInfo, format, args...)
}

// Pregularf print regular
func (s *AMPInfraManager) Pregularf(format string, args ...interface{}) {
	s.printf(colRegular, format, args...)
}

// Psuccessf print success
func (s *AMPInfraManager) Psuccessf(col int, format string, args ...interface{}) {
	s.printf(colSuccess, format, args...)
}

// Puserf print user
func (s *AMPInfraManager) Puserf(col int, format string, args ...interface{}) {
	s.printf(colUser, format, args...)
}

func (s *AMPInfraManager) printf(col int, format string, args ...interface{}) {
	if s.Silence {
		return
	}
	if !s.Verbose && col == colInfo {
		return
	}
	output := &servergrpc.TypedOutput{
		Output:     fmt.Sprintf(format, args...),
		OutputType: int32(col),
	}
	s.clientStream.Send(&servergrpc.ClientMes{
		ClientId: s.clientID,
		Function: "Print",
		Output:   output,
	})
}

func (s *AMPInfraManager) close() {
	//s.docker.Close()
}

//ampStack & ampService

func (s *ampStack) init() {
	s.serviceMap = make(map[string]*ampService)
	s.imageMap = make(map[string]string)
}

func (s *ampStack) addImage(imageLabel string, image string) {
	s.imageMap[imageLabel] = image
}

func (s *ampStack) addService(m *AMPInfraManager, name string, imageLabel string, replicas int, spec *swarm.ServiceSpec, dependencies ...string) {
	image, ok := s.imageMap[imageLabel]
	if !ok {
		m.printf(colError, "The image doesn't exist for service: %s", name)
		return
	}
	spec.Annotations.Name = name
	spec.TaskTemplate.ContainerSpec.Image = image
	if replicas > 0 {
		nb := uint64(replicas)
		spec.Mode = swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &nb,
			},
		}
	} else {
		spec.Mode = swarm.ServiceMode{
			Global: &swarm.GlobalService{},
		}
	}
	labels := spec.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	s.serviceMap[name] = &ampService{
		name:            name,
		image:           image,
		desiredReplicas: replicas,
		spec:            spec,
		dependencies:    dependencies,
		labels:          labels,
	}
}

func (s *ampStack) getService(name string) *ampService {
	if service, ok := s.serviceMap[name]; ok {
		return service
	}
	return nil
}

// display value in the left of a col
func col(value string, size int) string {
	if value == "-" {
		return separator[0:size]
	}
	if len(value) > size {
		return value[0:size]
	}
	return value + blank[0:size-len(value)]
}

// display value in the right of a col
func colr(value string, size int) string {
	if len(value) > size {
		return value[0:size]
	}
	return blank[0:size-len(value)] + value
}

// display value in the middle of a col
func colm(value string, size int) string {
	if len(value) > size {
		return value[0:size]
	}
	space := size - len(value)
	rest := space % 2
	return blank[0:space/2+rest] + value + blank[0:space/2]
}

// display time col
func colTime(val int64, size int) string {
	tm := time.Unix(val, 0)
	value := tm.Format("2006-01-02 15:04:05")
	return col(value, size)
}
