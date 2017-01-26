# Telegraf Configuration
#
# Telegraf is entirely plugin driven. All metrics are gathered from the
# declared inputs, and sent to the declared outputs.
#
# Plugins must be declared in here to be active.
# To deactivate a plugin, comment out the name and any variables.
#
# Use 'telegraf -config telegraf.conf -test' to see what metrics a config
# file would generate.
#
# Environment variables can be used anywhere in this config file, simply prepend
# them with $. For strings the variable must be within quotes (ie, "$STR_VAR"),
# for numbers and booleans they should be plain (ie, $INT_VAR, $BOOL_VAR)


# Global tags can be specified here in key="value" format.
[global_tags]
  # dc = "us-east-1" # will tag all metrics with dc=us-east-1
  # rack = "1a"
  ## Environment variables can be used as tags, and throughout the config file
  # user = "$USER"
  {{ range $key, $value := environment "TAG_" }}{{ $key }}="{{ $value }}"
  {{ end -}}


# Configuration for telegraf agent
[agent]
  ## Default data collection interval for all inputs
  interval = "{{ .INTERVAL | default "10s" }}"
  ## Rounds collection interval to 'interval'
  ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
  round_interval = {{ .ROUND_INTERVAL | default "true"  }}

  ## Telegraf will send metrics to outputs in batches of at most
  ## metric_batch_size metrics.
  ## This controls the size of writes that Telegraf sends to output plugins.
  metric_batch_size = {{ .METRIC_BATCH_SIZE | default "1000" }}

  ## For failed writes, telegraf will cache metric_buffer_limit metrics for each
  ## output, and will flush this buffer on a successful write. Oldest metrics
  ## are dropped first when this buffer fills.
  ## This buffer only fills when writes fail to output plugin(s).
  metric_buffer_limit = {{ .METRIC_BUFFER_LIMIT | default "10000" }}

  ## Flush the buffer whenever full, regardless of flush_interval.
  flush_buffer_when_full = true

  ## Collection jitter is used to jitter the collection by a random amount.
  ## Each plugin will sleep for a random time within jitter before collecting.
  ## This can be used to avoid many plugins querying things like sysfs at the
  ## same time, which can have a measurable effect on the system.
  collection_jitter = "{{ .COLLECTION_JITTER | default "1s" }}"

  ## Default flushing interval for all outputs. You shouldn't set this below
  ## interval. Maximum flush_interval will be flush_interval + flush_jitter
  flush_interval = "{{ .FLUSH_INTERVAL | default "10s" }}"
  ## Jitter the flush interval by a random amount. This is primarily to avoid
  ## large write spikes for users running a large number of telegraf instances.
  ## ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
  flush_jitter = "{{ .FLUSH_JITTER | default "3s" }}"

  ## By default, precision will be set to the same timestamp order as the
  ## collection interval, with the maximum being 1s.
  ## Precision will NOT be used for service inputs, such as logparser and statsd.
  ## Valid values are "ns", "us", "ms", "s".
  precision = ""

  ## Logging configuration:
  ## Run telegraf with debug log messages.
  debug = {{ .DEBUG_MODE | default "false" }}
  ## Run telegraf in quiet mode (error log messages only).
  quiet = {{ .QUIET_MODE | default "false" }}
  ## Specify the log file name. The empty string means to log to stderr.
  logfile = ""

  ## Override default hostname, if empty use os.Hostname()
  hostname = "{{ .HOSTNAME }}"
  ## If set to true, do no set the "host" tag in the telegraf agent.
  omit_hostname = false


###############################################################################
#                            OUTPUT PLUGINS                                   #
###############################################################################

