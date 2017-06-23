## Getting started

### Prequisites

We are currently targeting the next stable release of Docker, anticipated to be Docker 17.06.
For now, we recommend installing the latest Docker CE edge channel release (17.05) on your system:

 * [Docker for Mac](https://docs.docker.com/docker-for-mac/install/)
 * [Docker for Windows](https://docs.docker.com/docker-for-windows/install/)
 * [Docker for Linux](https://docs.docker.com/engine/installation/) (see specific distribution under the sidebar).
 
#### Linux Node:

If you're running on Linux, there are a few steps you must perform for your system in addition to installing Docker.

Check if you have IPVS kernel modules loaded:

    $ lsmod | grep ip_vs
    # verify ip_vs and ip_vs_rr
    #
    # if necessary, load IPVS kernel modules (ip_vs, ip_vs_rr) for swarm mode networking and load-balancing
    # Note: these modules are needed to run a local swarm using Docker-in-Docker on the host,
    # but this step is not necessary if the host is already a member of a swarm.
    # Note: loading ip_vs_rr will also ensure ip_vs is loaded.
    $ sudo modprobe ip_vs_rr

Increase virtual memory (needed for Elasticsearch):

    $ sudo sysctl -w vm.max_map_count=262144

To make these changes permanent, you can run the following and reboot:

    $ echo "ip_vs_rr" | sudo tee -a /etc/modules
    $ echo "vm.max_map_count = 262144" | sudo tee -a /etc/sysctl.conf

### Get the AMP CLI

Download the latest release of the CLI for your platform:

https://github.com/appcelerator/amp/releases

Place the `amp` binary from the release archive into some location in your system path. For example:

    $ sudo mv ./amp /usr/local/bin

> *Or, if you prefer, fork this repo and build your own CLI. The entire toolchain
has been containerized, so all you need on your system is Docker 17.03 or greater. To
build the CLI, add the appropriate path under the `./bin` directory to your shell path.
Then run `ampmake buildall-cli`; this will build cross-compiled versions of the cli
(`amp`) and place them in the appropriate locations under `./bin`. (As a special case,
`make build-cli` can be used to build your OS-specific executable, if you already
have gnu make installed on your system.)*

### Signup and login on hosted AMP

    $ amp user signup # check your email to activate your account
    $ amp login

### Create your own local AMP cluster

To use amp to deploy and test stacks locally, you can use the `--server|-s` option
with the amp command and specify `localhost`.

    $ amp -s localhost cluster create
    $ amp -s localhost user signup # Ignore the verification email
    $ amp -s localhost login


##### Tip:

If you will be primarily working with `localhost`, then you might want to
set an alias to make this your default:

```sh
    $ alias amp="$PWD/amp --server localhost"
    $ amp version  # => localhost
    $ \amp version # => use backslash to override alias and use default cloud.appcelerator.io
    $ unalias amp  # restore default (cloud.appcelerator.io)
```

### Deploy a stack

> NOTE: by default, the CLI will connect to `cloud.appcelerator.io:50101`.
This is currently for testing and evaluation only and anything you
create will be deleted periodically over the next few weeks.

To deploy a standard Docker Compose version 3 stackfile, use
`amp stack deploy`. There are sample stackfiles under `examples`.


Here is a simple example:

```sh
    $ curl -O https://raw.githubusercontent.com/appcelerator/amp/master/examples/stacks/pinger/pinger.yml
    $ amp -s localhost stack deploy -c pinger.yml
    $ amp -s localhost stack ls
    $ amp -s localhost service logs pinger
    $ curl http://pinger.examples.local.appcelerator.io/ping
```

Or browse to: https://pinger.examples.local.appcelerator.io/ping

## Monitoring

The `amp stats` and `amp logs` commands provide rich filtered
query support for both realtime feeds and historical queries.
You can connect to the hosted dashboard at https://dashboard.cloud.appcelerator.io/.
If you create a local cluster, you can connect to it at:
http://localhost:50106.

## Organization and teams

To test the organization and teams features,
use the CLI to create your ATOMIQ ID (`amp user signup`), then
log in (`amp login`).
