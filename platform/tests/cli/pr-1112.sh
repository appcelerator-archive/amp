#!/bin/bash

# verify the command 'stack services' runs without any error 

set -e

amp stack up -c examples/stacks/pinger/pinger.yml pinger -s localhost
sleep 1
amp stack services pinger -s localhost | grep -q pinger_pinger
ret=$?
amp stack rm pinger -s localhost
exit $ret
