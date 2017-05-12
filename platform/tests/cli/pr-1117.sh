#!/bin/bash

# verify the command 'service tasks' runs without any error

amp="amp -s localhost"
set -e
function cleanup {
  $amp user rm test1117
}
trap cleanup EXIT


$amp user signup --name test1117 --password test1117 --email test1117@email.amp
$amp login --name test1117 --password test1117
id=$(docker exec m1 docker service ls -q | head -n 1)
echo id=$id
$amp service tasks $id
