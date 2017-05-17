#!/bin/bash

set -e

amp="amp -s localhost"
user=test1117

test_setup() {
  $amp user signup --name $user --password pwd$user --email $user@email.amp
  $amp login --name $user --password pwd$user
}

test_main() {
  id=$(docker exec m1 docker service ls -q | head -n 1)
  $amp service tasks $id | grep 'DESIRED STATE'
}

test_teardown() {
  $amp user rm $user
}
