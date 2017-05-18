#!/bin/bash

set -e

amp="amp -s localhost"
user=test1117

test_setup() {
  $amp user signup --name $user --password pwd$user --email $user@email.amp
  $amp login --name $user --password pwd$user
}

test_main() {
  # TODO: enable it once we have service ls support in the CLI
  return 0
  #id=$(docker exec m1 docker service ls -q | head -n 1)
  #$amp service tasks $id | grep 'DESIRED STATE'
}

test_teardown() {
  $amp user rm $user
}
