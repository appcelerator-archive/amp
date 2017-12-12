#!/bin/bash

ES_CONF=/opt/elasticsearch/config/elasticsearch.yml

# Add elasticsearch as command if needed
if [ "${1:0:1}" = '-' ]; then
    echo "INFO - adding arguments $@ to elasticsearch"
    set -- elasticsearch "$@"
fi
program="$1"

if [[ -z "$ES_JAVA_OPTS" && -z "$JAVA_HEAP_SIZE" ]]; then
  # adjust max heap size to available memory
  if [[ -f /proc/meminfo ]]; then
    tmem=$(grep MemTotal /proc/meminfo | awk '{print int($2 * 0.001)}')
    echo "INFO - system memory is ${tmem}M"
    if [[ $tmem -lt 1024 ]]; then
        echo "INFO - set java heap size to floor value"
        JAVA_HEAP_SIZE=512
    elif [[ $tmem -lt 4096 ]]; then
        echo "INFO - set java heap size to ramping value"
        # 512 to 2048
        JAVA_HEAP_SIZE=$((512 + (tmem - 1024) / 2))
    else
        echo "INFO - set java heap size to 50%"
        JAVA_HEAP_SIZE=$((tmem / 2))
    fi
    export JAVA_HEAP_SIZE
    echo "Java Heap Size: $JAVA_HEAP_SIZE MB"
  else
    echo "WARN - can't read /proc/meminfo, using default java heap size value"
  fi
fi
if [[ -n "$JAVA_HEAP_SIZE" ]]; then
  ES_JAVA_OPTS="$ES_JAVA_OPTS -Xms${JAVA_HEAP_SIZE}M -Xmx${JAVA_HEAP_SIZE}M ${java_max_direct_mem_size:+"-XX:MaxDirectMemorySize=$java_max_direct_mem_size"} $java_opts"
fi
ES_JAVA_OPTS="${ES_JAVA_OPTS} -Des.cgroups.hierarchy.override=/ -Djava.security.policy=file:///opt/elasticsearch/config/java.policy"
export ES_JAVA_OPTS

echo "memory lock limit:"
ulimit -Hl
ulimit -l
echo -n "Hard limit max open fd:"
ulimit -Hn
echo -n "Max open fd:"
ulimit -n

SECONDS=0
echo "resolving the container IP with Docker DNS..."
while [ -z "$cip" ]; do
  cip=$(dig +short $(hostname))
  # checking that the returned IP is really an IP
  echo "$cip" | egrep -qe "^[0-9\.]+$"
  if [ -z "$cip" ]; then
    sleep 1
  fi
  [[ $SECONDS -gt 10 ]] && break
done
if [[ -z "$cip" ]]; then
  cip=$(grep $(hostname) /etc/hosts |awk '{print $1}' | head -1)
fi
if [[ -z "$cip" ]]; then
  echo "unable to get this container's ip"
  exit 1
fi
echo "this node ip is $cip"
export PUBLISH_HOST=$cip

# if the unicast hosts is the list of tasks from a swarm service, we need to substract this node IP
echo "$UNICAST_HOSTS" | grep -q "^tasks."
if [[ $? -eq 0 ]]; then
  # wait for master nodes to be available
  echo "waiting for other tasks to be available..."
  SECONDS=0
  typeset -i count=0
  while [[ $count -lt $MIN_MASTER_NODES ]]; do
    if [[ $SECONDS -gt 15 ]]; then
      echo "Expecting $MIN_MASTER_NODES tasks, only found $count after $SECONDS sec, abort"
      exit 1
    fi
    sleep 1
    tips=$(dig +short $UNICAST_HOSTS | grep -v "$cip")
    count=$(echo $tips | wc -w)
  done
  UNICAST_HOSTS=$(echo $tips | tr ' ' ',')
  echo "$count tasks found for unicast zen discovery ($UNICAST_HOSTS)"
fi

if [[ -f $ES_CONF.tpl ]]; then
    mv $ES_CONF $ES_CONF.bak
    envtpl -o $ES_CONF $ES_CONF.tpl && rm $ES_CONF.tpl
    if [[ $? -ne 0 ]]; then
        echo "WARNING - configuration file update failed"
        echo "WARNING - using default configuration instead"
        mv $ES_CONF.bak $ES_CONF
    fi
fi

runes=0
# Drop root privileges if we are running elasticsearch
# allow the container to be started with `--user`
if [[ "x$program" = "xelasticsearch" ]]; then runes=1; fi

if [[ $runes -eq 1 && "$(id -u)" = '0' ]]; then
    echo "INFO - setting user perms on elasticsearch data dir"
    # Change the ownership of /opt/elasticsearch/data to elasticsearch
    chown -R elastico:elastico /opt/elasticsearch/data

    echo "INFO - running $1 as user elastico"
    set -- gosu elastico "$@"
fi

echo "INFO - running $@"
exec "$@"
