#!/bin/bash

test_stack_deploy() {
  amp -k stack up -c tests/cli/service/pinger.yml another
}

test_service_list_based_on_stack() {
  amp -k service ls --stack pinger | grep -Evq "another"
}
