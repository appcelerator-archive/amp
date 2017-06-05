#!/bin/bash

export GOPATH=/go
NATS_HOST=${NATS_HOST:-nats}
NATS_PORT=4222
TMPFILE=$(mktemp)
MAX=2000

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

echo -n "run a subscriber... "
go run /bin/nats-sub.go -s nats://$NATS_HOST:$NATS_PORT msg.test > $TMPFILE 2>&1 &
r=1
i=0
while [[ $r -ne 0 ]]; do
  sleep 1
  grep -q "Listening on \[msg.test\]" $TMPFILE
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


echo -n "publish messages... "
r=0
sleep 1
for i in $(seq $MAX); do
  /go/bin/nats-pub -s nats://$NATS_HOST:$NATS_PORT msg.test "test message $i" 2>&1 | grep -q "Published \[msg.test\]"
  r=$?
  if [[ $r -ne 0 ]]; then
    break
  fi
  if [[ $((i % 200)) -eq 0 ]]; then echo -n "+"  ; fi
done
if [[ $r -ne 0 ]]; then
  echo
  echo "failed ($i/$MAX msg sent)"
  exit 1
fi
echo "[OK]"
sleep 10
echo -n "receive messages... "
n=$(egrep -c "Received on \[msg.test\].*:.*'test message .*'" $TMPFILE)
if [[ $n -ne $MAX ]]; then
  echo
  echo " failed ($n/$MAX messages received)"
  cat $TMPFILE
  exit 1
fi
echo " [OK]"

echo -n "benchmark...        "
go run /bin/nats-bench.go -s nats://$NATS_HOST:$NATS_PORT -np 10 -ns 1  -n $MAX -ms 10 msg.bench > $TMPFILE 2>&1
if [[ $? -ne 0 ]]; then
  echo
  echo "failed"
  exit 1
fi
echo " [OK]"
cat $TMPFILE

echo "all tests passed successfully"
