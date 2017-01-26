#!/bin/bash

ELASTICSEARCH_URL=http://elasticsearch-amp:9200

echo -n "test elasticsearch availability... "
maxretries=15
sleeptime=1
while [[ $try -lt $maxretries ]]; do
    version=$(curl -s "$ELASTICSEARCH_URL" | jq .version.number 2>/dev/null)
    if [[ $? -eq 0 && -n "$version" ]]; then
        break
    fi
    (( try++ ))
    echo -n "+"
    sleep $sleeptime
done
if [[ $try -ge $maxretries ]]; then
    echo
    echo "failed"
    curl "$ELASTICSEARCH_URL"
    exit 1
fi
echo " ($version) [OK]"

echo "all tests passed successfully"
