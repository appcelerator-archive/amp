# stack

Currently, Docker stack support is not part of the Docker engine; instead, stack
functionality is embedded in the current Docker CLI stack command implementation. 

The experimental project is a precursor to working with the actual `docker/docker`,
`docker/cli`, and `docker/docker-ce` repos and preparing formal pull requests.

The goal of this project is to create standalone packages for stack operations and
compose file support, add gRPC / protobuf support, and then add support for exposing
stack operations as part of the engine REST API.

You can generate a Docker CLI (that only supports the `stack` command) to demo.

## make

* `make` | `make build` - builds target `bin/stackcli`
* `make clean` - remove `bin/stackcli`
* `make test` - run tests
* `make image` - build `subfuzion/stackcli` using the multi-stage `Dockerfile`

You can run `bin/stackcli` directly or in a container, like the following:

    $ docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock subfuzion/stackcli stack ls
