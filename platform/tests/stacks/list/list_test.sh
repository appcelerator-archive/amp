#!/bin/bash

SECONDS=0

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user102 --password password --email email@user102.amp
  $amp login --name user102 --password password
}

test_stack_deploy() {
  $amp stack up -c platform/tests/stacks/list/global.service.yml global
  $amp stack up -c platform/tests/stacks/list/replicated.service.yml replicated
}

test_stack_starting() {
  $amp stack ls 2>/dev/null | pcregrep -q "\s*global\s*0/1\s*STARTING\s*"
  $amp stack ls 2>/dev/null | pcregrep -q "\s*replicated\s*0/1\s*STARTING\s*"
}

test_stack_global_running() {
  while true
  do
     if $amp stack ls 2>/dev/null | pcregrep -q "\s*global\s*1/1\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
     then
             break
     fi
     sleep 1
     SECONDS=$[$SECONDS+1]
  done
}

test_stack_replicated_running() {
  while true
  do
     if $amp stack ls 2>/dev/null | pcregrep -q "\s*replicated\s*1/1\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
     then
             break
     fi
     sleep 1
     SECONDS=$[$SECONDS+1]
  done
}

test_teardown() {
  $amp stack rm global
  $amp stack rm replicated
  $amp user rm user102
}
