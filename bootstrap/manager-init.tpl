#!/bin/bash
set -o errexit
set -o nounset
set -o xtrace

{{ source "default.ikt" }}
{{ source "file:///infrakit/env.ikt" }}
{{ source "provider.sh" }}
{{ source "install-docker.sh" }}
{{ source "attach-ebs-volume.sh" }}

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

if [ "x$provider" != "xdocker" ]; then
  systemctl start docker.service
  sleep 2
fi

{{ if ref "/certificate/ca/service" }}{{ include "request-certificate.sh" }}{{ end }}

# INSTANCE_LOGICAL_ID can be an IP of a hostname, we need an IP
IP="{{ INSTANCE_LOGICAL_ID }}"
if ! $(echo "$IP" | egrep -q "([0-9.]+){4}"); then
  IP=$(nslookup {{ INSTANCE_LOGICAL_ID }}  2>/dev/null | awk '$1 == "Address" {print $3}' | tail -1)
fi

{{ if and ( eq INSTANCE_LOGICAL_ID SPEC.SwarmJoinIP ) (not SWARM_INITIALIZED) }}
if [ "x$provider" != "xdocker" ]; then
  mkdir -p /etc/systemd/system/docker.service.d
  cat > /etc/systemd/system/docker.service.d/docker.conf <<EOF
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// -H 0.0.0.0:{{ if ref "/certificate/ca/service" }}{{ ref "/docker/remoteapi/tlsport" }} --tlsverify --tlscacert={{ ref "/docker/remoteapi/cafile" }} --tlscert={{ ref "/docker/remoteapi/srvcertfile" }} --tlskey={{ ref "/docker/remoteapi/srvkeyfile" }}{{else }}{{ ref "/docker/remoteapi/port" }}{{ end }} -H unix:///var/run/docker.sock
EOF

  # Restart Docker to let port listening take effect.
  systemctl daemon-reload
  systemctl restart docker.service
fi

  {{/* The first node of the special allocations will initialize the swarm. */}}
  docker swarm init --advertise-addr $IP

{{ else }}

  {{/* The rest of the nodes will join as followers in the manager group. */}}
  docker swarm join --token {{ SWARM_JOIN_TOKENS.Manager }} {{ SWARM_MANAGER_ADDR }}

{{ end }}
