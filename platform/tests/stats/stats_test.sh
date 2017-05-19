#!/bin/bash

set -e

amp="amp -s localhost"
user=test1705191557

test_setup() {
  $amp user signup --name $user --password pwd$user --email $user@email.amp
  $amp login --name $user --password pwd$user
}

test_main() {
  res=$($amp stats | wc -l)
  echo $res
  if [ "$res" -lt 1 ]; then
     exit 1
  fi
}

test_teardown() {
  $amp user rm $user
}
