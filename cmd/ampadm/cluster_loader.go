package main

import (
	"fmt"
	"github.com/appcelerator/amp/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"golang.org/x/net/context"
	"time"
)

const (
	//DockerURL docker url
	DockerURL = amp.DockerDefaultURL
	//DockerVersion docker version
	DockerVersion = amp.DockerDefaultVersion
	//RegistryToken token used for registry
	RegistryToken = ""
)

// ClusterLoader to load cluster service the right way
type ClusterLoader struct {
	docker  *docker.Client
	client  *clusterClient
	ctx     context.Context
	silence bool
	verbose bool
	force   bool
	local   bool
	ampTag  string
}

func (s *ClusterLoader) init(client *clusterClient, firstMessage string) error {
	s.ctx = context.Background()
	s.client = client
	defaultHeaders := map[string]string{"User-Agent": "adm-client"}
	cli, err := docker.NewClient(DockerURL, DockerVersion, nil, defaultHeaders)
	if err != nil {
		return fmt.Errorf("impossible to connect to Docker on: %s\n%v", DockerURL, err)
	}
	s.docker = cli
	return nil
}

func (s *ClusterLoader) startClusterServices() error {
	stack := clusterStack{}
	stack.init(s.local, s.ampTag)
	for _, name := range stack.networks {
		if err := s.createNetwork(name); err != nil {
			return err
		}
	}
	starting := false
	for _, service := range stack.serviceMap {
		if _, ok := s.doesServiceExist(service.name); ok {
			client.printfc(colInfo, "Service %s is already started\n", service.name)
		} else {
			if err := s.createService(service); err != nil {
				client.printfc(colError, "Service %s create error: %v\n", service.name, err)
			} else {
				if s.verbose {
					client.printfc(colSuccess, "Service %s starting using image %s\n", service.name, service.image)
				} else {
					client.printfc(colSuccess, "Service %s starting\n", service.name)
				}
				starting = true
			}
		}
	}
	t0 := time.Now()
	for {
		nb := 0
		for _, service := range stack.serviceMap {
			if s.isServiceRunning(service.name) {
				nb++
			}
		}
		if nb == 2 {
			break
		}
		if time.Now().Sub(t0).Seconds() > 30 {
			client.printfc(colError, "Cluster services startup timeout\n")
			s.stopClusterServices()
			return nil
		}
	}
	if starting {
		client.printfc(colSuccess, "Cluster services started\n")
	}
	return nil
}

func (s *ClusterLoader) stopClusterServices() error {
	stack := clusterStack{}
	stack.init(s.local, s.ampTag)
	for _, service := range stack.serviceMap {
		id, exist := s.doesServiceExist(service.name)
		if !exist {
			s.client.printfc(colSuccess, "Service %s stopped\n", service.name)
		} else {
			if err := s.removeService(id); err != nil {
				client.printfc(colError, "Service %s remove error: %v\n", service.name, err)
			} else {
				client.printfc(colSuccess, "Service %s removed\n", service.name)
			}
		}
	}
	return nil
}

// verify if service exist
func (s *ClusterLoader) doesServiceExist(name string) (string, bool) {
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

func (s *ClusterLoader) createService(service *clusterService) error {
	options := types.ServiceCreateOptions{}
	_, err := s.docker.ServiceCreate(s.ctx, *service.spec, options)
	if err != nil {
		return err
	}
	return nil
}

// verify if network already exist
func (s *ClusterLoader) doestNetworkExist(name string) (string, bool) {
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
func (s *ClusterLoader) createNetwork(name string) error {
	if _, exist := s.doestNetworkExist(name); exist {
		return nil
	}
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

func (s *ClusterLoader) removeService(id string) error {
	err := s.docker.ServiceRemove(s.ctx, id)
	if err != nil {
		return err
	}
	return nil

}

func (s *ClusterLoader) isServiceRunning(name string) bool {
	status := s.getServiceStatus(name, 1)
	if status == "running" {
		return true
	}
	return false
}

func (s *ClusterLoader) getServiceStatus(name string, replicas int) string {
	status := "stopped"
	containerOk := 0
	containerFailed := 0
	if _, exist := s.doesServiceExist(name); exist {
		serv, err := s.inspectService(name)
		if err != nil {
			status = "inspect error"
		} else {
			id := serv.ID
			status = "starting"
			taskList, err := s.getServiceTasks(id)
			if err != nil {
				status = "get task error"
			} else {
				for _, task := range taskList {
					if task.DesiredState == swarm.TaskStateRunning && task.Status.State == swarm.TaskStateRunning {
						containerOk++
					}
					if task.DesiredState == swarm.TaskStateShutdown || task.DesiredState == swarm.TaskStateFailed || task.DesiredState == swarm.TaskStateRejected {
						containerFailed++
					}
				}
				if containerOk > 0 {
					if containerOk >= 1 || replicas == 0 {
						status = "running"
					}
				} else if containerFailed > 0 {
					status = "failing"
				}
			}
		}
	}
	return status
}

func (s *ClusterLoader) inspectService(name string) (swarm.Service, error) {
	serv, _, err := s.docker.ServiceInspectWithRaw(s.ctx, name)
	return serv, err
}

func (s *ClusterLoader) getServiceTasks(id string) ([]swarm.Task, error) {
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
		if task.ServiceID == id {
			taskList = append(taskList, task)
		}
	}
	return taskList, nil
}
