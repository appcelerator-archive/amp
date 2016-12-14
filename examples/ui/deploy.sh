#!/bin/sh
amp stack rm -f gw-ui || true
docker build -t examples/gw-ui .
amp registry push examples/gw-ui
amp stack up -f stack.yml gw-ui
