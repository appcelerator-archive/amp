#!/bin/sh

tag=${1:-latest}
shift
images=$*
IMAGEDIR=$(dirname $0)
owner=appcelerator
issues=0

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
    docker build -t $name:$tag .
    if [ $? -ne 0 ]; then
      echo "Failed to build $name ($image)"
      exit 1
    fi
    ;;
  compose)
    echo "Building image $name:$tag (docker-compose)..."
    docker-compose -f docker-compose.test.yml build && \
    docker-compose -f docker-compose.test.yml run sut && \
    docker tag $name $name:$tag
    if [ $? -ne 0 ]; then
      echo "Failed to build $name ($image)"
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
