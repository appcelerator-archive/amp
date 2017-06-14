#!/usr/bin/env bash

amp="amp -s localhost"
set -e
function cleanup {
  $amp user rm user
}
trap cleanup EXIT

$amp user signup --name user --password password --email email@user.amp
$amp user ls -q | wc -l | grep 1
