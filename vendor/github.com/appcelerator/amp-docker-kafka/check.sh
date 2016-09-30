#!/bin/bash

set -euo pipefail

#TOPIC_LIST="amp-logs amp-service-start amp-service-stop amp-service-terminate amp-docker-events amp-service-events"

if [ -f /tmp/kafka-topics ]; then
    pidof java
else
    echo "Checking topic list..."
    list=$(bin/kafka-topics.sh --zookeeper $ZOOKEEPER_CONNECT --list)
    echo "Checking topic list : $list"
    for topic in $TOPIC_LIST
    do
        if [[ $list =~ ^.*$topic ]]; then
            echo "$topic exists"
        else
            echo "Creating topic $topic..."
            bin/kafka-topics.sh --zookeeper $ZOOKEEPER_CONNECT --create --partitions=1 --replication-factor=1 --topic $topic
            echo "Creating topic $topic Done"
        fi
    done
    touch /tmp/kafka-topics
fi

