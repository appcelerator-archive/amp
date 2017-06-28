#!/bin/bash
DIR="$(dirname $0)"
cd "$DIR"
DIR="$(pwd -P)"
image="williamyeh/ansible:alpine3"
image="appcelerator/ansible"

if [[ ! -f $HOME/.aws/credentials ]]; then
  echo "Please configure your aws credentials first"
  exit 1
fi
docker run --rm -v "$DIR:/data" -v "$HOME/.aws:/root/.aws:ro" "$image" ansible-playbook /data/build-ami.yml
