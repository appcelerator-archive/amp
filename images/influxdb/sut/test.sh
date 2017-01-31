#!/bin/bash

INFLUXDB_HOST=${INFLUXDB_HOST:-influxdb}

echo -n "test Influxdb Availability... "
# give time to influxdb to be up
r=1
i=0
while [[ $r -ne 0 ]]; do
  sleep 1
  curl -I $INFLUXDB_HOST:8083 2>/dev/null | grep -q "HTTP/1.1 200 OK"
  r=$?
  ((i++))
  if [[ $i -gt 60 ]]; then break; fi
  echo -n "+"
done
if [[ $r -ne 0 ]]; then
  echo
  echo "Influxdb:8083 failed"
  curl -I $INFLUXDB_HOST:8083
  exit 1
fi
echo " [OK]"

echo -n "test Influxdb ping... "
curl -I $INFLUXDB_HOST:8086/ping 2>/dev/null | grep -q "HTTP/1.1 204 No Content"
if [[ $? -ne 0 ]]; then
  echo
  echo "Influxdb:8086 ping failed"
  exit 1
fi
echo "[OK]"

echo -n "test measurements presence... "
r="false"
i=0
# give time to telegraf to send data
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  r=$(curl -GET "http://$INFLUXDB_HOST:8086/query" --data-urlencode "db=telegraf" --data-urlencode "q=SHOW MEASUREMENTS" 2>/dev/null | jq -r '.results[0] | has("series")')
  ((i++))
  if [[ $i -gt 30 ]]; then break; fi
  echo -n "+"
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "Influxdb telegraf db has no measurements ($r)"
  curl -GET "http://$INFLUXDB_HOST:8086/query" --data-urlencode "db=telegraf" --data-urlencode "q=SHOW MEASUREMENTS"
  exit 1
fi
echo " [OK]"

echo -n "test Docker cpu measurement presence... "
r=$(curl -GET "http://$INFLUXDB_HOST:8086/query" --data-urlencode "db=telegraf" --data-urlencode "q=SELECT usage_total FROM docker_container_cpu limit 1" 2>/dev/null | jq -r '.results[0] | has("series")')
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "Influxdb telegraf db has no docker_container_cpu measurements ($r)"
  curl -GET "http://$INFLUXDB_HOST:8086/query" --data-urlencode "db=telegraf" --data-urlencode "q=SELECT usage_total FROM docker_container_cpu limit 1"
  exit 1
fi
echo "[OK]"

echo "all tests passed successfully"
