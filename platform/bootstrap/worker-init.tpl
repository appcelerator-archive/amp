#!/bin/sh
set -o errexit
set -o nounset
set -o xtrace

{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ source "provider.sh" }}
{{ source "attach-ebs-volume.sh" }}
{{ source "install-docker.sh" }}

# if we're already in a Docker container, either it's DinD and Docker is already installed
# or it's a non DinD Docker container and we can't easily install Docker
if [ "x$provider" != "xdocker" ]; then
  _install_docker
  systemctl stop docker.service || service docker stop
fi
# Use an EBS volume for the devicemapper
if [ "x$provider" = "xaws" ]; then
  rm -rf /var/lib/docker
  _attach_ebs_volume /dev/sdn /var/lib/docker "Docker AUFS" {{ var "/docker/aufs/size" }}
fi

# if mirrorregistries is defined, it's a comma separated list (or just a list) of registries, it has to be correctly quoted to be inserted in the json file
{{ if var "/docker/mirrorregistries" }}
mirrorregistries="$(echo {{ var "/docker/mirrorregistries" }} | tr ',' ' ' | sed 's/  */", "/g')"
{{ end }}


mkdir -p /etc/docker
cat << EOF > /etc/docker/daemon.json
{
  "labels": {{ INFRAKIT_LABELS | jsonEncode }},
  "experimental": true,
  "metrics-addr": "0.0.0.0:9323",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }{{ if var "/docker/mirrorregistries" }},
  "registry-mirrors": [ "$mirrorregistries" ]{{ end }}
}
EOF

if [ "x$provider" != "xdocker" ]; then
  systemctl start docker.service || service docker start
  sleep 2
else
  # TODO: send kill -HUP to reload the labels, see appcelerator/amp#1123
  kill -s HUP $(cat /var/run/docker.pid)
fi
# TODO: send kill -HUP to reload the labels, see appcelerator/amp#1123

docker swarm join --token {{  SWARM_JOIN_TOKENS.Worker }} {{ SWARM_MANAGER_ADDR }}

# InfraKit sets labels on the engine, we want them on the node
nodeid=$(docker info 2>/dev/null| grep NodeID | awk '{print $2}')
labels="$(echo {{ INFRAKIT_LABELS }} | tr -d '[]')"
remote_api="$(echo {{ SWARM_MANAGER_ADDR }} | cut -d: -f1)"
for label in $labels; do
  docker -H $remote_api node update --label-add "$label" "$nodeid"
done
