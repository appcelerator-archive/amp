#!/bin/sh

docker version >/dev/null 2>&1
if [ $? -ne 0 ]; then
  DOCKER="sudo docker"
else
  DOCKER=docker
fi
BUILDER_DOCKER_FILE=Dockerfile
SHRINKED_DOCKER_FILE=Dockerfile.shrink
REPOSITORY_NAME=appcelerator/amp
TAG=${1:-latest}

echo -n "building builder image... "
$DOCKER build -f $BUILDER_DOCKER_FILE -t $REPOSITORY_NAME:builder . >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo -n "running builder container... "
$DOCKER kill amp-builder >/dev/null 2>&1; $DOCKER rm amp-builder >/dev/null 2>&1
$DOCKER run -d --name amp-builder $REPOSITORY_NAME:builder >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo -n "copy binary from container... "
$DOCKER cp amp-builder:/go/bin/amplifier ./amplifier >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo -n "building shrinked image... "
$DOCKER build -f $SHRINKED_DOCKER_FILE -t appcelerator/amplifier:$TAG . >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo -n "cleanup... "
rm -f amplifier
$DOCKER kill amp-builder >/dev/null 2>&1
$DOCKER rm amp-builder >/dev/null 2>&1
$DOCKER rmi $REPOSITORY_NAME:builder
echo "OK"
