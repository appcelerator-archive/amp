package stack

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

type ampService struct {
	id              string
	name            string
	image           string
	starting        int
	ready           int
	failing         int
	desiredReplicas int
	labels          map[string]string
}

// load all running services in s, create a ampService instance for each
func (s *Server) setCurrentServices(ctx context.Context) error {
	s.serviceMap = make(map[string]*ampService)
	if err := s.addUserServices(ctx); err != nil {
		return err
	}
	s.updateServiceInformation(ctx)
	return nil
}

// return the list of all infrastructure stacks needed images
func (s *Server) getImages(local bool) ([]string, error) {
	fileName := StackFileVarName
	if !local {
		fileName = fmt.Sprintf("%s/%s", stackFilePath, fileName)
	}
	mapVar, err := LoadInfraVariables(fileName, "")
	if err != nil {
		return nil, err
	}
	imageList := []string{}
	for name, val := range mapVar {
		if strings.HasPrefix(name, "image") {
			imageList = append(imageList, val)
		}
	}
	return imageList, nil
}

// WaitForServiceReady wait a given service to be reday or timedout
func (s *Server) WaitForServiceReady(ctx context.Context, serviceName string, timeoutSec int) error {
	s.setCurrentServices(ctx)
	t0 := time.Now()
	for {
		id, exist := s.DoesServiceExist(ctx, serviceName)
		if exist {
			ready, _, _ := s.getServiceStatus(ctx, id)
			if ready > 0 {
				time.Sleep(1 * time.Second)
				return nil
			}
		}
		time.Sleep(time.Second * 1)
		if time.Now().Sub(t0).Nanoseconds()/1000000000 > int64(timeoutSec) {
			return fmt.Errorf("timeout")
		}
	}
}

//pull an image
func (s *Server) pullImage(ctx context.Context, image string) error {
	options := types.ImagePullOptions{}
	options.RegistryAuth = "" //TODO: handle this
	reader, err := s.Docker.ImagePull(ctx, image, options)
	if err != nil {
		return fmt.Errorf("image %s pull error: %v", image, err)
	}
	data := make([]byte, 1000, 1000)
	for {
		_, err := reader.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("Pull image %s error: %v", image, err)
		}
	}
	return nil
}

// verify if service exist
func (s *Server) addUserServices(ctx context.Context) error {
	list, err := s.Docker.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil || len(list) == 0 {
		return err
	}
	for _, serv := range list {
		if service, ok := s.serviceMap[serv.Spec.Annotations.Name]; !ok {
			s.serviceMap[serv.Spec.Annotations.Name] = &ampService{
				id:              serv.ID,
				name:            serv.Spec.Annotations.Name,
				image:           serv.Spec.TaskTemplate.ContainerSpec.Image,
				desiredReplicas: s.getReplicas(serv.Spec),
				labels:          serv.Spec.Labels,
			}
		} else {
			service.desiredReplicas = s.getReplicas(serv.Spec)
		}
	}
	return nil
}

// extract replicas information from swarm serviceSPec struct
func (s *Server) getReplicas(spec swarm.ServiceSpec) int {
	mode := spec.Mode
	if mode.Replicated != nil {
		return int(*mode.Replicated.Replicas)
	}
	return 0
}

// GetMonitorLines build the list of text line to return to the monitoring command
func (s *Server) GetMonitorLines(ctx context.Context) ([]*MonitorService, error) {
	if err := s.setCurrentServices(ctx); err != nil {
		return nil, err
	}
	s.addUserServices(ctx)
	s.updateServiceInformation(ctx)
	listName := []string{}
	for name := range s.serviceMap {
		listName = append(listName, name)
	}
	sort.Strings(listName)
	ret := []*MonitorService{}
	for _, name := range listName {
		serv := s.serviceMap[name]
		ret = append(ret, s.getMonitorLine(serv))
	}
	return ret, nil
}

