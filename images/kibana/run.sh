#!/bin/bash

ES_ALIVE=/tmp/elastic_search_is_alive

# check if a certificate is present
if [ -n "$SERVER_SSL_CERTIFICATE" && -n "$SERVER_SSL_KEY" && -f "$SERVER_SSL_CERTIFICATE" && -f "$SERVER_SSL_KEY" ]; then
  echo "found $SERVER_SSL_CERTIFICATE and $SERVER_SSL_KEY"
  echo "enabling SSL"
  SERVER_SSL_ENABLED=true
else
  echo "disabling SSL"
  SERVER_SSL_ENABLED=false
fi

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


rm -f "$ES_ALIVE"
# Start pre-configuration in case ES is already available
# In background because it will need a connection from Kibana to ES
/preconfiguration.sh &

# wait a bit for elasticsearch
for w in $(seq 16); do
  if [[ -f "$ES_ALIVE" ]]; then
    break
  fi
  sleep 1
done

# Start kibana
CMD="kibana"
CMDARGS="$@"
exec "$CMD" $CMDARGS
