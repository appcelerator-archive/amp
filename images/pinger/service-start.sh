#!/bin/bash
# runs pinger as a service (assumes host is a swarm manager), ie:
# docker swarm init
#
docker service create -p 3000:3000 --name pinger appcelerator/pinger
