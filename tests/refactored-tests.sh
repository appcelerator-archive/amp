#!/bin/bash

NETWORK=swarmnet
TAG=local
SERVER_PORT=50101

echo "Starting the amplifier service... "
docker service create --name amplifier --network $NETWORK appcelerator/amplifier:$TAG
if [ $? -ne 0 ]; then
  echo "Failed to start amplifier service"
  exit 1
fi
echo "Waiting for the amplifier service to be up... "
maxretries=10
retries=0
while [ $retries -le $maxretries ]; do
  docker service ps  amplifier | awk '{print $6}' | grep -qw Running && break
  echo -n "."
  sleep 1
  ((retries+1))
done
if [ $retries -eq $maxretries ]; then
  echo "amplifier failed to start in a sensible time"
  docker service rm amplifier
  exit 1
fi
echo
echo "Connecting to amplifier with the CLI... "
# test the CLI and the connection to the server
# if connection fails, the container will return a non zero code
docker run --rm --name cli --network $NETWORK appcelerator/amp:$TAG --server amplifier:$SERVER_PORT version
if [ $? -ne 0 ]; then
  echo "Failed to connect"
  docker service rm amplifier
  exit 1
fi
docker service rm amplifier
echo "Passed"
