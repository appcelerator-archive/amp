package amp

import "time"

const (
	// AmplifierDefaultEndpoint is the default amplifier endpoint
	AmplifierDefaultEndpoint = "amplifier:50101"

	// EtcdDefaultEndpoint is the default etcd endpoint
	EtcdDefaultEndpoint = "http://etcd:2379"

	// ElasticsearchDefaultURL is the default elasticsearch endpoint
	ElasticsearchDefaultURL = "http://elasticsearch:9200"

	// NatsDefaultURL is the default nats endpoint
	NatsDefaultURL = "nats://nats:4222"

	// NatsClusterID is the id of the infrastructure nats cluster
	NatsClusterID = "test-cluster"

	// DockerDefaultURL is the default docker endpoint
	DockerDefaultURL = "unix:///var/run/docker.sock"

	// DockerDefaultVersion is the default docker version
	DockerDefaultVersion = "1.24"

	// DefaultTimeout is the default timeout
	DefaultTimeout = time.Minute

	// NatsFunctionTopic is the topic used for function calls
	NatsFunctionTopic = "amp-function-calls"

	// NatsLogsTopic is the topic used for log events
	NatsLogsTopic = "amp-logs"

	// NatsMetricsTopic is the topic used for metrics events
	NatsMetricsTopic = "amp-metrics"

	// InfrastructureRole is the label used for infrastructure services
	InfrastructureRole = "infrastructure"

	// AmplifierDefaultAddress is the default amplifier endpoint
	AmplifierDefaultAddress = "local.appcelerator.io"

	//DefaultEmailSender email sender
	DefaultEmailSender = "amp.appcelerator@axway.com"

	//DefaultSmsSender sms sender
	DefaultSmsSender = ""
)
