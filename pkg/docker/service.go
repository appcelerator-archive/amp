package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/service/progress"
	"github.com/appcelerator/amp/docker/docker/pkg/jsonmessage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type ServiceStatus struct {
	RunningTasks int32
	FailedTasks  int32
	TotalTasks   int32
	Status       string
}

// ServiceInspect inspects a service
func (d *Docker) ServiceInspect(ctx context.Context, service string) (swarm.Service, error) {
	serviceEntity, _, err := d.client.ServiceInspectWithRaw(ctx, service, types.ServiceInspectOptions{InsertDefaults: true})
	if err != nil {
		return swarm.Service{}, err
	}
	return serviceEntity, nil
}

// ServiceScale scales a service
func (d *Docker) ServiceScale(ctx context.Context, service string, scale uint64) error {
	serviceEntity, err := d.ServiceInspect(ctx, service)
	if err != nil {
		return err
	}
	serviceMode := &serviceEntity.Spec.Mode
	if serviceMode.Replicated == nil {
		return fmt.Errorf("scale can only be used with replicated mode")
	}
	serviceMode.Replicated.Replicas = &scale

	response, err := d.client.ServiceUpdate(ctx, serviceEntity.ID, serviceEntity.Version, serviceEntity.Spec, types.ServiceUpdateOptions{})
	if err != nil {
		return err
	}
	for _, warning := range response.Warnings {
		log.Warnln(warning)
	}
	log.Infof("service %s scaled to %d\n", service, scale)
	return nil
}

// ServiceStatus returns service status
func (d *Docker) ServiceStatus(ctx context.Context, service string) (*ServiceStatus, error) {
	taskMap, err := d.checkTasks(ctx, service)
	if err != nil {
		return &ServiceStatus{}, err
	}
	totalTasks, err := d.ExpectedNumberOfTasks(ctx, service)
	if err != nil {
		return &ServiceStatus{}, err
	}
	if totalTasks == 0 {
		return &ServiceStatus{
			RunningTasks: 0,
			TotalTasks:   0,
			FailedTasks:  0,
			Status:       StateNoMatchingNode,
		}, nil
	}
	if taskMap[StateError] != 0 && taskMap[StateRunning] != totalTasks {
		return &ServiceStatus{
			RunningTasks: int32(taskMap[StateRunning]),
			TotalTasks:   int32(totalTasks),
			FailedTasks:  int32(taskMap[StateError]),
			Status:       StateError,
		}, nil
	}
	if taskMap[StateRunning] == totalTasks {
		return &ServiceStatus{
			RunningTasks: int32(taskMap[StateRunning]),
			TotalTasks:   int32(totalTasks),
			FailedTasks:  int32(taskMap[StateError]),
			Status:       StateRunning,
		}, nil
	}
	return &ServiceStatus{
		RunningTasks: int32(taskMap[StateRunning]),
		TotalTasks:   int32(totalTasks),
		FailedTasks:  int32(taskMap[StateError]),
		Status:       StateStarting,
	}, nil
}

// ServiceList list the services
func (d *Docker) ServicesList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	return d.client.ServiceList(ctx, options)
}

// WaitOnService waits for the service to converge. It outputs a progress bar,
func (d *Docker) WaitOnService(ctx context.Context, serviceID string, quiet bool) error {
	errChan := make(chan error, 1)
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		errChan <- progress.ServiceProgress(ctx, d.client, serviceID, pipeWriter)
	}()

	if quiet {
		go io.Copy(ioutil.Discard, pipeReader)
		return <-errChan
	}

	err := jsonmessage.DisplayJSONMessagesToStream(pipeReader, command.NewOutStream(os.Stdout), nil)
	if err == nil {
		err = <-errChan
	}
	return err
}
