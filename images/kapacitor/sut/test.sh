#!/bin/bash

KAPACITOR_HOST=${KAPACITOR_HOST:-kapacitor}

# wait for init of kapacitor, to avoid testing the pre config start
sleep 5
echo -n "test kapacitor ping... "
i=0
r=1
while [[ $r -ne 0 ]]; do
  ((i++))
  sleep 1
  curl -I $KAPACITOR_HOST:9092/kapacitor/v1/ping 2>/dev/null | grep -q "HTTP/1.1 204 No Content"
  r=$?
  if [[ $i -gt 25 ]]; then break; fi
  echo -n "+"
done
if [[ $r -ne 0 ]]; then
  echo
  echo "failed"
  curl -I $KAPACITOR_HOST:9092/kapacitor/v1/ping
  echo "Running containers:"
  docker ps
  ci=$(docker ps -a | grep /influxdb | head -1 | awk '{print $1}')
  echo "logs from influxdb $ci:"
  docker logs $ci
  ck=$(docker ps -a | grep /kapacitor | head -1 | awk '{print $1}')
  echo "logs from kapacitor $ck:"
  docker logs $ck
  exit 1
fi
echo "[OK]"

echo -n "test tasks list... "
i=0
r=0
while [[ $r -lt 2 ]]; do
  ((i++))
  sleep 1
  r=$(curl $KAPACITOR_HOST:9092/kapacitor/v1/tasks 2>/dev/null | jq -r '.tasks | length')
  if [[ $i -gt 5 ]]; then break; fi
  echo -n "+"
done
if [[ $r -lt 2 ]]; then
  echo
  echo "failed ($r)"
  curl $KAPACITOR_HOST:9092/kapacitor/v1/tasks
  exit 1
fi
echo "($r) [OK]"

echo -n "test subscriptions... "
i=0
nb=0
while [[ $nb -eq 0 ]]; do
  ((i++))
  sleep 1
  r=$(curl $KAPACITOR_HOST:9092/kapacitor/v1/debug/vars 2>/dev/null | jq '.kapacitor | map(select(.name == "ingress") + select(.tags.database == "telegraf") + select(.tags.retention_policy == "default"))')
  nb=$(echo $r | jq 'length')
  if [[ $i -gt 50 ]]; then break; fi
  echo -n "+"
done
if [[ $nb -lt 1 ]]; then
  echo
  echo "failed ($nb subscriptions)"
  curl $KAPACITOR_HOST:9092/kapacitor/v1/debug/vars 2>/dev/null
  echo "Running containers:"
  docker ps
  ci=$(docker ps -a | grep /influxdb | head -1 | awk '{print $1}')
  echo "logs from influxdb $ci:"
  docker logs $ci
  ck=$(docker ps -a | grep /kapacitor | head -1 | awk '{print $1}')
  echo "logs from kapacitor $ck:"
  docker logs $ck
  exit 1
fi
echo "($nb) [OK]"

echo -n "test subscription data... "
i=0
nb=0
while [[ $nb -eq 0 ]]; do
  ((i++))
  sleep 1
  nb=$(curl $KAPACITOR_HOST:9092/kapacitor/v1/debug/vars 2>/dev/null | jq '.kapacitor | map(select(.name == "ingress") + select(.tags.database == "telegraf") + select(.tags.retention_policy == "default"))' | jq -r 'map(select(.tags.measurement == "cpu")) | .[0].values.points_received')
  if [[ $i -gt 25 ]]; then break; fi
  echo -n "+"
done
if [[ $nb -lt 1 ]]; then
  echo
  echo "failed ($nb subscriptions)"
  curl $KAPACITOR_HOST:9092/kapacitor/v1/debug/vars 2>/dev/null
  exit 1
fi
echo "($nb points) [OK]"

echo "all tests passed successfully"
