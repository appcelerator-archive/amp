#!/bin/sh
set -o errexit
amp -s localhost stack rm micro || true
docker build -t 127.0.0.1:5000/examples/micro .
docker push 127.0.0.1:5000/examples/micro
amp -s localhost stack deploy -c stack.yml micro

