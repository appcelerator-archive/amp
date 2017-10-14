#!/bin/bash

test_stack_up() {
  amp -k stack up -c tests/cli/service/pinger.yml stackname
  amp -k stack ls 2>/dev/null | grep -Eq "stackname"
}
