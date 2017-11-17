#!/bin/bash

# Add kibana as command if needed
if [ "${1:0:1}" = '-' ]; then
    echo "INFO - adding arguments $@ to kibana"
    set -- kibana "$@"
fi
program="$1"

# check if a certificate is present
if [[ -n "$SERVER_SSL_CERTIFICATE" && -n "$SERVER_SSL_KEY" && -f "$SERVER_SSL_CERTIFICATE" && -f "$SERVER_SSL_KEY" ]]; then
  echo "found $SERVER_SSL_CERTIFICATE and $SERVER_SSL_KEY"
  echo "enabling SSL"
  SERVER_SSL_ENABLED=true
else
  echo "disabling SSL"
  SERVER_SSL_ENABLED=false
fi
export SERVER_SSL_ENABLED

# Configuration file for kibana.yml
if [[ -f "/opt/kibana/config/kibana.yml.tpl" ]]; then
    echo "Kibana configuration file will be generated."
    envtpl -o "/opt/kibana/config/kibana.yml" "/opt/kibana/config/kibana.yml.tpl" && rm "/opt/kibana/config/kibana.yml.tpl"
else
    if [[ -f "/opt/kibana/config/kibana.yml" ]]; then
        echo "Kibana configuration file already generated."
    else
        echo "Kibana configuration file and template are missing."
        exit 1
    fi
fi

chown elastico /opt/kibana/config/kibana.yml

# Adding index to Kibana to set as defaultIndex
SECONDS=0
maxtime=60
rc=1
while [[ $rc -ne 0 ]]; do
  echo "Attempt to add index to Kibana ($((maxtime - SECONDS)) sec left)"
  curl -sf -XPUT $ELASTICSEARCH_URL/.kibana/index-pattern/ampbeat-* -d '{"title" : "ampbeat-*",  "timeFieldName": "@timestamp"}' &> /dev/null
  rc=$?
  if [ $SECONDS -gt $maxtime ]; then
    echo "[$0] FATAL - Can't reach Elasticsearch"
    exit 1
  fi
  sleep 1
done

# Update Kibana config to set index as default index
curl -sf -XPUT $ELASTICSEARCH_URL/.kibana/config/$KIBANA_VERSION -d '{"defaultIndex" : "ampbeat-*"}' &> /dev/null
if [[ $? -ne 0 ]]; then
  echo "Failed to configure Kibana default index"
  exit 1
fi
echo "Successfully configured Kibana default index"

# Drop root privileges if we are running kibana
# allow the container to be started with `--user`
[[ "x$program" = "xkibana" && "$(id -u)" = '0' ]] && set -- gosu elastico "$@"

# Start kibana
exec "$@"
