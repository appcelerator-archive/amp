#!/usr/bin/env bash

amp -k config create FOOBAR platform/tests/config/create/foobar | grep -o -w -E -q '[[:alnum:]]{25}'

amp -k config ls | grep -q 'FOOBAR'

amp -k config rm FOOBAR | grep -q 'FOOBAR'
