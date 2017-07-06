# ampagent

Runs a container in a target cluster environment for monitoring various
swarm events, as well as to perform various adminstrative tasks, which may include
running amp initialization and system checks, etc.

## Build

* `make` | `make build` - builds target `bin/ampctl`
* `make clean` - removes `bin/ampctl`
* `make test` - runs tests
* `make image` - builds `appcelerator/ampaagent` using the multi-stage `Dockerfile`
* `make run` - runs `ampctl monitor` in a container (to mount swarm control socket)

