package ampadmin

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

const (
	defaultURL        = "unix:///var/run/docker.sock"
	defaultVersion    = "1.29"
	minimumApiVersion = 1.29
	testNetwork       = "amptest"
	TestPassed        = "passed"
	TestFailed        = "failed"
)

type testServiceSpec struct {
	Name        string
	Image       string
	Command     []string
	Networks    []string
	Replicas    int
	Constraints []string
}

func VerifyDockerVersion() (string, error) {
	c, err := client.NewClient(defaultURL, defaultVersion, nil, nil)
	if err != nil {
		return TestFailed, err
	}
	version, err := c.ServerVersion(context.Background())
	log.Printf("Docker engine version %s\n", version.Version)
	apiVersion, err := strconv.ParseFloat(version.APIVersion, 32)
	if err != nil {
		return TestFailed, err
	}
	if apiVersion < minimumApiVersion {
		log.Printf("minimum expected: %.3g, observed: %.3g", minimumApiVersion, apiVersion)
		return TestFailed, errors.New("Docker engine doesn't meet the requirements (API Version)")
	}
	return fmt.Sprintf("%s - api version = %.3g\n", TestPassed, apiVersion), nil
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

func VerifyServiceScheduling() (string, error) {
	c, err := client.NewClient(defaultURL, defaultVersion, nil, nil)
	if err != nil {
		return TestFailed, err
	}

	nwId, err := createNetwork(c, testNetwork)
	if err != nil {
		return TestFailed, err
	}
	defer func() {
		log.Printf("removing network %s (%s)\n", testNetwork, nwId)
		if err := c.NetworkRemove(context.Background(), nwId); err != nil {
			log.Printf(fmt.Sprintf("network deletion failed: %s\n", err))
		}
	}()

	serverServiceId, err := createService(c, testServiceSpec{Name: "check-server", Image: "alpine:3.6", Command: []string{"nc", "-kvlp", "5968", "-e", "echo"}, Networks: []string{testNetwork}, Replicas: 3, Constraints: []string{"node.labels.amp.type.api==true"}})
	if err != nil {
		return TestFailed, err
	}
	defer func() {
		log.Printf("removing service %s (%s)\n", "check-server", serverServiceId)
		_ = c.ServiceRemove(context.Background(), serverServiceId)
		time.Sleep(2 * time.Second)
	}()
	time.Sleep(5 * time.Second)
	if err := assertServiceHasRunningTasks(c, "check-server", 3); err != nil {
		return TestFailed, err
	}
	clientServiceId, err := createService(c, testServiceSpec{Name: "check-client", Image: "alpine:3.6", Command: []string{"sh", "-c", "while true; do nc -zv check-server 5968; done"}, Networks: []string{testNetwork}, Replicas: 3, Constraints: []string{"node.labels.amp.type.core==true"}})
	if err != nil {
		return TestFailed, err
	}
	defer func() {
		log.Printf("removing service %s (%s)\n", "check-client", clientServiceId)
		_ = c.ServiceRemove(context.Background(), clientServiceId)
		time.Sleep(2 * time.Second)
	}()
	time.Sleep(5 * time.Second)
	if err := assertServiceHasRunningTasks(c, "check-client", 3); err != nil {
		return TestFailed, err
	}
	time.Sleep(5 * time.Second)
	log.Println("Counting request success rate")
	body, err := c.ServiceLogs(context.Background(), clientServiceId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return TestFailed, err
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
			return TestFailed, err
		}
		if matched {
			openCount++
		}
	}
	if lineCount < 50 {
		log.Printf("%d connections / %d success\n", lineCount, openCount)
		return TestFailed, errors.New("Connection test failed, expected more connections")
	}
	if openCount < (lineCount - 10) {
		log.Printf("%d connections / %d success\n", lineCount, openCount)
		return TestFailed, errors.New("Connection test failed, not enough successes")
	}
	return fmt.Sprintf("%s - %d connections / %d success\n", TestPassed, lineCount, openCount), nil
}
