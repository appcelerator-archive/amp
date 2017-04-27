#!/bin/bash

SECONDS=0

test_stack_deploy() {
  amp -k stack up -c platform/tests/stacks/list/global.service.yml global
  amp -k stack up -c platform/tests/stacks/list/replicated.service.yml replicated
}

test_stack_starting() {
  amp -k stack ls 2>/dev/null | pcregrep -q "\s*global\s*0/1\s*[0-9]\s*STARTING\s*"
  amp -k stack ls 2>/dev/null | pcregrep -q "\s*replicated\s*0/1\s*[0-9]\s*STARTING\s*"
}

test_stack_global_running() {
  while true
  do
     if amp -k stack ls 2>/dev/null | pcregrep -q "\s*global\s*1/1\s*[0-9]\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
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
     if amp -k stack ls 2>/dev/null | pcregrep -q "\s*replicated\s*1/1\s*[0-9]\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
     then
             break
     fi
     sleep 1
     SECONDS=$[$SECONDS+1]
  done
}
