#!/bin/bash

DATA_DIR=/var/data
OUTPUT_FILE=$DATA_DIR/output.dat
DOCKER_SOCKET=/var/run/docker-host.sock

_docker_logs(){
  ct=$(docker --host=unix://$DOCKER_SOCKET ps | grep /telegraf | awk '{print $1}')
  echo "logs from telegraf $ct:"
  docker --host=unix://$DOCKER_SOCKET logs $ct
  cat $OUTPUT_FILE
  > "$OUTPUT_FILE"
  rm "$OUTPUT_FILE"
}

# cleanup
if [[ -f "$OUTPUT_FILE" ]]; then
  echo "INFO - resetting the test file"
  > "$OUTPUT_FILE"
else
  echo "INFO - no test file found, yet"
fi

echo -n "test output file...                      "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  if [[ -f "$OUTPUT_FILE" ]]; then
    grep -q "docker" "$OUTPUT_FILE"
    if [[ $? -eq 0 ]]; then
      r="true"
    fi
  fi
  ((i++))
  if [[ $i -gt 60 ]]; then break; fi
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "telegraf didn't write anything"
  _docker_logs
  exit 1
fi
echo "[OK] ($i sec)"

echo -n "test docker measurement...               "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  grep -q "^docker," "$OUTPUT_FILE"
  if [[ $? -eq 0 ]]; then
    r="true"
  fi
  if [[ $i -gt 45 ]]; then break; fi
  ((i++))
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "failed (after $i sec)"
  _docker_logs
  exit 1
fi
echo "[OK] ($i sec)"

echo -n "test docker_container_cpu measurement... "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  grep -q "^docker_container_cpu," "$OUTPUT_FILE"
  if [[ $? -eq 0 ]]; then
    r="true"
  fi
  if [[ $i -gt 15 ]]; then break; fi
  ((i++))
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "failed (after $i sec)"
  _docker_logs
  exit 1
fi
echo "[OK] ($i sec)"
echo -n "test docker_container_mem measurement... "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  grep -q "^docker_container_mem," "$OUTPUT_FILE"
  if [[ $? -eq 0 ]]; then
    r="true"
  fi
  if [[ $i -gt 15 ]]; then break; fi
  ((i++))
done
if [[ $r -ne 0 ]]; then
  echo
  echo "failed (after $i sec)"
  _docker_logs
  exit 1
fi
echo "[OK] ($i sec)"
echo -n "test docker_container_net measurement... "
r="false"
i=0
while [[ "x$r" != "xtrue" ]]; do
  sleep 1
  grep -q "^docker_container_net," "$OUTPUT_FILE"
  if [[ $? -eq 0 ]]; then
    r="true"
  fi
  if [[ $i -gt 45 ]]; then break; fi
  ((i++))
done
if [[ "x$r" != "xtrue" ]]; then
  echo
  echo "failed (after $i sec)"
  _docker_logs
  exit 1
fi
echo "[OK] ($i sec)"
echo -n "test docker_container_blkio measurement... "
grep -q "^docker_container_blkio," "$OUTPUT_FILE"
if [[ $? -ne 0 ]]; then
  echo "[no data, ignore]"
else
echo "[OK]"
fi
echo -n "test net measurement...                  "
grep -q "^net," "$OUTPUT_FILE"
if [[ $? -ne 0 ]]; then
  echo
  echo "failed"
  _docker_logs
  exit 1
fi
echo "[OK]"
echo -n "test send to tcp listener...             "
echo '{"count": 1, "status": 0}' | nc telegraf 8094
if [[ $? -ne 0 ]]; then
  echo
  echo "failed (send)"
  _docker_logs
  exit 1
fi
echo "[OK]"
echo -n "test tcp listener measurement data...    "
sleep 3
grep -q "^tcp_listener," "$OUTPUT_FILE"
if [[ $? -ne 0 ]]; then
  echo
  echo "failed (no data)"
  _docker_logs
  exit 1
fi
echo "[OK]"

echo "cleaning up output file"
> "$OUTPUT_FILE"
rm "$OUTPUT_FILE"

echo "all tests passed successfully"
