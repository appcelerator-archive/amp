# docker-telegraf

Docker Image for [InfluxData Telegraf](https://influxdata.com/time-series-platform/telegraf/), based on Alpine linux.

## Run

You may need to replace the path to */var/run/docker.sock* depending on the location of your docker socket.

Most basic form:

    docker run -t -v /var/run/docker.sock:/var/run/docker.sock:ro appcelerator/telegraf

Custom InfluxDB location, additional tags, and retention policy for InfluxDB 1.0.0:

    docker run -t -v /var/run/docker.sock:/var/run/docker.sock:ro -v /var/run/utmp:/var/run/utmp:ro -e INFLUXDB_URL=http://influxdb:8086 -e TAG_datacenter=eu-central-1 -e TAG_type=core -e INFLUXDB_RETENTION_POLICY= appcelerator/telegraf

# Configuration (ENV, -e)

Variable | Description | Default value | Sample value 
-------- | ----------- | ------------- | ------------
HOSTNAME | To pass in the docker host's actual hostname | | localhost
TAG_\<name\> | Adds a tag with the given value to all measurements | | TAG_datacenter=eu-central-1
METRIC_BATCH_SIZE | size of writes that Telegraf sends to output plugins | 1000 |
METRIC_BUFFER_LIMIT | buffer limit for failed writes | 10000 |
DEBUG_MODE | Run telegraf in debug mode | false | true
QUIET_MODE | Run telegraf in quiet mode | false | true
INTERVAL | Data collection interval | 10s |
ROUND_INTERVAL | Round collection interval | true |
COLLECTION_JITTER | Collection jitter by a random amount | 1s |
FLUSH_INTERVAL | Default flushing interval for all outputs | 10s |
FLUSH_JITTER | Jitter the flush interval by a random amount | 3s |
OUTPUT_INFLUXDB_ENABLED | enable InfluxDB Output | true |
OUTPUT_CLOUDWATCH_ENABLED | enable Amazon Cloudwatch Output | false |
OUTPUT_KAFKA_ENABLED | enable Kafka Output | false |
OUTPUT_NATS_ENABLED | enable NATS Output | false |
OUTPUT_FILE_ENABLED | enable File Output | false |
INPUT_KAFKA_ENABLED | enable Kafka Input | false |
INPUT_NATS_ENABLED | enable Nats Input | false |
INPUT_CPU_ENABLED | enable cpu metrics | true |
INPUT_DISK_ENABLED | enable disk metrics | true |
INPUT_DISKIO_ENABLED | enable disk I/O metrics | true |
INPUT_KERNEL_ENABLED | enable kernel metrics | false |
INPUT_MEM_ENABLED | enable mem metrics | true |
INPUT_PROCESS_ENABLED | enable process metrics | true |
INPUT_SWAP_ENABLED | enable swap metrics | true |
INPUT_SYSTEM_ENABLED | enable system metrics | true |
INPUT_NETSTAT_ENABLED | enable netstat metrics | false |
INPUT_NET_ENABLED | enable net metrics | true |
INPUT_LISTENER_ENABLED | enable generic TCP listener | false |
INPUT_DOCKER_ENABLED | enable Docker metrics | true |
INPUT_HAPROXY_ENABLED | enable haproxy metrics | false |
INFLUXDB_URL | Where is your InfluxDB running? | http://localhost:8086 | http://influxdb:8086
INFLUXDB_RETENTION_POLICY | Set the name of the policy | default | autogen
INFLUXDB_USER | InfluxDB username | |
INFLUXDB_PASS | InfluxDB password | metrics |
INFLUXDB_TIMEOUT | InfluxDB timetout (in seconds) | 5 |
CLOUDWATCH_REGION | Amazon region | us-east-1 |
CLOUDWATCH_NAMESPACE | Namespace | InfluxData/Telegraf |
INPUT_NATS_URL | Url of NATS server | nats://localhost:4222 |
INPUT_NATS_SUBJECT | Subject to consume | telegraf |
INPUT_LISTENER_PORT | Port of the generic TCP listener | 8094 |
INPUT_LISTENER_DATA_FORMAT | Data format of the generic TCP listener | json |
INPUT_HAPROXY_SERVER | haproxy address | 127.0.0.1:1931/haproxy?stats | /var/run/haproxy/admin/sock
OUTPUT_NATS_URL | URL of NATS server | localhost:4222 |
OUTPUT_NATS_SUBJECT | NATS Subject for producer messages | telegraf |
OUTPUT_KAFKA_BROKER_URL | Kafka broker URL in output | localhost:9092 |
OUTPUT_KAFKA_TOPIC | Kafka topic on which to write | telegraf |
OUTPUT_KAFKA_RETRIES | Number of retries for the connection to Kafka | 3 |
OUTPUT_FILE_PATH | absolute path to the file, would better be mounted | stdout |
INPUT_KAFKA_BROKER_URL | Kafka broker URL in input | localhost:9092 |
INPUT_KAFKA_TOPIC | Kafka topic on which to read | telegraf |
INPUT_KAFKA_ZOOKEEPER_PEER | Zookeeper peers used by Kafka in input | zookeeper:2181 |
INPUT_KAFKA_ZOOKEEPER_CHROOT | Zookeeper chroot path | |
KAFKA_DATA_FORMAT | Kafka data format | influx |

## Tags

- telegraf-0.13
- telegraf-1.0.1, telegraf-1.0
- telegraf-1.1.2, telegraf-1.1, latest
