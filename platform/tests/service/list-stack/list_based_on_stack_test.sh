#!/bin/bash

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user106 --password password --email email@user106.amp
  $amp login --name user106 --password password
}

test_stack_deploy() {
   $amp stack up -c examples/stacks/pinger/pinger.yml
   $amp stack up -c examples/stacks/counter/counter.yml
 }

test_service_list_based_on_stack() {
  $amp service ls --stack pinger | grep -Evq "counter"
}

test_teardown() {
  $amp stack rm pinger
  $amp stack rm counter
  $amp user rm user106
}
