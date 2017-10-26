#!/bin/bash

test_stack_up() {
  amp -k stack up -c tests/cli/stack/pinger.yml stackname
  amp -k stack ls 2>/dev/null | grep -Eq "stackname"
  amp -k stack rm stackname
}

test_stack_up_env() {
  export VAR=2.0
  amp -k stack up -c tests/cli/stack/pinger.yml stackenv
  amp -k service ls 2>/dev/null | grep -Eq "2.0"
  amp -k stack rm stackenv
}
