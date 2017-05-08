ATOMIQ
======

## Getting started

> NOTE: by default, the CLI will connect to `cloud.appcelerator.io:50101`.
This is currently for testing and evaluation only and anything you
create will be deleted periodically over the next few weeks.

Download the latest release of the CLI for your platorm.

> *Or, if you prefer, fork this repo and build your own CLI. The entire toolchain
has been containerized, so all you need on your system is Docker 17.03 or greater. To
build the CLI, add the appropriate path under the `./bin` directory to your shell path.
Then run `ampmake buildall-cli`; this will build cross-compiled versions of the cli
(`amp`) and place them in the appropriate locations under `./bin`. (As a special case,
`make build-cli` can be used to build your OS-specific executable, if you already
have gnu make installed on your system.)*

To deploy a standard Docker Compose version 3 stackfile, use
`amp stack deploy`. There are sample stackfiles under `examples`.

## Linux

If you're running on Linux, there are a few things you must do first.

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

## Docker for Mac

Versions of Docker for Mac v.17.03 and higher work with the default settings. 

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

## Local cluster

You can create a cluster on your own system using `amp cluster create`.
Options for `create` and `update` allow you to specify the number of
managers and workers you want. Managers should be an odd number and
three is generally the ideal number. More than seven is not recommended.

To use the local cluster with the CLI, you need to specify it
with the `--server (-s)` option (ex: `amp -s localhost cluster status`).

The local cluster does not depend on the host system being a swarm manager.
The cluster nodes are created as Docker-in-Docker containers and then
joined to a swarm. The cluster includes its own registry. You can
push images that you build directly to the registry using the following sequence of steps:

    $ docker build -t foo/bar
    $ docker tag foo/bar 127.0.0.1:5000/foo/bar
    $ docker push 127.0.0.1:5000/foo/bar

The image is now available to be used by a stackfile you deploy to the
cluster. You can always use images available on Docker Hub. Support for
alternative and private registry is planned.
