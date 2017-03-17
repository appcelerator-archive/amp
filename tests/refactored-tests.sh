#!/bin/bash

NETWORK=hostnet # This is not configurable yet
TAG=local
SERVER_PORT=50101

echo "Starting the swarm ... "
docker network create $NETWORK 2>/dev/null || true
amp start

echo "Building and publishing amplifier image..."
docker push localhost:5000/appcelerator/amplifier:local

echo "Starting the amplifier stack..."
docker run -it --rm --network=$NETWORK -v $PWD/stacks:/stacks docker --host=m1 stack deploy -c /stacks/amplifier.stack.yml amplifier
if [ $? -ne 0 ]; then
  echo "Failed to start amplifier stack"
  exit 1
fi

#echo "Waiting for amplifier to be reachable ..."
#maxretries=30
#retries=0
#while [ $retries -le $maxretries ]; do
#  docker run --rm --name cli --network $NETWORK appcelerator/amp:$TAG --server m1:50101 version &> /dev/null && break
#  echo -n "."
#  sleep 1
#  ((retries++))
#done
#echo
#
#if [ $retries -gt $maxretries ]; then
#  echo " amplifier failed to start in a sensible time"
#  docker run -it --rm --network=hostnet docker --host=m1 service ls
#  exit 1
#fi

echo "Passed"
