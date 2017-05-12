#!/bin/bash

# verify the command 'stack services' runs without any error

amp="amp -s localhost"
set -e
function cleanup {
  $amp user rm test1112
}
trap cleanup EXIT

$amp user signup --name1112 test --password test1112 --email test1112@email.amp
$amp login --name test1112 --password test1112
$amp stack up -c examples/stacks/pinger/pinger.yml pinger1112
sleep 1
$amp stack services pinger1112 | grep -q pinger1112_pinger
ret=$?
$amp stack rm pinger1112
exit $ret
