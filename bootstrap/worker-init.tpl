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
  mkdir -p /etc/systemd/system/docker.service.d
  cat > /etc/systemd/system/docker.service.d/docker.conf <<EOF
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H 0.0.0.0:{{ if ref "/certificate/ca/service" }}{{ ref "/docker/remoteapi/tlsport" }} --tlsverify --tlscacert={{ ref "/docker/remoteapi/cafile" }} --tlscert={{ ref "/docker/remoteapi/srvcertfile" }} --tlskey={{ ref "/docker/remoteapi/srvkeyfile" }}{{else }}{{ ref "/docker/remoteapi/port" }}{{ end }} -H unix:///var/run/docker.sock{{ if ref "/bootstrap/ip" }} --registry-mirror=http://{{ ref "/bootstrap/ip" }}:5000 --insecure-registry=http://{{ ref "/bootstrap/ip" }}:5000{{ end }}
EOF
  systemctl daemon-reload
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
