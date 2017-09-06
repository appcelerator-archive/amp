#!/bin/bash

for id in $(amp -k stack ls -q)
do
  amp -k team resource ls | grep -q $id
  code=$?
  if [[ $code -ne 0 ]]; then
    echo "couldn't find resource $id in the team resources:" >&2
    amp -k team resource ls >&2
    exit $code
  fi
done
