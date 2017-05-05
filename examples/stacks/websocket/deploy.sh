#!/bin/sh
set -o errexit
amp -s localhost stack rm websocket || true
docker build -t 127.0.0.1:5000/examples/ws-bash-server server
docker build -t 127.0.0.1:5000/examples/ws-bash-web web
docker push 127.0.0.1:5000/examples/ws-bash-server
docker push 127.0.0.1:5000/examples/ws-bash-web
amp -s localhost stack deploy -c stack.yml websocket
