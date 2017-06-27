# Envoy

Runs a container in a target cluster environment to perform various tasks, which may include
set up, smoke tests, etc.

`Dockerfile.compiler` builds an image for compiling the Go source files for
the final target image (`appcelerator/amp-envoy`).

`Dockerfile` builds the image for running a container in the target cluster.

# Build

`make compiler` builds an Alpine image (`appcelerator/amp-envoy-compiler`) with
a Go compiler, the AWS SDK package, and other necessary packages to to build 
the `appcelerator/amp-envoy` image.

An automated build for the repo also creates the `appcelerator/amp-envoy-compiler`
image on [Docker Hub](https://hub.docker.com/r/appcelerator/amp-envoy-compiler/).

`make` (or `make build`) builds `appcelerator/amp-envoy`.

`make clean` removes the target binary (`envoy.alpine`) that is created by the
compiler to be copied into the `appcelerator/amp-envoy` image when building it.

`make test` uses `appcelerator/amp-envoy` to create, update, and remove test stacks.

