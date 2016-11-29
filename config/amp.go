package amp

const (
	// AmplifierDefaultEndpoint is the default amplifier endpoint
	AmplifierDefaultEndpoint = "amplifier:50101"

	// EtcdDefaultEndpoints is the default etcd endpoint
	EtcdDefaultEndpoint = "http://etcd:2379"

	// ElasticsearchDefaultURL is the default elasticsearch endpoint
	ElasticsearchDefaultURL = "http://elasticsearch:9200"

	// NatsDefaultURL is the default nats endpoint
	NatsDefaultURL = "nats://nats:4222"

	// NatsClusterID is the id of the infrastructure nats cluster
	NatsClusterID = "test-cluster"

	// InfluxDefaultURL is the default influxdb endpoint
	InfluxDefaultURL = "http://influxdb:8086"

	// DockerDefaultURL is the default docker endpoint
	DockerDefaultURL = "unix:///var/run/docker.sock"

	// DockerDefaultVersion is the default docker version
	DockerDefaultVersion = "1.24"
)
