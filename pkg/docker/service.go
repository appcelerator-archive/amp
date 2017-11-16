package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/filters"
	"docker.io/go-docker/api/types/swarm"
	"github.com/appcelerator/amp/docker/cli/cli/command"
	"github.com/appcelerator/amp/docker/cli/cli/service/progress"
	"github.com/appcelerator/amp/docker/docker/pkg/jsonmessage"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type ServiceStatus struct {
	RunningTasks   int32
	CompletedTasks int32
	FailedTasks    int32
	TotalTasks     int32
	Status         string
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
func (d *Docker) ServiceStatus(ctx context.Context, service *swarm.Service) (*ServiceStatus, error) {
	// Get expected number of tasks for the service
	expectedTaskCount, err := d.ExpectedNumberOfTasks(ctx, service.ID)
	if err != nil {
		return nil, err
	}
	if expectedTaskCount == 0 {
		return &ServiceStatus{
			RunningTasks: 0,
			TotalTasks:   0,
			Status:       StateNoMatchingNode,
		}, nil
	}

	// List all tasks of service
	args := filters.NewArgs()
	args.Add("service", service.ID)
	tasks, err := d.TaskList(ctx, types.TaskListOptions{Filters: args})
	if err != nil {
		return nil, err
	}

	// Sort tasks by slot, then by most recent
	sort.Stable(TasksBySlot(tasks))

	// Build a map with only the most recent task per slot
	mostRecentTasks := map[int]swarm.Task{}
	for _, task := range tasks {
		if _, exists := mostRecentTasks[task.Slot]; !exists {
			mostRecentTasks[task.Slot] = task
		}
	}

	// Computing service status based on task status
	taskMap := map[string]int32{}
	for _, task := range mostRecentTasks {
		switch task.Status.State {
		case swarm.TaskStatePreparing:
			taskMap[StatePreparing]++
		case swarm.TaskStateReady:
			taskMap[StateReady]++
		case swarm.TaskStateStarting:
			taskMap[StateStarting]++
		case swarm.TaskStateRunning:
			taskMap[StateRunning]++
		case swarm.TaskStateComplete:
			taskMap[StateComplete]++
		case swarm.TaskStateFailed, swarm.TaskStateRejected:
			taskMap[StateError]++
		}
	}

	// If any task has an ERROR status, the service status is ERROR
	if taskMap[StateError] > 0 {
		return &ServiceStatus{
			RunningTasks: taskMap[StateRunning],
			TotalTasks:   expectedTaskCount,
			Status:       StateError,
		}, nil
	}

	// If all tasks are PREPARING, the service status is PREPARING
	if taskMap[StatePreparing] == expectedTaskCount {
		return &ServiceStatus{
			RunningTasks: taskMap[StateRunning],
			TotalTasks:   expectedTaskCount,
			Status:       StatePreparing,
		}, nil
	}

	// If all tasks are READY, the service status is READY
	if taskMap[StateReady] == expectedTaskCount {
		return &ServiceStatus{
			RunningTasks: taskMap[StateRunning],
			TotalTasks:   expectedTaskCount,
			Status:       StateReady,
		}, nil
	}

	// If all tasks are RUNNING, the service status is RUNNING
	if taskMap[StateRunning] == expectedTaskCount {
		return &ServiceStatus{
			RunningTasks: taskMap[StateRunning],
			TotalTasks:   expectedTaskCount,
			Status:       StateRunning,
		}, nil
	}

	// If all tasks are COMPLETE, the service status is COMPLETE
	if taskMap[StateComplete] == expectedTaskCount {
		return &ServiceStatus{
			RunningTasks: taskMap[StateRunning],
			TotalTasks:   expectedTaskCount,
			Status:       StateComplete,
		}, nil
	}

	// Else the service status is STARTING
	return &ServiceStatus{
		RunningTasks: taskMap[StateRunning],
		TotalTasks:   expectedTaskCount,
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
