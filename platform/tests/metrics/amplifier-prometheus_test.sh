#!/bin/bash

docker run --rm --network ampnet appcelerator/alpine:3.6.0 curl -sf http://amplifier:5100/metrics
