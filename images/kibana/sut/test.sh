#!/bin/bash

KIBANA_URL=http://kibana:5601

echo -n "test kibana status... "
maxretries=30
sleeptime=1
try=0
while [[ $try -lt $maxretries ]]; do
    status=$(curl -sf -m 2 "$KIBANA_URL/api/status" 2>/dev/null | jq -r '.status.overall.state')

    if [[ $? -eq 0 && "x$status" = "xgreen" ]]; then
        break
    fi
    (( try++ ))
    echo -n "+"
    sleep $sleeptime
done
if [[ $try -ge $maxretries ]]; then
    echo
    echo "failed"
    curl "$KIBANA_URL/api/status"
    echo
    echo "ssl endpoint connection test:"
    curl -k https://kibana
    exit 1
fi
echo " ($status) [OK]"

echo "all tests passed successfully"
