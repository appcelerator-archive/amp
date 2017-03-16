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
  systemctl stop docker.service
fi
# Use an EBS volume for the devicemapper
if [ "x$provider" = "xaws" ]; then
  rm -rf /var/lib/docker
  _attach_ebs_volume /dev/sdn /var/lib/docker "Docker AUFS" {{ ref "/docker/aufs/size" }}
fi

mkdir -p /etc/docker
cat << EOF > /etc/docker/daemon.json
{
  "labels": {{ INFRAKIT_LABELS | to_json }}
}
EOF

{{ if ref "/certificate/ca/service" }}{{ include "request-certificate.sh" }}{{ end }}

if [ "x$provider" != "xdocker" ]; then
  systemctl start docker.service
  sleep 2
fi

docker swarm join --token {{  SWARM_JOIN_TOKENS.Worker }} {{ SWARM_MANAGER_ADDR }}
