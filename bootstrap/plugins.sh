#!/usr/bin/env bash

set -o errexit
set -o nounset

if [ $# -ne 1 ]; then
  echo "Usage: $0 <start|stop|status>"
  exit 0
fi

HERE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$HERE"

export PATH=$GOPATH/src/github.com/docker/infrakit/build:$PATH

INFRAKIT_HOME=${INFRAKIT_HOME:-~/.infrakit}
export INFRAKIT_HOME
logs=$INFRAKIT_HOME/logs
mkdir -p $logs
# set the leader -- for os / file based leader detection for manager
leaderfile=$INFRAKIT_HOME/leader
# set the inventory file -- for keeping trace of instances info in a single file
INVENTORY_FILE=$INFRAKIT_HOME/inventory

case $1 in
start)
  configstore=$INFRAKIT_HOME/configs
  mkdir -p $configstore
  rm -rf $configstore/*
  echo group > $leaderfile
  infrakit plugin start --config-url file:///$PWD/plugins.json --exec os \
	 manager \
	 group-stateless \
	 flavor-combo \
	 flavor-swarm \
	 flavor-vanilla \
	 instance-vagrant \
	 instance-terraform &
  sleep 3
  echo "Plugins started."
  echo "Do something like: infrakit manager commit file://$PWD/amp.json"
  ;;

stop)
  infrakit plugin ls -q | awk '{print $1}' | xargs infrakit plugin stop || true
  killall infrakit infrakit-manager infrakit-group-default infrakit-instance-vagrant infrakit-instance-terraform infrakit-flavor-combo infrakit-flavor-vanilla infrakit-flavor-swarm || true
  ;;

status)
  tfile=$(mktemp)
  code=0
  infrakit plugin ls | tee $tfile
  for p in manager group-stateless flavor-combo flavor-vanilla instance-vagrant instance-terraform; do
    grep -q $p $tfile || code=$((code+$?))
  done
  rm "$tfile"
  exit $code
  ;;

*)
  exit 1
  ;;
esac
