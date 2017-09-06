#!/bin/bash

_timeout=10
for i in $(seq 3); do
  docker run --rm --network ampnet appcelerator/alpine:3.6.0 curl -sfm $_timeout http://amplifier:5100/metrics && break
done
