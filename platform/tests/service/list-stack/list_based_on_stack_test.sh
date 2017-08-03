#!/bin/bash

test_stack_deploy() {
   amp -k stack up -c examples/stacks/counter/counter.yml
 }

test_service_list_based_on_stack() {
  amp -k service ls --stack pinger | pcregrep -vq "counter"
}
