# ampadmin

Runs a container in a target cluster environment to perform various tasks, which may include
set up, smoke tests, etc.

## Build

* `make` (or `make build`) - builds target `bin/ampadmin`
* `make clean` - removes `bin/aws`
* `make test` - runs tests
* `make build-image` - builds `appcelerator/ampadmin` using the multi-stage `Dockerfile`
