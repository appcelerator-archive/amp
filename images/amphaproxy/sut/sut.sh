#!/bin/bash

HAPROXY_HOST=${HAPROXY_HOST:-haproxy}
STATS_PORT=8082
STATS_URL="admin?stats"

echo -n "test Haproxy Availability... "
r=1
i=0
while [[ $r -ne 0 ]]; do
  sleep 1
  timeout -t 2 curl -I "$HAPROXY_HOST:$STATS_PORT/$STATS_URL" 2>/dev/null | grep -q "HTTP/1.1 200 OK"
  r=$?
  ((i++))
  if [[ $i -gt 12 ]]; then break; fi
  echo -n "+"
done
if [[ $r -ne 0 ]]; then
  echo
  echo "$HAPROXY_HOST:$STATS_PORT failed"
  curl -I "$HAPROXY_HOST:$STATS_PORT"
  exit 1
fi
echo " [OK]"

echo -n "test Haproxy stats... "
r=1
timeout -t 2 curl "$HAPROXY_HOST:$STATS_PORT/$STATS_URL;csv" 2>/dev/null | grep -q "http-in"
r=$?
if [[ $r -ne 0 ]]; then
  echo
  echo "can't find http-in in $HAPROXY_HOST:$STATS_PORT"
  curl "$HAPROXY_HOST:$STATS_PORT;csv"
  exit 1
fi
echo " [OK]"

echo "all tests passed successfully"