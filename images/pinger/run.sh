#!/bin/bash
# runs pinger in a container
#
docker run -it --rm -p 3000:3000 -e PORT=3000 --name pinger appcelerator/pinger
