package admin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	sk "github.com/appcelerator/amp/cluster/agent/swarm"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/docker/client"
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
)

const (
	DefaultURL        = "unix:///var/run/docker.sock"
	DefaultVersion    = "1.30"
	minimumApiVersion = "1.30"
	testNetwork       = "amptest"
)

type testServiceSpec struct {
	Name        string
	Image       string
	Command     []string
	Networks    []string
	Replicas    int
	Constraints []string
}

func VerifyDockerVersion() error {
	c, err := client.NewClient(DefaultURL, DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	version, err := c.ServerVersion(context.Background())
	apiVersion := version.APIVersion
	if versions.LessThan(apiVersion, minimumApiVersion) {
		log.Printf("Docker engine version %s\n", version.Version)
		log.Printf("API version - minimum expected: %.s, observed: %.s", minimumApiVersion, apiVersion)
		return errors.New("Docker engine doesn't meet the requirements (API Version)")
	}
	return nil
}

func VerifyLabels() error {
	labels := map[string]bool{}
	expectedLabels := []string{"amp.type.api=true", "amp.type.route=true", "amp.type.core=true", "amp.type.metrics=true",
		"amp.type.search=true", "amp.type.mq=true", "amp.type.kv=true", "amp.type.user=true"}
	missingLabel := false
	c, err := client.NewClient(DefaultURL, DefaultVersion, nil, nil)
	if err != nil {
		return err
	}
	nodes, err := c.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return err
	}
	// get the full list of labels
	for _, node := range nodes {
		nodeLabels := node.Spec.Annotations.Labels
		for k, v := range nodeLabels {
			labels[fmt.Sprintf("%s=%s", k, v)] = true
		}
	}
	// check that all expected labels are at least on one node
	for _, label := range expectedLabels {
		if !labels[label] {
			log.Printf("label %s is missing\n", label)
			missingLabel = true
		}
	}
	if missingLabel {
		return errors.New("At least one missing label")
	}
	return nil

}

func createNetwork(c *client.Client, name string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("name", name)
	res, err := c.NetworkList(context.Background(), types.NetworkListOptions{Filters: filter})
	if err != nil {
		return "", err
	}
	if len(res) == 1 {
		log.Printf("Network %s already exists\n", name)
		return res[0].ID, nil
	}
	log.Printf("creating network %s\n", name)
	nw, err := c.NetworkCreate(context.Background(), name, types.NetworkCreate{Driver: "overlay", Attachable: true})
	if err != nil {
		return "", err
	}
	return nw.ID, nil
}

func createService(c *client.Client, spec testServiceSpec) (string, error) {
	var networkAttachments []swarm.NetworkAttachmentConfig
	for _, n := range spec.Networks {
		networkAttachments = append(networkAttachments, swarm.NetworkAttachmentConfig{Target: n})
	}
	placement := swarm.Placement{Constraints: spec.Constraints}
	task := swarm.TaskSpec{
		ContainerSpec: swarm.ContainerSpec{
			Image:   spec.Image,
			Command: spec.Command,
		},
		Placement: &placement,
		Networks:  networkAttachments,
	}
	replicas := uint64(spec.Replicas)
	log.Printf("creating service %s\n", spec.Name)
	resp, err := c.ServiceCreate(context.Background(), swarm.ServiceSpec{Annotations: swarm.Annotations{Name: spec.Name}, Mode: swarm.ServiceMode{Replicated: &swarm.ReplicatedService{Replicas: &replicas}}, TaskTemplate: task}, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func assertServiceHasRunningTasks(c *client.Client, name string, count int) error {
	filter := filters.NewArgs()
	filter.Add("service", name)
	filter.Add("desired-state", "running")
	tasks, err := c.TaskList(context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return err
	}
	if len(tasks) != count {
		return errors.New(fmt.Sprintf("%d running task for service %s, expected %d", len(tasks), name, count))
	}
	log.Printf("assertServiceHasRunningTasks(%s) passed\n", name)
	return nil
}

// listen for events and write them in the channel
func listenSwarmEvents(eventChan chan *api.WatchMessage_Event, w api.Watch_WatchClient) {
	// until we receive the first empty message, the events should be considered as garbage
	dirty := true
	for {
		msg, err := w.Recv()
		if err == io.EOF {
			return
		}
		//&WatchMessage{Events:[&WatchMessage_Event{Action:WATCH_ACTION_CREATE,...
		if err != nil {
			log.Printf("Error while receiving events: %s\n", err)
			return
		}
		events := msg.Events
		if len(events) == 0 {
			if dirty {
				// Initial event
				dirty = false
				continue
			}
			log.Println("Error: received an extra empty event")
			return
		}
		if !dirty {
			for _, event := range events {
				//log.Printf("Action: %s\n", event.Action.String())
				eventChan <- event
			}
		}
	}
}

// returns true if the expected event count has been caught in the channel
func waitForEvents(eventChan chan *api.WatchMessage_Event, expectedEvent string, expectedCount int, seconds int) bool {
	count := 0
	// TODO: what if seconds==0
	timeout := time.After(time.Duration(seconds) * time.Second)
	for {
		select {
		case event := <-eventChan:
			if event.Action.String() == expectedEvent {
				count++
			}
			if expectedCount == count {
				log.Println("expected event count reached")
				return true
			}
			if expectedCount < count {
				log.Printf("expected event count over reached (%d/%d)\n", count, expectedCount)
				return true
			}
		case <-timeout:
			// log.Printf("timeout reached, count = %d/%d\n", count, expectedCount)
			log.Printf("timeout reached while waiting for %s events", expectedEvent)
			return false
		}
	}
}

func eventWatcher(eventType string) (chan *api.WatchMessage_Event, error) {
	// listen for swarm events on tasks
	_, conn, err := sk.Dial(sk.DefaultSocket())
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Println("Error: ", s)
		}
		return nil, err
	}
	watcher := api.NewWatchClient(conn)
	watchEntry := sk.NewWatchRequestEntry(eventType, sk.WatchActionKindAll, nil)
	watchEntries := []*api.WatchRequest_WatchEntry{
		watchEntry,
	}
	ctx := context.TODO()
	in := sk.NewWatchRequest(watchEntries, nil, true)
	w, err := watcher.Watch(ctx, in)
	if err != nil {
		return nil, err
	}
	// buffered channel for Swarm events
	eventChan := make(chan *api.WatchMessage_Event, 32)
	go listenSwarmEvents(eventChan, w)
	return eventChan, nil
}

