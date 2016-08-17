package stat

import (
	"encoding/json"
	"testing"

	"github.com/influxdata/influxdb/client/v2"
)

// Static Message for testing query results without influx
var (
	jsonMsg = []byte(`{"Results":[{"Series":[{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/amp-agent:latest","container_name":"amp-agent.0.a81rrghki2zby1swoc9sbvywn","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.4041869944273426,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/amp-log-worker:latest","container_name":"amp-log-worker.1.5bfvgpi8jtzh0pfkxtxalinar","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,1.5470597424366765,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/amp-log:latest","container_name":"amp-log.1.12j0y22kfkl6swg0wl8af4e2e","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.40712588333910166,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/amp-monitor:latest","container_name":"amp-monitor.1.abenfc3hniaf6dn270fpkg21w","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.4117169313689853,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/amp-ui:latest","container_name":"amp-ui.1.a3hwhv5hrjqmzosqzxddiej8u","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.30526519242175015,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/elasticsearch-amp:latest","container_name":"elasticsearch.1.7nhfi2gru0zeot3yvgcxag0us","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,8.572736784489871,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/grafana:latest","container_name":"grafana.1.9db7qe4lsq2m50g61n04ddxzj","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,1.045158087713783,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/influxdb:latest","container_name":"influxdb.1.0g6yoxuppjwaikiu9tycq3659","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,1.065678697077215,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/kafka:latest","container_name":"kafka.1.7n4j4m9dpmlt3zyqglllqijc9","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,2.019878996559292,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/kapacitor:latest","container_name":"kapacitor.1.b52urzhfe8inugqawwifo8fi1","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.6679577655186435,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/kibana:latest","container_name":"kibana.1.39hcg1piog4ghhhaxu072wohj","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.6120147356036376,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/nginx:latest","container_name":"nginx.1.a91dwdx7qex6evutm35kz9d38","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.3163646345055874,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/telegraf:latest","container_name":"telegraf-agent.0.32g08w0fdytkxj9kfs08hdcrp","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,1.591462501319946,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/telegraf:latest","container_name":"telemetry-worker.1.7enpz8qii5gs4qbx2z5ll741s","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,1.1040008410284445,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"appcelerator/zookeeper:latest","container_name":"zookeeper.1.0jlx509me7c3l1ddpfa8n7jo0","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,3.773735265682202,3.811409808474576e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"quay.io/coreos/etcd:v3.0.4","container_name":"etcd.1.42rj42ta8oyi9hswhov8b7r7o","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.40696223882120486,3.811010375862069e+14]]},{"name":"docker_container_cpu","tags":{"container_image":"sheepkiller/kafka-manager:latest","container_name":"kafka-manager.1.6c1rmg30i6t5lx0h91inqraaw","datacenter":"dc1"},"columns":["time","cpu_pct","usage_system"],"values":[[1471534751,0.5569247096606933,3.811409808474576e+14]]}],"Messages":null}]}`)
)

func TestMarshalCPU(t *testing.T) {
	res := &client.Response{}
	// Marshal JSON to Influx Type to simulate db access
	err := json.Unmarshal(jsonMsg, &res)
	if err != nil {
		t.Error(err)
	}
	if len(res.Results) == 0 {
		t.Errorf("The Results Length should be greater than zero %v\n", res)
	}
	size := len(res.Results[0].Series)
	if size == 0 {
		t.Errorf("The Results Length should be greater than zero %v\n", res)
	}
	// Create a new Marshaller for the Protobuf Query
	marsh := NewCPUContainerMarshaler(size)
	// Marshal the Influx Type to Proto
	MarshalInfluxToProto(res, marsh)
	if len(marsh.cpus) != size {
		t.Errorf("The Protobuf message size %v does not match expected size of %v\n", len(marsh.cpus), size)
	}

}
