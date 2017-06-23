#!/bin/bash

SECONDS=0

test_service_scale() {
  amp service scale --service pinger_pinger --replicas 4 2>/dev/null
   while true
    do
       if amp service ls 2>/dev/null | pcregrep -q "\s*replicated\s*4/4\s*[0-9]\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
       then
               break
       fi
       sleep 1
       SECONDS=$[$SECONDS+1]
    done
}
