package main

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

	"github.com/appcelerator/amp/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"golang.org/x/net/context"
)

const (
	//DockerURL docker url
	DockerURL = amp.DockerDefaultURL
	//DockerVersion docker version
	DockerVersion = amp.DockerDefaultVersion
	//RegistryToken token used for registry
	RegistryToken = ""
	//ServiceFailedTimeout max time allowed to start a service
	ServiceFailedTimeout = 30 //seconds
)

type ampManager struct {
	docker     *client.Client
	ctx        context.Context
	stack      *ampStack
	silence    bool
	verbose    bool
	force      bool
	local      bool
	status     string
	printColor [6]*color.Color
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

func (s *ampManager) init(firstMessage string) error {
	s.ctx = context.Background()
	s.setColors()
	if firstMessage != "" {
		s.printf(colRegular, firstMessage+"\n")
	}
	if s.force {
		s.printf(colWarn, "Force mode: on\n")
	}
	defaultHeaders := map[string]string{"User-Agent": "amplifier"}
	cli, err := client.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("impossible to connect to Docker on: %s\n%v", DockerURL, err)
	}
	s.docker = cli
	return nil
}

func (s *ampManager) start(stack *ampStack) error {
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
			s.printf(colWarn, "Impossible to delete volume %s\n", name)
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
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *ampManager) stop(stack *ampStack) error {
	s.stack = stack
	s.addUserServices(stack)
	s.updateServiceInformation()
	for _, service := range stack.serviceMap {
		if _, ok := s.doesServiceExist(service.name); !ok {
			s.printf(colInfo, "Service %s already stopped\n", service.name)
		} else {
			if err := s.removeService(service); err != nil {
				s.printf(colError, "Error stopping service %s: %v\n", service.name, err)
			}
		}
	}
	return nil
}

func (s *ampManager) pull(stack *ampStack) error {
	for _, image := range stack.imageMap {
		if s.local && image == "appcelerator/amp:local" {
			continue
		}
		s.printf(colInfo, "Pulling image %s\n", image)
		options := types.ImagePullOptions{}
		if RegistryToken != "" {
			options.RegistryAuth = RegistryToken
		}
		reader, err := s.docker.ImagePull(s.ctx, image, options)
		if err != nil {
			s.printf(colError, "image %s pull error: %v\n", image, err)
		}
		data := make([]byte, 1000, 1000)
		for {
			_, err := reader.Read(data)
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				s.printf(colWarn, "Pull image %s error: %v\n", image, err)
			}
			//TODO: good display
			//s.printf(colorInfo, "%s", string(data))
			s.printf(colInfo, "+")
		}
		s.printf(colInfo, "\n")
		s.printf(colSuccess, "Image %s pulled\n", image)

	}
	return nil
}

func (s *ampManager) computeStatus(stack *ampStack) {
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
		s.status = "running"
	} else if nbStarted == 0 {
		s.status = "stopped"
	} else {
		s.status = "partially running"
	}
}

