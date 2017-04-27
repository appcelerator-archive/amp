#!/bin/bash

# not finding the version is an error
test_has_version() {
  result=$(amp -k version)
  version="v[[:digit:]]{1,3}\.[[:digit:]]{1,3}\.[[:digit:]]{1,3}"
  echo $result | grep -E "Version:[[:space:]]+$version"
}

# finding "not connected" is an error
test_is_connected() {
  result=$(amp -k version)
  echo $result | grep "not connected"
  (( $? == 0 )) && return 1 || return 0
}

