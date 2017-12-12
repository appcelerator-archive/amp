# Deploying a stack

## What is a stackfile?
A stackfile is a Docker compose v3 [file](https://docs.docker.com/compose/compose-file/) which defines one or more services to be deployed. 
  
## Stack deploy
Here is how you can deploy your own stack file:
  
1. Create a simple YAML file and save it in your working directory as `pinger.yml`
```
version: "3"
  
networks:
  public:
    external: true
  
services:
  
  pinger:
    image: appcelerator/pinger:${TAG:-latest}
    networks:
      - public
    environment:
      SERVICE_PORTS: "3000"
      VIRTUAL_HOST: "pinger.examples.*,https://pinger.examples.*"
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
      placement:
        constraints: [node.labels.amp.type.user == true]
```
 
2. Simply run the following command:
```
$ amp stack deploy -c pinger.yml
[user user1 @ 127.0.0.1:50101]
Deploying stack pinger using pinger.yml
Creating service pinger_pinger
```
  
You can also deploy one of our many example stackfiles in the `examples/stacks` directory:
  
1. Clone the `amp` repository from GitHub.
```
$ git clone https://github.com/appcelerator/amp.git
```
  
2. Go into the `examples/stacks` directory with various applications designed for demos and quick starts:
```
$ cd amp/examples/stacks
```
  
For help with deploying the stacks, simply follow the instructions in the README for each of the examples in their respective directories.
 
## Dashboard

To help monitoring user deployed stacks, we provide access to Grafana Dashboard and Kibana Dashboard.

### Grafana

local: http://dashboard.local.appcelerator.io

hosted: https://dashboard.YOUR.DOMAIN

### Kibana

local: http://kibana.local.appcelerator.io

hosted: https://kibana.YOUR.DOMAIN
  
## Viewing and filtering logs

The `amp logs` command allows for the querying and filtering of both real-time and historical logs.

To fetch the logs for your cluster:
```
$ amp logs
[user user1 @ 127.0.0.1:50101]
         pinger_pinger.2 | listening on :3000
         pinger_pinger.1 | listening on :3000
         pinger_pinger.3 | listening on :3000
```

Following on from the previous example of deploying a stack, if you want to fetch the logs of the `pinger` stack,
you can run:
```
$ amp stack logs pinger
[user user1 @ 127.0.0.1:50101]
         pinger_pinger.2 | listening on :3000
         pinger_pinger.1 | listening on :3000
         pinger_pinger.3 | listening on :3000
```
This will get the logs of every service in the stack.

If you want to get the logs for the individual services in the stack, you can run:
```
$ amp service logs pinger
[user user1 @ 127.0.0.1:50101]
         pinger_pinger.2 | listening on :3000
         pinger_pinger.1 | listening on :3000
         pinger_pinger.3 | listening on :3000
```

> NOTE: The above example is just a sample of the output of `amp logs` command and its sub-commands. The output may vary on different machines.

For more detailed examples of the querying and filtering options for logs, see the [logs documentation](reference/logs.md)  
  