# Configuration for influxdb server to send metrics to
{{ if eq .OUTPUT_INFLUXDB_ENABLED "true" }}
[[outputs.influxdb]]
  ## The full HTTP or UDP endpoint URL for your InfluxDB instance.
  ## Multiple urls can be specified as part of the same cluster,
  ## this means that only ONE of the urls will be written to each interval.
  # urls = ["udp://localhost:8089"] # UDP endpoint example
  urls = ["{{ .INFLUXDB_URL }}"] # required
  ## The target database for metrics (telegraf will create it if not exists).
  database = "telegraf" # required
  ## Retention policy to write to.
  retention_policy = "{{ .INFLUXDB_RETENTION_POLICY | default "default" }}"
  ## Precision of writes, valid values are "ns", "us", "ms", "s", "m", "h".
  ## note: using "s" precision greatly improves InfluxDB compression.
  precision = "s"
  ## Write consistency (clusters only), can be: "any", "one", "quorum", "all"
  write_consistency = "any"

  ## Write timeout (for the InfluxDB client), formatted as a string.
  ## If not provided, will default to 5s. 0s means no timeout (not recommended).
  timeout = "{{ .INFLUXDB_TIMEOUT | default "5" }}s"
  {{ if .INFLUXDB_USER }}
  username = "{{ .INFLUXDB_USER }}"
  password = "{{ .INFLUXDB_PASS | default "metrics" }}"
  {{ end }}
  ## Set the user agent for HTTP POSTs (can be useful for log differentiation)
  # user_agent = "telegraf"
  ## Set UDP payload size, defaults to InfluxDB UDP Client default (512 bytes)
  # udp_payload = 512
{{ else }}
# InfluxDB output is disabled
{{ end }}

# Configuration for AWS CloudWatch output.
{{ if eq .OUTPUT_CLOUDWATCH_ENABLED "true" }}
[[outputs.cloudwatch]]
  ## Amazon REGION
  region = "{{ .CLOUDWATCH_REGION | default "us-east-1" }}"

  ## Namespace for the CloudWatch MetricDatums
  namespace = "{{ .CLOUDWATCH_NAMESPACE | default "InfluxData/Telegraf" }}"
{{ else }}
# Cloudwatch output is disabled
{{ end }}

# # Configuration for the Kafka server to send metrics to
{{ if eq .OUTPUT_KAFKA_ENABLED "true" }}
[[outputs.kafka]]
  ## URLs of kafka brokers
  brokers = ["{{ .OUTPUT_KAFKA_BROKER_URL | default "localhost:9092" }}"]
  ## Kafka topic for producer messages
  topic = "{{ .OUTPUT_KAFKA_TOPIC | default "telegraf" }}"
  ## Telegraf tag to use as a routing key
  ##  ie, if this tag exists, it's value will be used as the routing key
  routing_tag = "host"

  ## CompressionCodec represents the various compression codecs recognized by
  ## Kafka in messages.
  ##  0 : No compression
  ##  1 : Gzip compression
  ##  2 : Snappy compression
  compression_codec = 0

  ##  RequiredAcks is used in Produce Requests to tell the broker how many
  ##  replica acknowledgements it must see before responding
  ##   0 : the producer never waits for an acknowledgement from the broker.
  ##       This option provides the lowest latency but the weakest durability
  ##       guarantees (some data will be lost when a server fails).
  ##   1 : the producer gets an acknowledgement after the leader replica has
  ##       received the data. This option provides better durability as the
  ##       client waits until the server acknowledges the request as successful
  ##       (only messages that were written to the now-dead leader but not yet
  ##       replicated will be lost).
  ##   -1: the producer gets an acknowledgement after all in-sync replicas have
  ##       received the data. This option provides the best durability, we
  ##       guarantee that no messages will be lost as long as at least one in
  ##       sync replica remains.
  required_acks = -1

  ##  The total number of times to retry sending a message
  max_retry = {{ .OUTPUT_KAFKA_RETRIES | default "3" }}

  ## Optional SSL Config
  # ssl_ca = "/etc/telegraf/ca.pem"
  # ssl_cert = "/etc/telegraf/cert.pem"
  # ssl_key = "/etc/telegraf/key.pem"
  ## Use SSL but skip chain and host verification
  # insecure_skip_verify = false

  ## Data format to output.
  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
  data_format = "{{ .KAFKA_DATA_FORMAT | default "influx" }}"
{{ else }}
# Kafka output is disabled
{{ end }}

