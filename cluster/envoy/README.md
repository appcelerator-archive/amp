# Envoy

Runs a container in a target cluster environment to perform various tasks, which may include
set up, smoke tests, etc.

## Build

* `make` (or `make build`) - builds target `bin/aws`
* `make clean` - removes `bin/aws`
* `make test` - runs go tests that create, update, and remove test stacks on AWS.
* `make build-image` - builds `appcelerator/amp-aws` using the multi-stage `Dockerfile`
