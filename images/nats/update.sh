#!/bin/bash

set -e

if [ $# -eq 0 ] ; then
	echo "Usage: ./update.sh <nats-io/nats-streaming-server tag or branch>"
	exit
fi

VERSION=$1

# cd to the current directory so the script can be run from anywhere.
cd `dirname $0`

cleanup(){
  [[ -n "$TEMP" && -d "$TEMP" ]] && rm -rf "$TEMP"
  docker rm -f nats-streaming-builder
  docker rmi nats-streaming-builder
}

trap cleanup EXIT

echo "Fetching and building nats-streaming-server $VERSION..."

# Create a tmp build directory.
TEMP=$(mktemp -d)

git clone -b $VERSION https://github.com/nats-io/nats-streaming-server $TEMP

docker build -t nats-streaming-builder $TEMP

# Create a dummy nats streaming builder container so we can run a cp against it.
docker create --name nats-streaming-builder nats-streaming-builder

# Update the local binary.
docker cp nats-streaming-builder:/go/bin/nats-streaming-server .

echo "Done"
