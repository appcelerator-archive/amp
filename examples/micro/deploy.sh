#!/bin/sh
amp stack rm -f micro || true
docker build -t examples/micro .
amp registry push examples/micro
amp stack up -f stack.yml micro

