#!/bin/bash
# this runs the go compiler in a container to build an alpine binary (./pinger)
#
docker run -it --rm -v $PWD:/go/src/github.com/appcelerator/pinger -w /go/src/github.com/appcelerator/pinger golang:alpine go build
