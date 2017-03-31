#!/bin/bash

ES_ALIVE=/tmp/elastic_search_is_alive

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

# Start kibana
CMD="kibana"
CMDARGS="$@"
exec "$CMD" $CMDARGS
