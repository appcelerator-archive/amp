#!/bin/bash

dryrun=0
tag=latest
while getopts ":t:n" opt; do
  case $opt in
  t)  tag=$OPTARG
      echo "tag provided: $tag"
      ;;
  n)  echo "Dry run - no operation will be performed"
      dryrun=1
      ;;
  \?) echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))

images=$*
IMAGEDIR=$(dirname $0)
owner=appcelerator
issues=0
makefile=Makefile.refactor.make


if [ -z "$images" ]; then
  images=$(find $IMAGEDIR -type d \! -name \. -depth 1 | sed -e "s/^\.\///")
fi

echo "Building AMP images"
echo "tag: $tag"
echo "images to build: ${images}" | tr '\n' ' ' ; echo

for image in $images; do
  image=$(basename $image)
  pushd $IMAGEDIR/$image >/dev/null
  if [ -f .travis.yml ]; then
    echo "Warning: a travis configuration is present in $image, it's not supported yet, other build methods will be tried"
  fi
  # check if the image needs a binary that can be built from the root Makefile
  binaries=$(ls make.* 2>/dev/null)
  if [ -n "$binaries" ]; then
    pushd ../.. >/dev/null
    for b in $binaries; do
      b=${b#make.}
      grep -q "^$b:" $makefile
      if [ $? -ne 0 ]; then
        echo "Target $b is not defined in $makefile"
        continue
      fi
      echo "Building $b"
      if [ $dryrun -eq 0 ]; then
        hack/amptools make -f $makefile $b
      fi
    done
    popd >/dev/null
  fi
  # check for a docker compose test file
  if [ -f docker-compose.test.yml ]; then
    method=compose
    name=$(egrep -e "image: .*${image}.*" docker-compose.test.yml | grep -v sut | head -1 | sed "s/.*image:[ 	]*\(.*\)$/\1/")
    if [ -z "${name}" ] && [ -f docker-compose.yml ]; then
      # the service may be extended from another docker compose file
      name=$(egrep -e "image: .*${image}.*" docker-compose.yml | grep -v sut | head -1 | sed "s/.*image:[ 	]*\(.*\)$/\1/")
    fi
    if [ -z "${name}" ]; then
      echo "Unable to guess the image name for tagging, abort"
      ((issues++))
      continue
    fi
  else
    method=docker
    name=$owner/$(echo $image | sed "s/^docker-//")
  fi
  case $method in
  docker)
    echo "Building image $name:$tag (docker build)..."
    if [ $dryrun -eq 0 ]; then
      docker build -t appcelerator/$name:$tag .
    fi
    if [ $? -ne 0 ]; then
      echo "Failed to build $name ($image)"
      exit 1
    fi
    ;;
  compose)
    echo "Building image $name:$tag (docker-compose)..."
    if [ $dryrun -eq 0 ]; then
      docker-compose -f docker-compose.test.yml build && \
      docker-compose -f docker-compose.test.yml run sut && \
      docker tag $name $name:$tag
    fi
    if [ $? -ne 0 ]; then
      echo "Failed to build $name ($image)"
      docker-compose -f docker-compose.test.yml down
      exit 1
    fi
    docker-compose -f docker-compose.test.yml down
    ;;
  *)
    echo "Unknown method: $method"
    exit 1
    ;;
  esac
  popd >/dev/null
done

if [ $issues -eq 0 ]; then
  echo "Images were built successfully"
else
  echo "$issues images couldn't be built"
  exit 1
fi
