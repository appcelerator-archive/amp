#!/bin/bash

ES_URL=${ELASTICSEARCH_URL:-http://elasticsearch:9200}
ES_ALIVE=/tmp/elastic_search_is_alive
maxretries=10
sleeptime=3
# Wait for elasticsearch availability
try=0

echo "INFO: waiting for elasticsearch availability ($ES_URL)..."
while [[ $try -lt $maxretries ]]; do
    version=$(curl -s "$ES_URL" | jq .version.number 2>/dev/null)
    if [[ $? -eq 0 && -n "$version" ]]; then
        echo
        echo "[$0] ElasticSearch $version is available after $try x $sleeptime sec"
        break
    fi
    (( try++ ))
    sleep $sleeptime
done
echo
if [[ $try -ge $maxretries ]]; then
    echo "[$0] WARN - elasticsearch is not available"
    curl -s "$ES_URL"
    exit 0
else
    touch $ES_ALIVE
fi
