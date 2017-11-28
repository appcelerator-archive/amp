#!/bin/bash

# Add kibana as command if needed
if [ "${1:0:1}" = '-' ]; then
    echo "[$0] INFO - adding arguments $@ to kibana"
    set -- kibana "$@"
fi
program="$1"
# fix cgroup membership, see https://github.com/elastic/kibana-docker/blob/master/build/kibana/bin/kibana-docker
echo $@ | grep -vq cgroup && set -- kibana --cpu.cgroup.path.override=/ --cpuacct.cgroup.path.override=/ "${@:2}"
KIBANA_CONFIG=/opt/kibana/config/kibana.yml

# Drop root privileges if we are running kibana
# allow the container to be started with `--user`
[[ "x$program" = "xkibana" && "$(id -u)" = '0' ]] && set -- gosu elastico "$@"

# check if a certificate is present
if [[ -n "$SERVER_SSL_CERTIFICATE" && -n "$SERVER_SSL_KEY" && -f "$SERVER_SSL_CERTIFICATE" && -f "$SERVER_SSL_KEY" ]]; then
  echo "[$0] INFO - found $SERVER_SSL_CERTIFICATE and $SERVER_SSL_KEY"
  echo "enabling SSL"
  SERVER_SSL_ENABLED=true
else
  echo "[$0] INFO - disabling SSL"
  SERVER_SSL_ENABLED=false
fi
export SERVER_SSL_ENABLED

# Configuration file for kibana.yml
if [[ -f "${KIBANA_CONFIG}.tpl" ]]; then
  echo "[$0] INFO -  Kibana configuration file generation from template."
  envtpl -o "$KIBANA_CONFIG" "${KIBANA_CONFIG}.tpl" && rm "${KIBANA_CONFIG}.tpl"
else
  if [[ -f "$KIBANA_CONFIG" ]]; then
    echo "[$0] INFO - Kibana configuration file already generated."
  else
    echo "[$0] FATAL - Kibana configuration file and template are missing."
    exit 1
  fi
fi

chown elastico /opt/kibana/config/kibana.yml

echo "[$0] INFO - Starting Kibana for pre configuration"
mv $KIBANA_CONFIG $KIBANA_CONFIG.disabled
kibana --server.host=localhost --server.port=6501 --elasticsearch.url="${ELASTICSEARCH_URL}" --elasticsearch.startupTimeout=900 --logging.silent=true &
kibanapid=$!
mv $KIBANA_CONFIG.disabled $KIBANA_CONFIG
# Create kibana index pattern
url="http://localhost:6501"
index_pattern="ampbeat-*"
time_field="@timestamp"
SECONDS=0
maxtime=40
id=""
while true; do
  if [[ $SECONDS -ge $maxtime ]]; then
    echo "[$0] FATAL - failed to create the index pattern"
    exit 1
  fi
  if [[ -z "$id" ]]; then
    echo "[$0] INFO - Attempting to create the Kibana index pattern ($((maxtime - SECONDS)) sec left)"
    id=$(curl -f -m 2 -XPOST -H "Content-Type: application/json" -H "kbn-xsrf: amp" "$url/api/saved_objects/index-pattern" \
      -d "{\"attributes\":{\"title\":\"$index_pattern\",\"timeFieldName\":\"$time_field\"}}" 2>/dev/null | jq -r '.id')
    [[ -n "$id" ]] && echo "[$0] INFO - Index pattern successfully created: $id"
  else
    # Create Kibana default index
    echo "[$0] INFO - Create default index ($((maxtime - SECONDS)) sec left)"
    curl -sf -m 2 -XPOST -H "Content-Type: application/json" -H "kbn-xsrf: amp" "$url/api/kibana/settings/defaultIndex" \
      -d "{\"value\":\"$id\"}"
    if [[ $? -ne 0 ]]; then
      echo "[$0] INFO - Failed to set the default index, will try again"
    else
      # we're done: index pattern + default index have been configured
      break
    fi
  fi
  sleep 1
done
echo
echo "[$0] INFO - Successfully configured Kibana default index"

# Stop the pre configuration kibana
echo "[$0] INFO - Stopping the pre configuration kibana"
kill $kibanapid
if [[ $? -ne 0 ]]; then
  kill -9 $kibanapid
fi
if [[ $? -ne 0 ]]; then
  echo "[$0] FATAL - Failed to stop the pre configuration Kibana"
fi

echo "[$0] INFO - Replacing this process with kibana"
exec "$@"
