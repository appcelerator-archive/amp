package pilot

import (
	"os"
	"os/exec"
	"strconv"
)

const (
	defaultZookeeperHost = "zookeeper:2181"
)

// KafkaPilot struct
type KafkaPilot struct {
	cmd *exec.Cmd
}

// New instantiate a new KafkaPilot
func New() *KafkaPilot {
	return &KafkaPilot{}
}

// StartKafka start Kafka
func (pilot *KafkaPilot) StartKafka() error {
	pilot.cmd = exec.Command("/opt/kafka/bin/kafka-server-start.sh", "/opt/kafka/config/server.properties")
	pilot.cmd.Stdout = os.Stdout
	pilot.cmd.Stderr = os.Stderr
	err := pilot.cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// Wait wait for Kafka application to exit
func (pilot *KafkaPilot) Wait() error {
	return pilot.cmd.Wait()
}

// CreateTopic create a topic
func (pilot *KafkaPilot) CreateTopic(topic string, partitions uint64, replicationFactor uint64) error {
	cmd := exec.Command("/opt/kafka/bin/kafka-topics.sh", "--zookeeper", defaultZookeeperHost, "--topic", topic, "--create", "--partitions", strconv.FormatUint(partitions, 10), "--replication-factor", strconv.FormatUint(replicationFactor, 10))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// DeleteTopic delete a topic
func (pilot *KafkaPilot) DeleteTopic(topic string) error {
	cmd := exec.Command("/opt/kafka/bin/kafka-topics.sh", "--zookeeper", defaultZookeeperHost, "--topic", topic, "--delete")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
