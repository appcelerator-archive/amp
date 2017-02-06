#!/bin/sh
amp stack rm websocket || true
docker build -t examples/ws-bash-server server
amp registry push examples/ws-bash-server
docker build -t examples/ws-bash-web web
amp registry push examples/ws-bash-web
amp stack up websocket -c stack.yml
