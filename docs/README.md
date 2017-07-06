# Getting started

* [Installation](#installation)
  * [Prerequisites](#prerequisites)
    * [MacOS](#macos)
    * [Windows](#windows)
    * [Linux](#linux)
  * [Getting the CLI](#getting-the-cli)
  * [Using the unstable version](#using-the-unstable-version)
* [Usage](#usage)
  * [Target an AMP cluster](#target-an-amp-cluster)
    * [Creating a local AMP cluster](#creating-a-local-amp-cluster)
  * [Signing up and logging in](#signing-up-and-logging-in)
  * [UI and Dashboard](#ui-and-dashboard)
    * [UI](#ui)
    * [Grafana](#grafana)
  * [Examples](#examples)
    * [Deploying a stack](#deploying-a-stack)
    * [Viewing and filtering logs](#viewing-and-filtering-logs)
    * [Displaying metrics and statistics](#displaying-metrics-and-statistics)
    * [Organization and team management](#organization-and-team-management)

## Installation

### Prerequisites

We recommend installing the Docker CE stable channel release 17.06 on your system.
Make sure you check the "What to know before you install" section on the Docker site to ensure your system meets the requirements.

>NOTE: We currently do not support Docker Toolbox on any OS.

#### MacOS

On MacOS, Docker can be installed via the Docker website.

[Docker for Mac](https://docs.docker.com/docker-for-mac/install/)

#### Windows

On Windows, Docker can be installed via the Docker website.

[Docker for Windows](https://docs.docker.com/docker-for-windows/install/)

#### Linux

On Linux, Docker can be installed via the Docker website. You can find your specific distribution in the tables.

[Docker for Linux](https://docs.docker.com/engine/installation/)

In addition, there is another step you must perform for your system.

Increase virtual memory (needed for Elasticsearch):
```
$ sudo sysctl -w vm.max_map_count=262144
```
To make this change permanent, you can run the following and reboot:
```
$ echo "vm.max_map_count = 262144" | sudo tee -a /etc/sysctl.conf
```

### Downloading the CLI

Download the latest release of the CLI for your platform:

https://github.com/appcelerator/amp/releases

Place the `amp` binary from the release archive into some location in your system path. For example:
```
$ sudo mv ./amp /usr/local/bin
```

### Certificates

The connection between the CLI and the ANP server is secured with TLS.
In the case the certificate on the server is not valid (self signed or expired), you can use the `-k` option.
For local deployment you can also add the CA to your local key store.

### Using the unstable version

Alternatively, if you wish to work with the latest development version directly from `master` on Github,
you can fork the repo and build your own CLI. The entire toolchain has been containerized so you do not need
`go` or `gnu` on your system, just docker.

To get the repo setup, you can run the following set of commands:
```
$ cd $GOPATH/src
$ mkdir -p github.com/appcelerator
$ cd github.com/appcelerator
$ git clone https://github.com/appcelerator/amp
$ cd amp
$ export PATH=bin/{YourOS}/amd64:$PATH
```
This will clone the repository into your `go` workspace and add the CLI path to your system `$PATH`.

To build the CLI, you can then run:
```
$ ampmake buildall-cli
```
This will build cross-compiled versions of the CLI and place them in the appropriate locations under `./bin`.

If you already have `gnu` installed on your system, you can instead run:
```
$ make build-cli
```
This will build your OS-specific CLI in `./bin` and it is much faster.

## Usage

Run `amp --help` to get the CLI help menu for amp.

```
$ amp --help

Usage:	amp [OPTIONS] COMMAND 

Deploy, manage, and monitor container stacks.

Examples:
amp version

Options:
  -h, --help            Print usage
  -k, --insecure        Control whether amp verifies the server's certificate chain and host name
  -s, --server string   Specify server (host:port)
  -v, --version         Print version information and quit

Management Commands:
  cluster     Cluster management operations
  org         Organization management operations
  password    Password management operations
  service     Service management operations
  stack       Stack management operations
  team        Team management operations
  user        User management operations

Commands:
  config      Display configuration
  login       Login to account
  logout      Logout of account
  logs        Display logs matching provided criteria
  stats       Display statistics
  version     Show version information
  whoami      Display currently logged in account name

Run 'amp COMMAND --help' for more information on a command.
```

>*A large number of AMP commands are interactive, so if necessary options are not provided
you will be prompted to provide them before the command runs*.

### Target an AMP cluster

When using `amp`, every command takes the `--server|-s` option.
This option is used to define the target cluster you will be running your commands with.
The target server is displayed at the top of every CLI command output.
You can also add a default `--server|-s` value to your `amp.yml` config file.

>For more information on creating your configuration file, see the [config documentation](config.md).

The server location of our hosted AMP is `COMING SOON`. This location is currently
acting as a playground for anyone to signup and use as they like.

The default value for the `--server|-s` option is `localhost:50101`, which points to a local
cluster that can be deployed on your system's docker engine.

In order to make sure you're connected to the specified server when running commands,
you can run the `amp version` command to test this.

```
$ amp -s your.server.com version
[your.server.com:50101]
Client:
 Version:       v0.12.0-dev
 Build:         fe0947b5
 Server:        your.server.com:50101
 Go version:    go1.8.1
 OS/Arch:       darwin/amd64

Server:
 Version:       v0.12.0-dev
 Build:         fd732802
 Go version:    go1.8.1
 OS/Arch:       linux/amd64
```

#### Creating a local AMP cluster

```
$ amp cluster create
```
>*If you are using the unstable version of the CLI, you will need to build the `amp-bootstrap`
image using the command `make build-bootstrap` or `ampmake build-bootstrap` if your don't have the build tools,
and then specify the `--tag=local` option to the CLI*

For more examples of cluster commands and deployment in different environments, see the [cluster documentation](cluster.md).

### Signing up and logging in

In order to be able to use AMP, you will need to signup for a user account, verify the account
and login.

>NOTE: Don't forget to specify your desired server target!

To signup for amp use:
```
$ amp user signup
username: sample
email: sample@user.com
password: [password]
```
After signing up, you will then be sent an email to your registered address. In this email, you will
be sent a link to verify your account with or you can verify your account with the provided CLI command.
>*The verification step is skipped for local deployment, you'll be logged automatically*

To verify your account using the token in verification email.
```
$ amp user verify [token]
```
>NOTE: You do not need to verify your account on a local cluster

To login to amp
```
$ amp login
username: sample
password: [password]
```

If you are having issues signing up or logging in, see the [user documentation](user.md).

If you are have forgotten or want to change your password, see the [password documentation](password.md)

### UI and Dashboard

Along with our CLI, after deploying a cluster to your desired environment, you can get access to our Custom
UI, a Grafana Dashboard and Kibana Dashboard.

#### UI

local: http://local.appcelerator.io

hosted: https://cloud.YOUR.DOMAIN

Note: TLS can be use also for a local deployment, but you'll have to add the self signed certificate to your key store.

#### Grafana

local: http://dashboard.local.appcelerator.io

hosted: https://dashboard.YOUR.DOMAIN

#### Kibana

local:  http://kibana.local.appcelerator.io

hosted:  https://kibana.YOUR.DOMAIN

### Examples

#### Deploying a stack

To deploy a standard Docker Compose version 3 stackfile into your cluster, use `amp stack deploy`.
There are a number of sample stackfiles under `examples/stacks` that you can try out.

Here is a simple example:
```
$ curl -O https://raw.githubusercontent.com/appcelerator/amp/master/examples/stacks/pinger/pinger.yml
$ amp stack deploy -c pinger.yml
$ amp stack ls
$ amp service logs pinger
$ curl http://pinger.examples.local.appcelerator.io/ping
```
Or browse to: https://pinger.examples.local.appcelerator.io/ping.

For more information on what you can do with your stack when it is deployed, see the [stacks documentation](stack.md).

For more information on inspecting and manipulating the services within your stack, see the [service documentation](service.md)

#### Viewing and filtering logs

The `amp logs` command allows for the querying and filtering of both realtime and historical logs.

To fetch the logs of your cluster:
```
$ amp logs
...
Cluster status: healthy
2017/06/29 17:59:32 listening on :3000
2017/06/29 17:59:37 listening on :3000
2017/06/29 17:59:37 listening on :3000
...
```

Following on from the previous example of deploying a stack, if you want to fetch the logs of the `pinger` stack,
you can run:
```
$ amp stack logs pinger
2017/06/29 17:59:32 listening on :3000
...
```
This will get the logs of every service in the stack.

If you want to get the logs for the individual services in the stack, you can run:
```
$ amp service logs pinger
2017/06/29 17:59:32 listening on :3000
...
```

For more detailed examples of the querying and filtering options for logs, see the [logs documentation](logs.md)

#### Displaying metrics and statistics

The `amp stats` command allows for the querying and filtering of realtime and historical statistics.

To fetch the stats of your cluster
```
$ amp stats
Stats on services period=now-10m, No filter
SERVICE         CPU %   MEM USAGE   LIMIT    MEM %   IO READ   IO WRITE   NET RX       NET TX
```

Following on from the previous example of deploying a stack, if you want to display the stats of the `pinger` stack,
you can run:
```
$ amp stats --stack
Stats on stacks period=now-10m, No filter
STACK     CPU %   MEM USAGE   LIMIT    MEM %   IO READ   IO WRITE   NET RX       NET TX
pinger    0.00%   6.7 MB      5.9 GB   0.3%    0/s       0/s        212.0 B/s   112.0 B/s
```

If you want to display the stats for the individual services in the stack, you can run:
```
$ amp stats --service
Stats on services period=now-10m, No filter
SERVICE         CPU %   MEM USAGE   LIMIT    MEM %   IO READ   IO WRITE   NET RX       NET TX
pinger_pinger   0.00%   6.7 MB      5.9 GB   0.3%    0/s       0/s        214.0 B/s   112.0 B/s
```

If you want to display the stats for all the containers replicating the service, you can run:
```
$ amp stats --container
Stats on containers period=now-10m, No filter
CONTAINER              CPU %   MEM USAGE   LIMIT    MEM %   IO READ   IO WRITE   NET RX      NET TX
pinger_pinger_6bc02a   0.00%   1.2 MB      2.0 GB   0.1%    0/s       0/s        73.0 B/s   38.0 B/s
pinger_pinger_6e8e55   0.00%   960.0 KB    2.0 GB   0.0%    0/s       0/s        73.0 B/s   38.0 B/s
pinger_pinger_cbf5cd   0.00%   4.5 MB      2.0 GB   0.2%    0/s       0/s        72.0 B/s   38.0 B/s
```

For more detailed examples of the querying and filtering options for stats, see the [stats documentation](stats.md)

#### Organization and team management

Once you have signed up with `amp signup`, you can start creating organizations and within those organizations,
you can add members or create teams of members.

To create an organization:
```
$ amp org create
organization name: org
email: sample@org.com
Organization has been created.
```
>NOTE: The majority of the org and team commands are interactive, you can look at the org and team documentation or the `--help`
option to see which commands take arguments or options.

After creating an organization, you will be the owner of that organization and you can then add other users to that organization.

With an organization created, you can "log in" to into that organization by switching to the context of the organization using:
```
$ amp org switch org
You are now logged in as: org
```
Whilst logged in on behalf of the organization, you can deploy stacks within the context of that organization that cannot be accessed
by users outside of that organization. In addition, by creating new teams within that organization, you can delegate resources in the
form of stacks using the `amp team resource` commands.

For more information on organization and team management, see the [org documentation](org.md) and the [team documentation](team.md)
