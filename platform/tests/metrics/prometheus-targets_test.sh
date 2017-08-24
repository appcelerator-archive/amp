#!/bin/bash
tmpfile=$(mktemp)
code=1
SECONDS=0
TIMEOUT=25
while [[ $code -ne 0 ]]; do
  [[ $SECONDS -gt $TIMEOUT ]] && break
  code=0
  docker run --rm --network ampnet appcelerator/alpine:3.6.0 curl -sf http://prometheus:9090/targets > $tmpfile
  countstate=$(grep -wc "State" $tmpfile)
  countup=$(grep -wc "up" $tmpfile)
  [[ $countstate -eq $countup && $countup -gt 0 ]] && break
  code=1
done
[[ $code -ne 0 ]] && (echo "issue with prometheus targets" ; cat $tmpfile)
rm $tmpfile
exit $code
