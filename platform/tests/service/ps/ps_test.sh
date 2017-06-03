#!/bin/bash

SECONDS=0

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user103 --password password --email email@user103.amp
  $amp login --name user103 --password password
}

test_stack_deploy() {
  $amp stack up -c examples/stacks/pinger/pinger.yml
}

test_service_ps_running() {
  id=$($amp service ls 2>/dev/null | grep -o -w -E '^[[:alnum:]]{25}')
  $amp service ps $id 2>/dev/null | pcregrep -q "\s*subfuzion/pinger\s*RUNNING\s*"
}

test_teardown() {
  $amp stack rm pinger
  $amp user rm user103
}
