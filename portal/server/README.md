###amp-agent

amp-agent is an infrastructure service installed on each node sending local node Docker events and local containers logs to kafka.

### version 1.0.0

Get Docker events related to container create, start, stop , kill, die, destroy, use it to maintain internally a list of running containers with their id, node_id, service_id, service_name and send the docker events to kafka.
For each running container, open a log stream and send them to Kafka.

### api v1.0.0

api doc:


    * /api/v1/health: return code 200 if amp-agent is ready, maining log-agent has open an event stream with docker and this stream is still active

    * /api/v1/containers: return a json containing the active containers list with service_id, service_name, status, and health status
    