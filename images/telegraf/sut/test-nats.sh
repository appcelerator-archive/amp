#!/bin/bash

DATA_DIR=/var/data
OUTPUT_FILE=$DATA_DIR/output.dat
DOCKER_SOCKET=/var/run/docker-host.sock
CONSUMER=telegraf-consumer
AGENT=telegraf-agent
nm=3000
waitforqueuer=${1:-2}
export GOPATH=/go
NATS_HOST=${NATS_HOST:-nats}
NATS_PORT=4222
TMPFILE=$(mktemp)

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

echo -n "test NATS Availability... "
r=1
i=0
while [[ $r -ne 0 ]]; do
  sleep 1
  curl -sI $NATS_HOST:8222 2>/dev/null | grep -q "HTTP/1.1 200 OK"
  r=$?
  ((i++))
  if [[ $i -gt 10 ]]; then break; fi
  echo -n "+"
done
if [[ $r -ne 0 ]]; then
  echo
  echo "$NATS_HOST:8222 failed"
  curl -I $NATS_HOST:8822
  exit 1
fi
echo " [OK]"

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

echo -n "input plugin test, publish messages... "
r=0
for i in $(seq $nm); do
  #go run /bin/nats-pub.go -s nats://$NATS_HOST:$NATS_PORT telegraf.in "fake_measurement,host=server$i value=$i 1422568543702900257" 2>&1 | grep -q "Published \[telegraf.in\]"
  # compiled will be faster
  /go/bin/nats-pub -s nats://$NATS_HOST:$NATS_PORT telegraf.in "fake_measurement,host=server$i value=$i 1422568543702900257" 2>&1 | grep -q "Published \[telegraf.in\]"
  r=$?
  if [[ $r -ne 0 ]]; then
    break
  fi
done
if [[ $r -ne 0 ]]; then
  echo
  echo "failed ($i/$nm msg sent)"
  exit 1
fi
echo " [OK] ($nm messages)"

echo -n "input plugin test - test output file...                      "
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

echo -n "input plugin test - test measurement data...    "
sleep 5
n=$(grep -c "^fake_measurement," "$OUTPUT_FILE")
if [[ $n -ne $nm ]]; then
  echo
  echo "failed ($n/$nm msg)"
  _docker_logs
  exit 1
fi
echo "[OK] ($nm msg)"
#grep "^fake_measurement," "$OUTPUT_FILE"

echo -n "output plugin test, run subscriber... "
go run /bin/nats-sub.go -s nats://$NATS_HOST:$NATS_PORT telegraf.out > $TMPFILE 2>&1  &
r=1
i=0
while [[ $r -ne 0 ]]; do
  sleep 1
  grep -q "Listening on \[telegraf.out\]" $TMPFILE
  r=$?
  ((i++))
  if [[ $i -gt 3 ]]; then break; fi
  echo -n "+"
done
if [[ $r -ne 0 ]]; then
  echo
  echo "failed"
  cat $TMPFILE
  exit 1
fi
echo " [OK]"

echo -n "output plugin test - send to tcp listener...             "
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

# @todo
echo -n "output plugin test - read messages... "
sleep 2
n=$(egrep -c "Received on \[telegraf.out\].*:.*'tcp_listener,.* msg=.*'" $TMPFILE)
if [[ $n -ne $nm ]]; then
  echo
  echo " failed ($n/$nm messages received)"
  cat $TMPFILE
  exit 1
fi
echo " [OK]"
#egrep "Received on \[telegraf.out\].*:.*'tcp_listener,.* msg=.*'" $TMPFILE

echo "cleaning up output file"
> "$OUTPUT_FILE"
rm "$OUTPUT_FILE"

echo "all tests passed successfully"
