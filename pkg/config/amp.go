package amp

import "time"

// AMP platform default values
const (
	AmplifierDefaultPort          = ":50101"
	AmplifierDefaultPublicAddress = "local.appcelerator.io"
	EmailDefaultSender            = "amp@atomiq.io"
	SmsDefaultSender              = "amp"

	EtcdDefaultEndpoint = "http://etcd:2379"

	ElasticsearchDefaultURL = "http://elasticsearch:9200"

	DockerDefaultURL     = "unix:///var/run/docker.sock"
	DockerDefaultVersion = "1.24"

	NatsDefaultURL    = "nats://nats:4222"
	NatsClusterID     = "test-cluster"
	NatsFunctionTopic = "amp-function-calls"
	NatsLogsTopic     = "amp-logs"
	NatsMetricsTopic  = "amp-metrics"

	DefaultTimeout     = time.Minute
	InfrastructureRole = "infrastructure"
)
