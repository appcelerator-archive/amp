#!/bin/bash

set -e

amp="amp -s localhost"
user=test1196

test_setup() {
  $amp user signup --name $user --password pwd$user --email $user@email.amp
  $amp login --name $user --password pwd$user
}

test_main() {
  $amp logs -m --infra | grep -q 'timestamp:'
}

test_teardown() {
  $amp user rm $user
}
