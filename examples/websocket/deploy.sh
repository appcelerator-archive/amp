#!/bin/sh
AMPOPTS="--server localhost:8080"
amp $AMPOPTS stack rm -f websocket || true
docker build -t examples/ws-bash-server server
amp $AMPOPTS registry push examples/ws-bash-server
docker build -t examples/ws-bash-web web
amp $AMPOPTS registry push examples/ws-bash-web
amp $AMPOPTS stack up -f stack.yml websocket
