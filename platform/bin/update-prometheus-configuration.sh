#!/bin/bash

METRICS_PORT=9323
PROMETHEUS_VOLUME=amp_prometheus
PROMETHEUS_DIR=/etc/prometheus
PROMETHEUS_FILE=prometheus.yml
PROMETHEUS_PORT=9090
ALERTMANAGER_DIR=/etc/alertmanager
ALERTMANAGER_FILE=config.yml
D4MIP="192.168.65.1"
DOCKER_METRICS_PORT=9323
TELEGRAF_METRICS_PORT=9126
DRY_RUN=0

manager_list(){
  local _remote_managers
  local _managers
  local _m
  docker node ls >/dev/null 2>&1
  if [[ $? -eq 0 ]]; then
    echo "unix:///var/run/docker.sock"
    return 0
  fi
  _remote_managers=$(docker info -f '{{range .Swarm.RemoteManagers}} {{.Addr}} {{end}}')
  for _m in $_remote_managers; do
    _managers="$_managers ${_m%:*}"
  done
  echo $_managers
}

get_telegraf_remotes(){
  local _telegraf_service=telegraf
  local _remotes
  _remotes=$(dig +short tasks.$_telegraf_service)
  [[ -n "$_remotes" ]] && echo $_remotes

}

prepare_prometheus_conf(){
  local _remotes=$*
  local _remote
  local _docker_remotes
  local _telegraf_remotes

  if [[ $# -eq 1 && "$_remotes" = "127.0.0.1" ]]; then
    # Docker for Mac/Windows: loopback address won't work, fix it
    _remotes=$(ifconfig $(netstat -nr | awk 'NF==6 && $1 ~/default/ {print $6}' | tail -1) | awk '$1 == "inet" {print $2}' | grep -v "127.0.0.1" | sed -e 's/addr://')
  fi
  if [[ $# -eq 1 && $(uname) != "Linux" ]]; then
    # special case: Docker for Mac/Windows expose metrics on the VM IP
    _docker_remotes="${D4MIP}"
  else
    _docker_remotes=$_remotes
  fi

  cat > $prometheus_conf << EOF
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
      monitor: 'amp-monitor'

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  - "/run/secrets/*.rules"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
#  - job_name: 'prometheus'
#    static_configs:
#      - targets: ['localhost:9090']
  - job_name: 'docker-engine'
    static_configs:
      - targets:
EOF
  for _remote in $_docker_remotes; do
    echo "        - '${_remote}:$DOCKER_METRICS_PORT'" >> $prometheus_conf
  done
  cat >> $prometheus_conf << EOF
  - job_name: 'system'
    static_configs:
      - targets:
EOF
  for _remote in $_remotes; do
    echo "        - '${_remote}:$TELEGRAF_METRICS_PORT'" >> $prometheus_conf
  done
}

prepare_alertmanager_conf(){
  cat > $alertmanager_conf << EOF
global:
  slack_api_url:
route:
  receiver: 'default-receiver'
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  group_by: [cluster, alertname]
EOF
}

update_prometheus_container(){
  local _container
  _container=$(docker ps --format '{{.Names}}' | grep prometheus) || return 1
  [[ $(wc -w <<< "$_container") -ne 1 ]] && return 1

  if [[ $DRY_RUN -eq 1 ]]; then
    echo "Dry run, container $_container would have been updated with this configuration:" >&2
    cat "$prometheus_conf"
    return 0
  fi
  echo "copying prometheus configuration in $_container" >&2
  docker cp $prometheus_conf $_container:$PROMETHEUS_DIR/$PROMETHEUS_FILE || return 1

  echo "reloading prometheus" >&2
  #curl -XPOST http://localhost:$PROMETHEUS_PORT/-/reload
  docker exec "$_container" kill -s HUP 1
}
update_alertmanager_container(){
  local _container
  _container=$(docker ps --format '{{.Names}}' | grep alertmanager) || return 1
  [[ $(wc -w <<< "$_container") -ne 1 ]] && return 1

  echo "copying alertmanager configuration in $_container" >&2
  docker cp $alertmanager_conf $_container:$ALERTMANAGER_DIR/$ALERTMANAGER_FILE || return 1

  echo "reloading alertmanager" >&2
  docker exec "$_container" kill -s HUP 1
}

cleanup(){
  [[ -n "$prometheus_conf" && -f "$prometheus_conf" ]] && rm -f "$prometheus_conf"
  [[ -n "$alertmanager_conf" && -f "$alertmanager_conf" ]] && rm -f "$alertmanager_conf"
}

prometheus_conf=$(mktemp)
alertmanager_conf=$(mktemp)
trap cleanup EXIT

while getopts ":nh" opt; do
  case $opt in
  n)
    DRY_RUN=1
    ;;
  h)
    echo "Usage:" >&2
    echo "$0 [-n] [-h]" >&2
    ;;
  esac
done
shift "$((OPTIND-1))"

managers=$(manager_list)
for manager in $managers; do
  nodes=$(docker -H $manager node ls -q) && break
done
[[ -z "$nodes" ]] && exit 1

for node in $nodes; do
  for manager in $managers; do
    remote=$(docker -H $manager node inspect $node -f '{{.Status.Addr}}')
    if [[ $? -eq 0 && -n "$remote" && "$remote" != "0.0.0.0" ]]; then
      break
    elif [[ -n "$remote" && "$remote" = "0.0.0.0" ]]; then
      # try the manager address instead
      remote=$(docker -H $manager node inspect $node -f '{{.ManagerStatus.Addr}}' | cut -d: -f1) 
      if [[ ${PIPESTATUS[0]} -ne 0 || "x$remote" = "x0.0.0.0" ]]; then
        echo "Failed to get IP of node $node, abort" >&2
        exit 1
      fi
      [[ -n "$remote" ]] && break
    fi
  done
  remotes="$remotes ${remote}"
done
[[ -z "$remotes" ]] && exit 1

prepare_prometheus_conf $remotes
update_prometheus_container