{{ if eq .OUTPUT_NATS_ENABLED  "true" }}
[[outputs.nats]]
## URLs of NATS servers
  servers = ["{{ .OUTPUT_NATS_URL | default "nats://localhost:4222" }}"]
  ## Optional credentials
   # username = ""
   # password = ""
   ## NATS subject for producer messages
   subject = "{{ .OUTPUT_NATS_SUBJECT | default "telegraf" }}"
   ## Optional TLS Config
   ## CA certificate used to self-sign NATS server(s) TLS certificate(s)
   # tls_ca = "/etc/telegraf/ca.pem"
   ## Use TLS but skip chain and host verification
   # insecure_skip_verify = false
   ## Data format to output.
   ## Each data format has it's own unique set of configuration options, read
   ## more about them here:
   ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
   data_format = "influx"
{{ else }}
 # Nats output is disabled
{{ end }}
# # Configuration for the file output
{{ if eq .OUTPUT_FILE_ENABLED "true" }}
# # Send telegraf metrics to file(s)
[[outputs.file]]
#   ## Files to write to, "stdout" is a specially handled file.
#   files = ["stdout", "/tmp/metrics.out"]
   files = ["{{ .OUTPUT_FILE_PATH | default "stdout" }}"]
#
#   ## Data format to output.
#   ## Each data format has it's own unique set of configuration options, read
#   ## more about them here:
#   ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_OUTPUT.md
   data_format = "influx"
#   data_format = "json"
{{ else }}
# File output is disabled
{{ end }}
###############################################################################
#                            INPUT PLUGINS                                    #
###############################################################################

# Read metrics about cpu usage
{{ if eq .INPUT_CPU_ENABLED "true" }}
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics.
  collect_cpu_time = false
  ## Comment this line if you want the raw CPU time metrics
  fielddrop = ["time_*"]
{{ else }}
  # CPU input is disabled
{{ end }}

# Read metrics about disk usage by mount point
{{ if eq .INPUT_DISK_ENABLED "true" }}
[[inputs.disk]]
  ## By default, telegraf gather stats for all mountpoints.
  ## Setting mountpoints will restrict the stats to the specified mountpoints.
  # mount_points = ["/"]

  ## Ignore some mountpoints by filesystem type. For example (dev)tmpfs (usually
  ## present on /run, /var/run, /dev/shm or /dev).
  ignore_fs = ["tmpfs", "devtmpfs"]
{{ else }}
  # Disk input is disabled
{{ end }}

# Read metrics about disk IO by device
{{ if eq .INPUT_DISKIO_ENABLED "true" }}
[[inputs.diskio]]
  ## By default, telegraf will gather stats for all devices including
  ## disk partitions.
  ## Setting devices will restrict the stats to the specified devices.
  # devices = ["sda", "sdb"]
  ## Uncomment the following line if you need disk serial numbers.
  # skip_serial_number = false
{{ else }}
  # Disk IO input is disabled
{{ end }}

# Get kernel statistics from /proc/stat
{{ if eq .INPUT_KERNEL_ENABLED "true" }}
[[inputs.kernel]]
  # no configuration
{{ else }}
  # Kernel input is disabled
{{ end }}

# Read metrics about memory usage
{{ if eq .INPUT_MEM_ENABLED "true" }}
[[inputs.mem]]
  # no configuration
{{ else }}
  # Memory input is disabled
{{ end }}

# Get the number of processes and group them by status
{{ if eq .INPUT_PROCESS_ENABLED "true" }}
[[inputs.processes]]
  # no configuration
{{ else }}
  # Process input is disabled
{{ end }}

# Read metrics about swap memory usage
{{ if eq .INPUT_SWAP_ENABLED "true" }}
[[inputs.swap]]
  # no configuration
{{ else }}
  # Swap input is disabled
{{ end }}

# Read metrics about system load and uptime
{{ if eq .INPUT_SYSTEM_ENABLED "true" }}
[[inputs.system]]
  # no configuration
{{ else }}
  # System input is disabled
{{ end }}

# Read metrics about docker containers
{{ if eq .INPUT_DOCKER_ENABLED "true" }}
[[inputs.docker]]
  ## Docker Endpoint
  ##   To use TCP, set endpoint = "tcp://[ip]:[port]"
  ##   To use environment variables (ie, docker-machine), set endpoint = "ENV"
  endpoint = "unix:///var/run/docker.sock"
  ## Only collect metrics for these containers, collect all if empty
  container_names = []
{{ else }}
  # Docker input is disabled
{{ end }}

# # Read metrics from one or more commands that can output to stdout
# [[inputs.exec]]
#   ## Commands array
#   commands = [
#     "/tmp/test.sh",
#     "/usr/bin/mycollector --foo=bar",
#     "/tmp/collect_*.sh"
#   ]
#
#   ## Timeout for each command to complete.
#   timeout = "5s"
#
#   ## measurement name suffix (for separating different commands)
#   name_suffix = "_mycollector"
#
#   ## Data format to consume.
#   ## Each data format has it's own unique set of configuration options, read
#   ## more about them here:
#   ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
#   data_format = "influx"


# # Read stats about given file(s)
# [[inputs.filestat]]
#   ## Files to gather stats about.
#   ## These accept standard unix glob matching rules, but with the addition of
#   ## ** as a "super asterisk". ie:
#   ##   "/var/log/**.log"  -> recursively find all .log files in /var/log
#   ##   "/var/log/*/*.log" -> find all .log files with a parent dir in /var/log
#   ##   "/var/log/apache.log" -> just tail the apache log file
#   ##
#   ## See https://github.com/gobwas/glob for more examples
#   ##
#   files = ["/var/log/**.log"]
#   ## If true, read the entire file and calculate an md5 checksum.
#   md5 = false


# # Read flattened metrics from one or more GrayLog HTTP endpoints
# [[inputs.graylog]]
#   ## API endpoint, currently supported API:
#   ##
#   ##   - multiple  (Ex http://<host>:12900/system/metrics/multiple)
#   ##   - namespace (Ex http://<host>:12900/system/metrics/namespace/{namespace})
#   ##
#   ## For namespace endpoint, the metrics array will be ignored for that call.
#   ## Endpoint can contain namespace and multiple type calls.
#   ##
#   ## Please check http://[graylog-server-ip]:12900/api-browser for full list
#   ## of endpoints
#   servers = [
#     "http://[graylog-server-ip]:12900/system/metrics/multiple",
#   ]
#
#   ## Metrics list
#   ## List of metrics can be found on Graylog webservice documentation.
#   ## Or by hitting the the web service api at:
#   ##   http://[graylog-host]:12900/system/metrics
#   metrics = [
#     "jvm.cl.loaded",
#     "jvm.memory.pools.Metaspace.committed"
#   ]
#
#   ## Username and password
#   username = ""
#   password = ""
#
#   ## Optional SSL Config
#   # ssl_ca = "/etc/telegraf/ca.pem"
#   # ssl_cert = "/etc/telegraf/cert.pem"
#   # ssl_key = "/etc/telegraf/key.pem"
#   ## Use SSL but skip chain and host verification
#   # insecure_skip_verify = false

{{ if eq .INPUT_NET_ENABLED "true" }}
[[inputs.net]]
  # no configuration
{{ else }}
  # Net input is disabled
{{ end }}

# # Read TCP metrics such as established, time wait and sockets counts.
{{ if eq .INPUT_NETSTAT_ENABLED "true" }}
[[inputs.netstat]]
  # no configuration
{{ else }}
  # Netstat input is disabled
{{ end }}

# Read metrics from Kafka topic(s)
{{ if eq .INPUT_KAFKA_ENABLED "true" }}
[[inputs.kafka_consumer]]
  ## topic(s) to consume
  topics = [ "{{ .INPUT_KAFKA_TOPIC | default "telegraf" }}" ]
  ## an array of Zookeeper connection strings
  zookeeper_peers = ["{{ .INPUT_KAFKA_ZOOKEEPER_PEER | default "zookeeper:2181" }}"]
  zookeeper_chroot = "{{ .INPUT_KAFKA_ZOOKEEPER_CHROOT | default "" }}"
  ## the name of the consumer group
  consumer_group = "telegraf_metrics_consumers"
  ## Maximum number of metrics to buffer between collection intervals
  metric_buffer = 100000
  ## Offset (must be either "oldest" or "newest")
  offset = "oldest"

  ## Data format to consume.

  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "{{ .KAFKA_DATA_FORMAT | default "influx" }}"
{{ else }}
  # Kafka input is disabled
{{ end }}

# Generic TCP listener
{{ if eq .INPUT_LISTENER_ENABLED "true" }}
[[inputs.tcp_listener]]
  ## Address and port to host TCP listener on
  service_address = ":{{ .INPUT_LISTENER_PORT | default "8094" }}"

  ## Number of TCP messages allowed to queue up. Once filled, the
  ## TCP listener will start dropping packets.
  allowed_pending_messages = 10000

  ## Maximum number of concurrent TCP connections to allow
  max_tcp_connections = 250

  ## Data format to consume.
  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "{{ .INPUT_LISTENER_DATA_FORMAT | default "json" }}"
{{ else }}
  # TCP listener input is disabled
{{ end }}

# Read metrics from NATS subject(s)
{{ if eq .INPUT_NATS_ENABLED "true" }}
[[inputs.nats_consumer]]
  ## urls of NATS servers
  servers = ["{{ .INPUT_NATS_URL | default "nats://localhost:4222" }}"]
  ## Use Transport Layer Security
  secure = false
  ## subject(s) to consume
  subjects = ["{{ .INPUT_NATS_SUBJECT | default "telegraf" }}"]
  ## name a queue group
  queue_group = "telegraf_consumers"
  ## Maximum number of metrics to buffer between collection intervals
  metric_buffer = 100000
  ## Sets the limits for pending msgs and bytes for each subscription
  ## These shouldn't need to be adjusted except in very high throughput scenarios
  # pending_message_limit = 65536
  # pending_bytes_limit = 67108864

  ## Data format to consume.

  ## Each data format has it's own unique set of configuration options, read
  ## more about them here:
  ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
  data_format = "influx"
{{ else }}
  # NATS consumer input is disabled
{{ end }}

#
#   ## Data format to consume.
#   ## Each data format has it's own unique set of configuration options, read
#   ## more about them here:
#   ## https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md
#   data_format = "influx"

# # Read metrics of haproxy, via socket or csv stats page
{{ if eq .INPUT_HAPROXY_ENABLED "true" }}
[[inputs.haproxy]]
#   ## An array of address to gather stats about. Specify an ip on hostname
#   ## with optional port. ie localhost, 10.10.3.33:1936, etc.
#   ## Make sure you specify the complete path to the stats endpoint
#   ## including the protocol, ie http://10.10.3.33:1936/haproxy?stats
#   #
#   ## If no servers are specified, then default to 127.0.0.1:1936/haproxy?stats
#   servers = ["http://myhaproxy.com:1936/haproxy?stats"]
#   ##
#   ## You can also use local socket with standard wildcard globbing.
#   ## Server address not starting with 'http' will be treated as a possible
#   ## socket, so both examples below are valid.
#   ## servers = ["socket:/run/haproxy/admin.sock", "/run/haproxy/*.sock"]
  {{ if .INPUT_HAPROXY_SERVER }}
  servers = ["{{ .INPUT_HAPROXY_SERVER }}"]
  {{ end }}
{{ else }}
  # haproxy input is disabled
{{ end }}
