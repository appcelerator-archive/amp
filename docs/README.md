AMP: A... Microservices Platform ![WIP](https://cdn.rawgit.com/appcelerator/amp/master/docs/static_files/amp--docs-WIP-yellow.svg)
============================

AMP is currently under development and this section is here to help you get started based on latest stable tagged version of the project. If you're here, then you **ARE** pioneering with us and we encourage you to [get in touch](./contributing.md) !

# Installation

AMP can be installed either on your computer for building solutions or on a cluster of servers for running them. To get started, [check out the installation instructions in the documentation](./installation/index.md)

# Configuration

## Docker Swarm cluster

AMP requires a Docker Swarm cluster to operate.
If you use an automated way to deploy AMP in a cluster, this step is already done, you can skip this paragraf.
If you installed AMP locally, here are the instructions to create a local Docker Swarm cluster:

    $ docker swarm init
    
It should provide the following output:

    Swarm initialized: current node (ej9yivb39rrq2iyk3itdqvcq1) is now a manager.
      
    To add a worker to this swarm, run the following command:
      
        docker swarm join \
        --token SWMTKN-1-08xe2j6h2y812exq4rw5cj7j98112gn2ar88s9kkniimmn4i74-1bkwl472uc7llf4divn7k3bkv \
        10.128.27.12:2377
      
    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
 
Make sure you Docker Swarm cluster is up and ready by typing the following command:

    $ docker node ls
    
## AMP CLI

The AMP Command Line Interface is a gRPC client connecting to the AMP daemon on `127.0.0.1:8080` by default, which is fine if you're using AMP locally.
If you have deployed your AMP Swarm cluster somewhere else and you want to control it remotely, please make sure to specify the valid AMP daemon address either:

 * as a command line parameter:

        $ amp --server amp.some.where:8080 [command]

 * In the configuration file:

        $ cat ~/.amp.yaml
        ServerAddress: amp.some.where:8080
        $ amp [command]
  
# Starting AMP

Use the AMP command `amp platform` (or the `amp pf` short version) to manipulate AMP swarm services.

On first usage, you need to pull AMP images to your local docker installation using:
    
    $ amp pf pull
    
After this step, you can start AMP with the following command:
                                                                                    
    $ amp pf start

Finally, you can monitor AMP status by using:

    $ amp pf monitor

## AMP platform commands

AMP platform commands are useful to manage your AMP swarm deployment:

 * pull
 * start
 * stop
 * status
 * monitor

The usual workflow looks like this:

    $ amp pf pull
    $ amp pf start
    $ amp pf monitor (better in a separate console)
    $ amp pf stop

Options:

    -v --verbose    To have more information messages
    -s --silence    To not have message at all
    -f --force      Only for 'start', to force amp restart if amp is already started or do not exit on error if a service doesn't start


## CLI

`amp --help` displays helps for available AMP commands.

### Running a service

    amp service create
    amp service rm

### Running a stack

    amp stack ls
    amp stack start
    amp stack rm
    amp stack stop
    amp stack up

### Logs

The `amp logs` command is used to query or stream logs. It provides useful filtering options to manage what is presented.

    $ amp logs --help

    Usage:  amp logs [OPTIONS] [SERVICE]

    Fetch log entries matching provided criteria. If provided, SERVICE can be a partial or full service id or service name.

    Options:
          --config string      Config file (default is $HOME/.amp.yaml)
          --container string   Filter by the given container
      -f, --follow             Follow log output
      -h, --help               help for logs
          --message string     Filter the message content by the given pattern
      -m, --meta               Display entry metadata
          --node string        Filter by the given node
      -n, --number string      Number of results (default "100")
          --server string      Server address
          --stack string       Filter by the given stack
      -v, --verbose            Verbose output


A few useful examples:

* To fetch and follow all the logs from the whole AMP platform:
```
  $ amp logs -f
```

* To fetch and follow the logs for a specific service, with the message content only:
```
  $ amp logs -f etcd
```

* To search for a specific pattern through all the logs of the platform:
```
  $ amp logs --message error
```

* To fetch and follow the logs for a `elasticsearch`, using partial service name:
```
  $ amp logs -f ela
```

### Stats

The `amp stats` command provides useful information about resource consumption. There is a comprensive set of options
to query and monitor specfic metrics that complements and extends what is visible in the web dashboard (http://localhost:6001).

    $ amp stats --help

    Get statistics on containers, services, nodes about cpu, memory, io, net.

    Usage:
      amp stats [flags]

    Flags:
          --container               display stats on containers
          --container-id string     filter on container id
          --container-name string   filter on container name
          --cpu                     display cpu stats
          --datacenter string       filter on datacenter
      -f, --follow                  Follow stat output
          --host string             filter on host
          --image string            filter on container image name
          --io                      display disk io stats
          --mem                     display memory stats
          --net                     display net rx/tx stats
          --node                    display stats on nodes
          --node-id string          filter on node id
          --period string           historic period of metrics extraction, duration + time-group as 1m, 10m, 4h, see time-group
          --service                 displat stats on services
          --service-id string       filter on service id
          --service-name string     filter on service name
          --since string            date defining when begin the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm
          --task                    display stats on tasks
          --task-id string          filter on task id
          --task-name string        filter on task name
          --time-group string       historic extraction group can be: s:seconds, m:minutes, h:hours, d:days, w:weeks
          --until string            date defining when stop the historic metrics extraction, format: YYYY-MM-DD HH:MM:SS.mmm

    Global Flags:
          --Config string   Config file (default is $HOME/.amp.yaml)
          --target string   target environment ("local"|"virtualbox"|"aws") (default "local")
      -v, --verbose         verbose output

A few useful examples:

* To display list of services with cpu, mem, io, net metrics and follow them
```
  $ amp stats --service -f
```

* To display last 10 minutes of historic of the containers of service nats with cpu, mem, io, net metrics and follow them:
```
  $ amp stats --container --service-name=nats --period=10m  -f
```

* To display list of tasks with only cpu and mem metrics
```
  $ amp stats --task --cpu --mem
```

### Managing Docker images

    amp registry ls
    amp registry push

More details in the [userguide](userguide/registry.md)
