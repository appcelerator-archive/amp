#!/bin/bash

ES_URL=${ELASTICSEARCH_URL:-http://elasticsearch:9200}
ES_ALIVE=/tmp/elastic_search_is_alive
maxretries=60
# Wait for elasticsearch availability
try=0

SECONDS=0
echo "INFO: waiting for elasticsearch availability ($ES_URL)..."
while [[ $try -lt $maxretries ]]; do
    version=$(curl -s "$ES_URL" | jq .version.number 2>/dev/null)
    if [[ $? -eq 0 && -n "$version" ]]; then
        echo
        echo "[$0] ElasticSearch $version is available (after $SECONDS sec)"
        break
    fi
    sleep 1
done
echo
if [[ $SECONDS -ge $maxretries ]]; then
    echo "[$0] WARN - elasticsearch is not available, abort the index and dashboard configuration"
    curl -s "$ES_URL"
    exit 0
else
    touch $ES_ALIVE
fi

# Add index to Kibana which you want to set as defaultIndex
echo "[$0]__1. Setting Kibana index __________"
SECONDS=0
rc=1
while [ $rc != 0 ]; do
  curl -sf -XHEAD "$ES_URL/.kibana"
  rc=$?
  sleep 1
  if [ $SECONDS -gt $maxretries ]; then
    echo "[$0] WARN - Can't find .kibana index"
    exit 1
  fi
done
for ipd in $(ls /var/lib/kibana/index-patterns/*.json); do
  name=$(basename $ipd .json)
  echo "[$0]__   index pattern $name __________"
  curl -sf -XPUT "$ES_URL/.kibana/index-pattern/${name}-*" -d @${ipd} > /tmp/debug.log
  if [[ $? -ne 0 ]]; then
        echo "[$0] WARN: FAILED to configure Kibana index $name"
        echo "[$0] $(cat /tmp/debug.log)"
        exit 1
  else
        echo "[$0] Ok"
        rm /tmp/debug.log
  fi
done

# Check the Kibana config Id
echo "[$0]__ 2. Looking for .kibana/config ID __________"
sleep 2
retry=1
while [[ $retry -lt $maxretries ]]; do
    curl -s -XGET $ES_URL/.kibana/config/_search?stored_fields= || true
    configVersion=$(curl -s -XGET $ES_URL/.kibana/config/_search?stored_fields= | jq .hits.hits[0]._id | sed 's/"//g') || true
    if [[ "$configVersion" == "null" || "x$configVersion" == 'x' ]]; then
        echo "[$0] WARN: FAILED to fetch Kibana config version (retry #$retry returns $configVersion)"
        sleep 1
        ((retry++))
    else
        echo "[$0] OK. Value found: .kibana/config = $configVersion"
        break
    fi
done

echo "[$0]__ 3. Retrieve fields of .kibana/config/$configVersion __________"
configValues=$(curl -s -XGET $ES_URL/.kibana/config/$configVersion | jq "._source") || true
if [[ "$configValues" == "null" || "x$configValues" == 'x' ]]; then
    echo "[$0] WARN: FAILED to fetch Kibana config value"
else
    echo "[$0] OK. Value found: .kibana/config/$configVersion = $configValues"
fi
# Add index key to JSON
configValues=$(echo "$configValues" '{ "defaultIndex" : "ampbeat-*" }' | jq -s add)
if [[ "$configValues" == "null" || "x$configValues" == 'x' ]]; then
    configValues="null"
    echo "[$0] ERROR: FAILED to merge Kibana config value"
else
    echo "[$0] Ok. Merged."
fi

if [[ "$configValues" == "null" ]]; then
  exit 1
fi

# Change your Kibana config to set index added above as defaultIndex
echo "[$0]__ 4. Setting index as default __________"
echo "[$0] Will update .kibana/config/$configVersion with value: $(echo $configValues | sed 's/\r\n/ /g')"
curl -s -XPUT "$ES_URL/.kibana/config/$configVersion" -d "$configValues" > /dev/null
if [[ $? -ne 0 ]]; then
    echo "[$0] ERROR: FAILED to set Kibana index as default"
else
    echo "[$0] Ok."
fi

# Workaround for: https://github.com/elastic/beats-dashboards/issues/94
curl -XPUT "$ES_URL/.kibana/_mapping/search" -d'{"search": {"properties": {"hits": {"type": "integer"}, "version": {"type": "integer"}}}}'

echo "[$0]__ 5. Import saved objects __________"
for f in $(ls /var/lib/kibana/saved-objects/search_*.json); do
  name=$(basename $f .json)
  t=$(echo $name | cut -d_ -f1)
  echo "[$0]__  import $name __________"
  curl -s -XPUT "$ES_URL/.kibana/$t/$name" -d @${f} > /tmp/debug.log || true
  if [[ "x$(cat /tmp/debug.log | jq '.created')" != "xtrue" ]]; then
        echo "[$0] WARN: FAILED to import $t $name"
        echo "[$0] $(cat /tmp/debug.log)"
        exit 1
  else
        echo "[$0] Ok"
        rm /tmp/debug.log
  fi
done
for f in $(ls /var/lib/kibana/saved-objects/visualization_*.json); do
  name=$(basename $f .json)
  t=$(echo $name | cut -d_ -f1)
  echo "[$0]__  import $name __________"
  curl -s -XPUT "$ES_URL/.kibana/$t/$name" -d @${f} > /tmp/debug.log || true
  if [[ "x$(cat /tmp/debug.log | jq '.created')" != "xtrue" ]]; then
        echo "[$0] WARN: FAILED to import $t $name"
        echo "[$0] $(cat /tmp/debug.log)"
        exit 1
  else
        echo "[$0] Ok"
        rm /tmp/debug.log
  fi
done
for f in $(ls /var/lib/kibana/saved-objects/dashboard_*.json); do
  name=$(basename $f .json)
  t=$(echo $name | cut -d_ -f1)
  echo "[$0]__  import $name __________"
  curl -s -XPUT "$ES_URL/.kibana/$t/$name" -d @${f} > /tmp/debug.log || true
  if [[ "x$(cat /tmp/debug.log | jq '.created')" != "xtrue" ]]; then
        echo "[$0] WARN: FAILED to import $t $name"
        echo "[$0] $(cat /tmp/debug.log)"
        exit 1
  else
        echo "[$0] Ok"
        rm /tmp/debug.log
  fi
done
