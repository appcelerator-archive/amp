# Local cluster Plugin

This is the plugin for creating and initializing a swarm on the local machine.

`Dockerfile.compiler` builds the image for compiling the Go source files for
the plugin.

`Dockerfile` builds the image for the actual plugin that will be used by the
AMP cli.

For more details about the design and use, see the
[wiki](https://github.com/appcelerator/amp/wiki/AMP-Clusters).

### Options

The parameters available for the local plugins are:

 * advertise-addr (interface for the Swarm initialization, default=eth0)

Parameters should be passed as followed:

    --advertise-addr en0 ...

## Prerequisites

The host should have these prerequisites:

* a proper storage driver (AUFS, overlay, overlay2, devicemapper on thin-pool)
* max map count limit at least equals to 262144 (Linux only)
* enough system resources (4GB memory, 20GB storage)

## Trying it out

From the `cluster/plugin/local` directory, run the following:

    $ make vendor
    $ make image
    $ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/docker:/var/run/docker appcelerator/amp-local init

Verify output similar to the following:

```
$ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v /var/run/docker:/var/run/docker appcelerator/amp-local init

2017/07/29 05:16:49 ampctl (version: latest, build: 43437f20)
2017/07/29 05:16:49 Version test: PASS
2017/07/29 05:16:49 Labels test: PASS
```


### Update Swarm

N/A

### Destroy Swarm

Execute the `destroy` command:

    $ docker run -it --rm appcelerator/amp-local destroy

## Tests

    $ make test
