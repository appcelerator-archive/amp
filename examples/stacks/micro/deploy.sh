#!/bin/sh
amp stack rm micro || true
docker build -t examples/micro .
amp registry push examples/micro
amp stack up micro -c stack.yml

