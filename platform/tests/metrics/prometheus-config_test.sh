#!/bin/bash
tmpfile=$(mktemp)
code=0
docker run --rm --network ampnet appcelerator/alpine:3.6.0 curl -sf http://prometheus:9090/config > $tmpfile
grep -q -- "- job_name: nodes" $tmpfile
code=$((code+$?))
grep -q -- "- job_name: docker-engine" $tmpfile
code=$((code+$?))
grep -q -- "taskname: amp_amplifier." $tmpfile
code=$((code+$?))
grep -q -- "- amp_haproxy_exporter:9101" $tmpfile
code=$((code+$?))
[[ $code -ne 0 ]] && (echo "$code issues found in prometheus configuration file" ; cat $tmpfile)
rm $tmpfile
exit $code
