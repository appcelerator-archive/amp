#!/bin/bash

docker run --rm --network monit appcelerator/alpine:3.6.0 curl -sf http://amplifier:5100/metrics
