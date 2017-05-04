#!/bin/bash
# simulate that a cluster creation is already in progress
cid=$(docker run --rm -d --name amp-bootstrap alpine:3.5 sleep 3 2>/dev/null)
# we should get a hint to terminate the container
amp -s localhost cluster create -t local 2>&1 | grep -q terminate
ret=$?
docker kill $cid >/dev/null 2>&1
exit $ret
