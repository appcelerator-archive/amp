#!/bin/sh
amp stack rm -f websocket || true
docker build -t examples/ws-bash-server server
amp registry push examples/ws-bash-server
docker build -t examples/ws-bash-web web
amp registry push examples/ws-bash-web
amp stack up -f stack.yml websocket