//build one text line of the monitoring command for a given service
func (s *Server) getMonitorLine(service *ampService) *MonitorService {
	stackName := service.labels["com.docker.stack.namespace"]
	line := &MonitorService{
		Id:         service.id[0:12],
		Stack:      stackName,
		Service:    service.name,
		Status:     s.getServiceStatusString(service.ready, service.starting, service.failing),
		Mode:       "replicated",
		Replicas:   fmt.Sprintf("%d/%d", service.ready, service.desiredReplicas),
		FailedTask: fmt.Sprintf("%d", service.failing),
	}
	if service.desiredReplicas == 0 {
		line.Mode = "global"
	}
	return line
}

//update all the services status
func (s *Server) updateServiceInformation(ctx context.Context) {
	for name, service := range s.serviceMap {
		if id, exist := s.DoesServiceExist(ctx, name); exist {
			ready, starting, failling := s.getServiceStatus(ctx, id)
			service.ready = ready
			service.starting = starting
			service.failing = failling
		}
	}
}

// DoesServiceExist verify if service exist, return id if so
func (s *Server) DoesServiceExist(ctx context.Context, name string) (string, bool) {
	list, err := s.Docker.ServiceList(ctx, types.ServiceListOptions{
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

// return service status: ready, starting, failling
func (s *Server) getServiceStatus(ctx context.Context, id string) (ready int, starting int, failling int) {
	taskList, err := s.getServiceTasks(ctx, id)
	if err != nil {
		failling++
		return
	}
	//Verify fist that there's at least a failed container
	for _, task := range taskList {
		if task.DesiredState == swarm.TaskStateRunning && task.Status.State == swarm.TaskStateRunning {
			ready++
		} else if task.DesiredState == swarm.TaskStateShutdown || task.DesiredState == swarm.TaskStateFailed || task.DesiredState == swarm.TaskStateRejected {
			starting++
			failling++
		} else {
			starting++
		}
	}
	return
}

//return service tasks
func (s *Server) getServiceTasks(ctx context.Context, id string) ([]swarm.Task, error) {
	list, err := s.Docker.TaskList(ctx, types.TaskListOptions{})
	if err != nil {
		return nil, err
	}
	taskList := []swarm.Task{}
	//get only service id task, name is not enough to discriminate
	for _, task := range list {
		if task.ServiceID == id {
			taskList = append(taskList, task)
		}
	}
	return taskList, nil
}

//return a string value of a service status
func (s *Server) getServiceStatusString(ready int, starting int, failing int) string {
	if ready > 0 {
		return "running"
	} else if failing > 0 {
		return "failing"
	} else if starting > 0 {
		return "starting"
	}
	return "stopped"
}

// EvalMappingString parse a value of io.amp.mapping label
// used in stack only for evaluation: verify the syntax is ok
// used in /cmd/haproxy to get the result values of the parsing.
func EvalMappingString(mapping string) ([]string, error) {
	ret := make([]string, 4, 4)
	if mapping == "" {
		return ret, fmt.Errorf("mapping is empty")
	}
	data := strings.Split(mapping, ":")
	if len(data) < 2 {
		return ret, fmt.Errorf("mapping format error should be: label:port[:tcp:port]: %s", mapping)
	}
	ret[0] = data[0]
	if _, err := strconv.Atoi(data[1]); err != nil {
		return ret, fmt.Errorf("mapping format error, port should be a number: %s", mapping)
	}
	ret[1] = data[1]
	if len(data) == 3 {
		return ret, fmt.Errorf("mapping format error should be: label:port[:tcp:port]: %s", mapping)
	}
	if len(data) >= 4 {
		if data[2] != "tcp" {
			return ret, fmt.Errorf("mapping format error mode can be only 'tcp' to specify a grpc mapping: %s", mapping)
		}
		ret[2] = data[2]
		if _, err := strconv.Atoi(data[3]); err != nil {
			return ret, fmt.Errorf("mapping format error, internal tcp port should be a number: %s", mapping)
		}
		ret[3] = data[3]

	}
	return ret, nil
}
