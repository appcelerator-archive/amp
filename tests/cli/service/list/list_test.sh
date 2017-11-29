#!/bin/bash

SECONDS=0

test_stack_deploy() {
  amp -k stack up -c tests/cli/service/list/global.service.yml global || return 1
  amp -k stack up -c tests/cli/service/list/replicated.service.yml replicated
}

test_service_starting() {
  amp -k service ls 2>/dev/null | pcregrep -q "global_pinger\s*.*s*(STARTING|RUNNING)\s*appcelerator/pinger" || return 1
  amp -k service ls 2>/dev/null | pcregrep -q "replicated_pinger\s*.*\s*(STARTING|RUNNING)\s*appcelerator/pinger"
}

# FIXME: this test is not reliable, test condition and usage of the internal variable SECONDS
test_service_global_running() {
  while true
  do
     if amp -k service ls 2>/dev/null | grep -q "\s*global_pinger\s*global\s*1/1\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
     then
             break
     fi
     sleep 1
     SECONDS=$[$SECONDS+1]
  done
}

# FIXME: this test is not reliable, test condition and usage of the internal variable SECONDS
test_service_replicated_running() {
  while true
  do
     if amp -k service ls 2>/dev/null | grep -q "\s*replicated_pinger\s*replicated\s*1/1\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
     then
             break
     fi
     sleep 1
     SECONDS=$[$SECONDS+1]
  done
}