func VerifyServiceScheduling() error {
	c, err := client.NewClient(DefaultURL, DefaultVersion, nil, nil)
	if err != nil {
		return err
	}

	nwId, err := createNetwork(c, testNetwork)
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("removing network %s (%s)\n", testNetwork, nwId)
		if err := c.NetworkRemove(context.Background(), nwId); err != nil {
			log.Printf("network deletion failed: %s\n", err)
		}
	}()

	// listening for task events
	eventChan, err := eventWatcher("task")
	if err != nil {
		return err
	}
	serverServiceId, err := createService(c, testServiceSpec{Name: "check-server", Image: "alpine:3.6", Command: []string{"nc", "-kvlp", "5968", "-e", "echo"}, Networks: []string{testNetwork}, Replicas: 3, Constraints: []string{"node.labels.amp.type.api==true"}})
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("Removing service %s (%s)\n", "check-server", serverServiceId)
		_ = c.ServiceRemove(context.Background(), serverServiceId)
		time.Sleep(2 * time.Second)
	}()
	// look for task creation events
	if observed := waitForEvents(eventChan, "WATCH_ACTION_CREATE", 3, 10); !observed {
		return errors.New("failed to read the server task creation events")
	} else {
		log.Println("Task creation events successfully read")
	}
	clientServiceId, err := createService(c, testServiceSpec{Name: "check-client", Image: "alpine:3.6", Command: []string{"sh", "-c", "while true; do nc -zv check-server 5968; done"}, Networks: []string{testNetwork}, Replicas: 3, Constraints: []string{"node.labels.amp.type.core==true"}})
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("Removing service %s (%s)\n", "check-client", clientServiceId)
		_ = c.ServiceRemove(context.Background(), clientServiceId)
		time.Sleep(2 * time.Second)
	}()
	if observed := waitForEvents(eventChan, "WATCH_ACTION_CREATE", 3, 10); !observed {
		return errors.New("failed to read the client task creation events")
	} else {
		log.Println("Task creation events successfully read")
	}
	// wait 6 seconds to make sure no tasks are dropped
	if dropped := waitForEvents(eventChan, "WATCH_ACTION_REMOVE", 1, 6); dropped {
		return errors.New("tasks have been dropped")
	} else {
		log.Println("No dropped task")
	}
	time.Sleep(5 * time.Second)
	log.Println("Counting request success rate")
	body, err := c.ServiceLogs(context.Background(), clientServiceId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}
	defer body.Close()
	scanner := bufio.NewScanner(body)
	var lineCount int
	var openCount int
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		matched, err := regexp.MatchString(".*open$", line)
		if err != nil {
			return err
		}
		if matched {
			openCount++
		}
	}
	if lineCount < 50 {
		log.Printf("%d connections / %d success\n", lineCount, openCount)
		return errors.New("Connection test failed, expected more connections")
	}
	if openCount < (lineCount - 10) {
		log.Printf("%d connections / %d success\n", lineCount, openCount)
		return errors.New("Connection test failed, not enough successes")
	}
	return nil
}
