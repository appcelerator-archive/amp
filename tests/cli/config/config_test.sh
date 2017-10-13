#!/usr/bin/env bash

amp -k config create FOOBAR tests/cli/config/foobar | grep -o -w -E -q '[[:alnum:]]{25}'

amp -k config ls | grep -q 'FOOBAR'

amp -k config rm FOOBAR | grep -q 'FOOBAR'
