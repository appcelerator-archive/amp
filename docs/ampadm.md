# AMP - ampadm command line documentation


# help on line

to get help on line used: `ampadm --help`
or for a specific command: `ampadm [command] --help`



# ampadm main commands

- version
- cluster
- platform
- node
- bootstrap


# ampadm common options

- --verbose (-v): display more information during the current command execution
- --silence (to be replaced by --quiet): display no message at all, the process return shell compatible 0 if ok, 1 if error
- --config: define which config file to use, default is: $HOME/.config/amp/amp.yaml
- --server: define the admin server address (override config): format [swarmManager address]:[port], default: 127.0.0.1:30313


## apadm version

display the ampadm version and build, ad for instance:

```
$ ampadm version
ampadm version: v0.5.1-dev, build: 8b1d0b87)
```

## admadm cluster

`cluster` is a set of commands to start and stop the amp admin services: 
- adm-server: locate in the swarm manager machine, this service get command from client, send orders to all adm-agents and aggregate their answers to send back answer to client.
- adm-agent: locate on each swarm-machine, this service get adm-server command to execute them on the node context.

these commands are usable only on the same machine where admadm is installed.

### ampadm cluster start

Create and start service adm-server and adm-agent on the local machine.

options:
- --force (-f): do not remove all services if one service failed to start
- --local (-l): use image appcelerator/amp:local to start adm services
- --tag [name]: use image appcelerator/amp:[name] to start adm service (so --tag local is equivalent to --local)

### ampadm cluster stop

remove local services adm-server and adm-agent


## ampadm platform (pf)

`platform` is a set of command to manage the amp infrastructure services. Considering the adm-server location (set in the configuration file or using option --server), these commands can be executed remotely, so from a single machine it's possible to handle several amp clusters.


### ampadm pf pull

ask all node in the targeted amp cluster to pull all needed infrastructure docker images


### ampadm pf start

Create and start all amp infrastructure services on the targeted amp cluster
Available options are:
- --force (-f): do not remove all services if one service failed to start
- --local (-l): use image appcelerator/amp:local to start infrastructure services included in this image
- --tag [name]: use image appcelerator/amp:[name] to start infrastructure services included in this image (so --tag local is equivalent to --local)

### ampadm pf stop

Remove all created services included the users ones on the targeted amp cluster

### ampadm pf monitor

display information about all the started service (infrastructure and user) running on the targeted amp cluster
Information is: 
- docker service id, 
- service name, 
- service global status: running, partially running (if some services are missing or failing), stopped
- mode: replicated or global (one task on each cluster node) 
- number of running replicas
- number of task failed: number of container crashed and restarted.


This information is refreshed every second until user break it by a ctrl-C

### ampadm pf status

display the global status of the targeted amp cluster:
- running: all infrastructure service are up and running
- partially running: some infrastructure service are missing or failing
- stopped: all infrastructure services are removed

## ampadm node

`node` is a set of command to manage the amp cluster nodes. Considering the adm-server location (set in the configuration file or using option --server), these commands can be executed remotely, so from a single machine it's possible to handle several amp clusters.

### ampadm node ls

display information of each targeted swarm cluster nodes 
Information is:
- docker node id, 
- role: "manager" (swarm manager), or "node" (swarm node)
- hostname: hostname of the virtual machine which host the swarm agent.
- address: adm-agent address
- version: docker version
- status: node status
- cpu: number of available cpu
- mem: size of the available ram

Available options:
- follow (-f): the information is updated every second and stop only at user Ctrl-C
- more: in addition to the current information display: Agent containerId, Host OS, Architecture 32 or 64 bits


### ampadm node count

display aggregated information for each targeted swarm cluster nodes and a global total:
informations is:
- Total number of containers
- Number of running containers
- Number of paused containers
- Number of stopped containers
- number of images

### ampadm node purge

purge containers and/or images on one or all targeted swarm cluster nodes
available options are:
- --node [id], if specified purge only the node having id [id], otherwise purge all cluster nodes
- --container, if set purge all stopped containers
- --images, if set purge all images not running in a container
- --force (-f): if set force to remove images and/or containers even linked to active resources. 


### ampadm bootstrap (bs)

`bootstrap` is a set of commands to install swarm cluster and prepare amp startup
it should start on machine having docker installed and other prerequisite detailed in md files in `docs/installation`
bootstrap command verify prerequisite and install swarm manager or node considering the options:
available options are:
- --manager (-m): to bootstrap a manager node, without a swarm node is bootstrapped 
- --create (-c) : to initialize the swarm, only possible with the --manager (-m) option
- --host [host] set the swarm manager hostname to [host]
- --port [port] set the swarm manager port to [port] (default 2377)
- --tag [tag] indicate to use amp image with tag [tag], default: latest
- --token [token] set the token used by a manager or a node to join the swarm cluster

