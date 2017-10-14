#!/bin/bash

SECONDS=0

test_service_scale() {
  amp -k service scale --service pinger_pinger --replicas 2 2>/dev/null
   while true
    do
       if amp -k service ls 2>/dev/null | pcregrep -q "\s*replicated\s*2/2\s*[0-9]\s*RUNNING\s*" || [ $SECONDS -eq 5 ]
       then
               break
       fi
       sleep 1
       SECONDS=$[$SECONDS+1]
    done
}
