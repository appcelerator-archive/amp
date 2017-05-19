#!/bin/bash

amp="amp -s localhost"

test_main() {
  res=$($amp stats | wc -l)
  echo $res
  if [ "$res" -lt 1 ]; then
     exit 1
  fi
}
