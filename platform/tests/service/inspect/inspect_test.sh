#!/bin/bash

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user104 --password password --email email@user104.amp
}

test_stack_deploy() {
  $amp stack up -c examples/stacks/pinger/pinger.yml
}

test_service_inspect() {
  $amp service inspect pinger_pinger 2>/dev/null | pcregrep -q "pinger"
}

test_teardown() {
  $amp stack rm pinger
  $amp user rm user104
}
