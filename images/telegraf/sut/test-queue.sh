#!/bin/bash

DATA_DIR=/var/data
OUTPUT_FILE=$DATA_DIR/output.dat
DOCKER_SOCKET=/var/run/docker-host.sock
CONSUMER=telegraf-consumer
AGENT=telegraf-agent
nm=3000
waitforqueuer=${1:-2}

_docker_logs(){
  ct=$(docker --host=unix://$DOCKER_SOCKET ps -a | grep dockertelegraf_telegraf | awk '{print $1}' | xargs -I {} docker --host=unix://$DOCKER_SOCKET logs {})
}

# cleanup
if [[ -f "$OUTPUT_FILE" ]]; then
  echo "INFO - resetting the test file"
  > "$OUTPUT_FILE"
else
  echo "INFO - no test file found, yet"
fi

echo -n "wait for telegraf to be available...     "
r1=""
r2=""
i=0
while [[ -z "$r1" || -z "$r2" ]]; do
  r1=$(dig +short $AGENT)
  r2=$(dig +short $CONSUMER)
  ((i++))
  if [[ $i -gt 30 ]]; then
    echo
    echo "failed"
    dig +short $AGENT
    dig +short $CONSUMER
    exit 1
  fi
  sleep 1
done
echo "[OK]"

echo -n "wait $waitforqueuer sec for queuer...    "
sleep $waitforqueuer
echo "[OK]"

echo -n "test send to tcp listener...             "
for i in $(seq $nm); do
  echo '{"msg": '$i', "ts": '$(date +%s)'}' | nc -w 1 $AGENT 8094
  if [[ $? -ne 0 ]]; then
    echo
    echo "failed"
    _docker_logs
    exit 1
  fi
done
echo "[OK] ($nm msg)"

echo -n "test output file...                      "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  if [[ -f "$OUTPUT_FILE" ]]; then
    r=true
    break
  fi
  if [[ $i -gt 6 ]]; then break; fi
  ((i++))
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "telegraf didn't write anything"
  exit 1
fi
echo "[OK] ($i sec)"

echo -n "test tcp listener measurement data...    "
sleep 5
n=$(grep -c "^tcp_listener," "$OUTPUT_FILE")
if [[ $n -ne $nm ]]; then
  echo
  echo "failed ($n/$nm msg)"
  _docker_logs
  cat $OUTPUT_FILE
  exit 1
fi
echo "[OK] ($nm msg)"

echo "cleaning up output file"
> "$OUTPUT_FILE"
rm "$OUTPUT_FILE"

echo "all tests passed successfully"
