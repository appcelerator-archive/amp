# AMP - the open source CaaS & microsevice platform for Docker

AMP is an open-source Container-as-a-Service ([CaaS](https://blog.docker.com/2016/02/containers-as-a-service-caas/))
platform for managing and monitoring containerized applications,  microservices, and microservice workers as part
of a unified [serverless computing](https://en.wikipedia.org/wiki/Serverless_computing) environment.
It is based on the swarm mode features introduced with Docker 1.12, providing developers with
a straightforward development strategy for building on these features without the need to stray from the
core Docker ecosystem and adopt complex third party frameworks.

There is a 5 minute introductory video under the `docs` directory.

This current version of this project is only a few months old and under heavy development. It is not
ready for production.

## Join the conversation

Talk to other AMP users and contributors

## Contributing

See [contributing](project/CONTRIBUTING.md)

## Getting started

See

## Starting AMP

Use the `swarm` shell script to launch amp swarm services. Available commands are:

 * pull
 * start
 * ls
 * restart
 * stop
 * monitor

The usual workflow looks like this:

    $ ./swarm pull
    $ sudo ./swarm start
    $ ./swarm monitor
    $ sudo ./swarm restart (equivalent to stop, pull, start)
    $ sudo ./swarm stop

## CLI

`amp --help` displays helps for available AMP commands.

### Running a service

    amp service create
    amp service rm

### Running a stack

    amp stack ls
    amp stack restart
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
