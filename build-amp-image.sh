#!/bin/sh

docker version >/dev/null 2>&1
if [ $? -ne 0 ]; then
  DOCKER="sudo docker"
else
  DOCKER=docker
fi
BUILDER_DOCKER_FILE=Dockerfile
SHRINK_DOCKER_FILE=Dockerfile.shrink
REPOSITORY_NAME=appcelerator/amp
TAG=${1:-latest}

echo "removing previous image... "
$DOCKER rmi -f appcelerator/amp:$TAG >&2
echo "OK"

echo "building builder image... "
$DOCKER build -f $BUILDER_DOCKER_FILE -t $REPOSITORY_NAME:builder . >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo "running builder container... "
$DOCKER kill amp-builder >/dev/null 2>&1; $DOCKER rm amp-builder >/dev/null 2>&1
$DOCKER run -d --name amp-builder $REPOSITORY_NAME:builder >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo "copy binary from container... "
$DOCKER cp amp-builder:/go/bin/amp ./amp >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amplifier ./amplifier >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amp-agent ./amp-agent >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amp-log-worker ./amp-log-worker >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amplifier-gateway ./amplifier-gateway >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/adm-server ./adm-server >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/adm-agent ./adm-agent >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amp-function-listener ./amp-function-listener >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
$DOCKER cp amp-builder:/go/bin/amp-function-worker ./amp-function-worker >&2
if [ $? -ne 0 ]; then
  echo "failed"
  exit 1
fi
echo "OK"

echo "building shrunk image... "
cp -p .dockerignore .dockerignore.bak
echo "vendor" >> .dockerignore
$DOCKER build -f $SHRINK_DOCKER_FILE -t $REPOSITORY_NAME:$TAG . >&2
if [ $? -ne 0 ]; then
  mv .dockerignore.bak .dockerignore
  echo "failed"
  exit 1
fi
echo "OK"
mv .dockerignore.bak .dockerignore

echo -n "cleanup... "
rm -f amp amplifier amp-agent amp-log-worker amplifier-gateway adm-server adm-agent amp-function-listener amp-function-worker
$DOCKER kill amp-builder >/dev/null 2>&1
$DOCKER rm amp-builder >/dev/null 2>&1
$DOCKER rmi $REPOSITORY_NAME:builder
echo "OK, image is $REPOSITORY_NAME:$TAG"
