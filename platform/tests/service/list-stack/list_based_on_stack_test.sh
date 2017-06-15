#!/bin/bash

amp="amp -s localhost"

test_stack_deploy() {
   $amp stack up -c examples/stacks/counter/counter.yml
 }

test_service_list_based_on_stack() {
  $amp service ls --stack pinger | grep -Evq "counter"
}
