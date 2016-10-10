AMP: A... Microservices Platform ![WIP](./static_files/amp--docs-WIP-yellow.svg)
============================

AMP is currently under development and this section is here to help you get started based on latest stable tagged version of the project. If you're here, then you **ARE** pioneering with us and we encourage you to [get in touch](#Contribute) !


Prerequisites
===============

All of the infrastructure services and almost all tooling is completely
containerized:

* [Docker](https://www.docker.com/products/docker) - the container engine
* [Go](https://golang.org/) - Binary distribution for Go programming language
* [Glide](https://glide.sh/) - Package Management for Go
* `make` (almost all the Makefile rules execute in containers now)

AMP - Getting Started
==============

AMP can be installed either on your computer for building solutions or on servers for running them. To get started, [check out the installation instructions in the documentation](./installation/index.md)

AMP Usage
==============

AMP comes with a complete toolbox to help you manage & operate your solutions. Find below AMP client sample usage:

Swarm
--------------

Use the `swarm` shell script to launch amp cluster services. Available commands are:

 * pull
 * start
 * ls
 * restart
 * stop
 * monitor

The usual workflow looks like this:

    $ ./swarm pull
    $ sudo ./swarm start --min
    $ ./swarm monitor
    $ sudo ./swarm restart --min (equivalent to stop, pull, start)
    $ sudo ./swarm stop

AMP Client
--------------

AMP currently uses GitHub authentication for logging in. The CLI command `amp login` will use your GitHub credentials to authenticate
you via GitHub. You need to be a member of the AMP team in the Appcelerator organization to be authorized. For the GitHub OAuth flow
to work, you need to start `amplifier` with a client ID and secret before you can run the `amp login` command. Follow these instructions
for now:

* [Create](https://github.com/settings/applications/new) an OAuth application (to simulate the amp backend)
* Note the `Client ID` and `Client Secret`
* Start `amplifier` with the following flags:
    $ amplifier --clientid <Client ID> --clientsecret <Client Secret> &
* Use the CLI to login
    $ amp login
* When prompted, enter username, password, and two-factor auth code

### Logs

The `amp logs` command is used to query or stream logs. It provides useful filtering options to manage what is presented.

    $ amp logs --help

    Search through all the logs of the system and fetch entries matching provided criteria.

    Usage:
      amp logs [flags]

    Flags:
          --container string      Filter by the given container id
      -f, --follow                Follow log output
          --from string           Fetch from the given index (default "-1")
          --message string        Filter the message content by the given pattern
      -m, --meta                  Display entry metadata
          --node string           Filter by the given node id
      -n, --number string         Number of results (default "100")
          --service-id string     Filter by the given service id
          --service-name string   Filter by the given service name

    Global Flags:
          --Config string   Config file (default is $HOME/.amp.yaml)
          --server string   Server address (default "localhost:50101")
          --target string   target environment ("local"|"virtualbox"|"aws") (default "local")
      -v, --verbose         verbose output


A few useful examples:

* To fetch and follow all the logs from the platform:
```
  $ amp logs -f
```

* To fetch and follow the logs for a specific service, with the message content only:
```
  $ amp logs -f --service_name etcd
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

* To display last 10 minutes of historic of the containers of service kafka with cpu, mem, io, net metrics and follow them:
```
  $ amp stats --container --service-name=kafka --period=10m  -f
```

* To display list of tasks with only cpu and mem metrics
```
  $ amp stats --task --cpu --mem
```


<a name="Contribute"></a>Contributing to AMP ![WIP](static_files/amp--docs-WIP-yellow.svg)
======================

Want to hack on AMP? Awesome! We have [instructions to help you get
started contributing code or documentation is on the way](misc/who-written-for.md).

These instructions are probably not perfect, please let us know if anything
feels wrong or incomplete. Better yet, submit a PR and improve them yourself.


### Talking to other AMP users and contributors

<table class="tg">
  <col width="45%">
  <col width="65%">
  <tr>
    <td>Gitter&nbsp;Chat&nbsp;</td>
    <td>
      <p>
        <a href="https://gitter.im/appcelerator/amp" target="_blank">our chat room</a>.
      </p>
    </td>
  </tr>
  <tr>
    <td>AMP Slack Channel</td>
    <td>
      Come join us on the #amp-community channel.
    </td>
  </tr>
</table>