// verify if service exist
func (s *ampManager) addUserServices(stack *ampStack) {
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

func (s *ampManager) getReplicas(spec swarm.ServiceSpec) int {
	mode := spec.Mode
	if mode.Replicated != nil {
		return int(*mode.Replicated.Replicas)
	}
	return 0
}

func (s *ampManager) monitor(stack *ampStack) {
	s.stack = stack
	s.verbose = true
	cols := []int{14, 20, 10, 12, 10, 12}
	fmt.Println("\033[2J\033[0;0H")
	infraListName := []string{}
	for name := range stack.serviceMap {
		infraListName = append(infraListName, name)
	}
	sort.Strings(infraListName)
	for {
		s.addUserServices(stack)
		s.updateServiceInformation()
		fmt.Println("\033[0;0H")
		for _, serv := range stack.serviceMap {
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
				userListName = append(userListName, name)
			}
		}
		sort.Strings(userListName)
		listName := []string{}
		listName = append(listName, infraListName...)
		listName = append(listName, userListName...)
		//fmt.Println("list: %v\n", listName)
		fmt.Printf("%s%s%s%s%s%s\n", col("ID", cols[0]), col("SERVICE", cols[1]), col("STATUS", cols[2]), col("MODE", cols[3]), col("REPLICAS", cols[4]), col("TASK FAILED", cols[5]))
		fmt.Printf("%s%s%s%s%s%s\n", col("-", cols[0]), col("-", cols[1]), col("-", cols[2]), col("-", cols[3]), col("-", cols[4]), col("-", cols[5]))
		for _, name := range listName {
			serv := stack.serviceMap[name]
			if serv.status == "running" {
				if serv.user {
					s.printf(colUser, "%s\n", s.displayService(serv, cols))
				} else {
					s.printf(colSuccess, "%s\n", s.displayService(serv, cols))
				}
			} else if serv.status == "failing" {
				s.printf(colError, "%s\n", s.displayService(serv, cols))
			} else if serv.status == "partially running" {
				s.printf(colWarn, "%s\n", s.displayService(serv, cols))
			} else if serv.status == "starting" {
				s.printf(colRegular, "%s\n", s.displayService(serv, cols))
			} else {
				s.printf(colInfo, "%s\n", s.displayService(serv, cols))
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *ampManager) displayService(service *ampService, cols []int) string {
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

func (s *ampManager) updateServiceInformation() {
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

func (s *ampManager) updateServiceStates(service *ampService) error {
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
				s.printf(colSuccess, "Service %s is ready\n", service.name)
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
				s.printf(colWarn, "Warning: Service %s is failing\n", service.name)
			} else {
				if time.Now().Sub(service.failedTime) > ServiceFailedTimeout*time.Second {
					if s.force {
						s.forceService(service)
						return nil
					}
					return fmt.Errorf("service %s startup timeout", service.name)
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

func (s *ampManager) forceService(service *ampService) {
	s.printf(colWarn, "Force mode: service %s is considered as ready\n", service.name)
	service.forced = true
	service.ready = true
	service.failed = true
}

func (s *ampManager) updateServiceDependencies(service *ampService) error {
	service.readyToStart = true
	for _, depName := range service.dependencies {
		serv, ok := s.stack.serviceMap[depName]
		if !ok {
			return fmt.Errorf("the dependency %s doesn't exist for the service %s", depName, service.name)
		}
		if !serv.ready {
			service.readyToStart = false
			return nil
		}
	}
	return nil
}

func (s *ampManager) getServiceTasks(service *ampService) ([]swarm.Task, error) {
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

func (s *ampManager) inspectService(service *ampService) (swarm.Service, error) {
	serv, _, err := s.docker.ServiceInspectWithRaw(s.ctx, service.name)
	return serv, err
}

func (s *ampManager) createService(service *ampService) error {
	service.id = ""
	service.starting = false
	options := types.ServiceCreateOptions{}
	if len(service.dependencies) == 0 {
		s.printf(colInfo, "Starting service %s dependency: none\n", service.name)
	} else {
		dep := ""
		for _, depName := range service.dependencies {
			dep += depName + " "
		}
		if len(service.dependencies) == 1 {
			s.printf(colInfo, "Starting service %s dependency: %s\n", service.name, dep)
		} else {
			s.printf(colInfo, "Starting service %s dependencies: %s\n", service.name, dep)
		}
	}
	if _, ok := s.doesImageExist(service.image); !ok {
		if s.force {
			s.printf(colWarn, "Service %s image %s doesn't exist\n", service.name, service.image)
			s.forceService(service)
			return nil
		}
		return fmt.Errorf("service %s image %s doesn't exist", service.name, service.image)
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
func (s *ampManager) doesServiceExist(name string) (string, bool) {
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
func (s *ampManager) doesImageExist(image string) (string, bool) {
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
func (s *ampManager) doestNetworkExist(name string) (string, bool) {
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
func (s *ampManager) createNetwork(name string) error {
	if _, exist := s.doestNetworkExist(name); exist {
		s.printf(colInfo, "Network %s already exist\n", name)
		return nil
	}
	s.printf(colInfo, "Create network %s\n", name)
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

func (s *ampManager) removeService(service *ampService) error {
	err := s.docker.ServiceRemove(s.ctx, service.name)
	if err == nil {
		s.printf(colSuccess, "Service %s removed\n", service.name)
	}
	return err

}

func (s *ampManager) removeVolume(name string) error {
	//s.printf(colorInfo, "Remove volume %s\n", name)
	//docker bug on remove volume, even using force a volume is not removed if used somewhere
	time.Sleep(3 * time.Second)
	return s.docker.VolumeRemove(s.ctx, name, true)
}

func (s *ampManager) cleanVolume(name string) error {
	s.printf(colInfo, "Remove volume %s\n", name)
	started := time.Now()
	filter := filters.NewArgs()
	filter.Add("name", name)
	for {
		s.docker.VolumeRemove(s.ctx, name, true)
		ret, err := s.docker.VolumeList(s.ctx, filter)
		if err == nil {
			return fmt.Errorf("failed to get volume list: %v", err)
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
			return fmt.Errorf("timeout waiting for all services removed")
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *ampManager) printf(col int, format string, args ...interface{}) {
	if s.silence {
		return
	}
	colorp := s.printColor[0]
	if col > 0 && col < len(s.printColor) {
		colorp = s.printColor[col]
	}
	if !s.verbose && col == colInfo {
		return
	}
	colorp.Printf(format, args...)
}

func (s *ampManager) close() {
	//s.docker.Close()
}

func (s *ampManager) setColors() {
	//theme := s.getTheme()
	theme := AMP.Configuration.CmdTheme
	if theme == "dark" {
		s.printColor[0] = color.New(color.FgHiWhite)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
		s.printColor[5] = color.New(color.FgHiGreen)
	} else {
		s.printColor[0] = color.New(color.FgMagenta)
		s.printColor[1] = color.New(color.FgHiBlack)
		s.printColor[2] = color.New(color.FgYellow)
		s.printColor[3] = color.New(color.FgRed)
		s.printColor[4] = color.New(color.FgGreen)
		s.printColor[5] = color.New(color.FgHiGreen)
	}
	//add theme as you want.
}

// system prerequisites
func (m *ampManager) systemPrerequisites() error {
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
			if m.verbose {
				m.printf(colRegular, "setting max virtual memory areas\n")
			}
			cmd = exec.Command("sysctl", "-w", "vm.max_map_count=262144")
			err = cmd.Run()
		} else if m.verbose {
			m.printf(colRegular, "max virtual memory areas is already at a safe value\n")
		}
	}
	return nil
}

//ampStack & ampService

func (s *ampStack) init() {
	s.serviceMap = make(map[string]*ampService)
	s.imageMap = make(map[string]string)
}

func (s *ampStack) addImage(imageLabel string, image string) {
	s.imageMap[imageLabel] = image
}

func (s *ampStack) addService(m *ampManager, name string, imageLabel string, replicas int, spec *swarm.ServiceSpec, dependencies ...string) {
	image, ok := s.imageMap[imageLabel]
	if !ok {
		m.printf(colError, "The image doesn't exist for service: %s\n", name)
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

	s.serviceMap[name] = &ampService{
		name:            name,
		image:           image,
		desiredReplicas: replicas,
		spec:            spec,
		dependencies:    dependencies,
	}
}

func (s *ampStack) getService(name string) *ampService {
	if service, ok := s.serviceMap[name]; ok {
		return service
	}
	return nil
}
