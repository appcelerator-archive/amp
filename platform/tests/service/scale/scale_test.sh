#!/bin/bash

SECONDS=0

test_setup() {
  amp="amp -s localhost"
  $amp user signup --name user105 --password password --email email@user105.amp
  $amp login --name user105 --password password
}

test_stack_deploy() {
  $amp stack up -c examples/stacks/pinger/pinger.yml
}

test_service_scale() {
  $amp service scale --service pinger_pinger --replicas 4 2>/dev/null
   while true
    do
       if $amp service ls 2>/dev/null | pcregrep -q "\s*replicated\s*4/4\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
       then
               break
       fi
       sleep 1
       SECONDS=$[$SECONDS+1]
    done
}

test_teardown() {
  $amp stack rm pinger
  $amp user rm user105
}
