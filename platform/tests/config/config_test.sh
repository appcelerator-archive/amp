#!/usr/bin/env bash

# on success, amp config create returns an alphanumeric ID of length 25
amp -k config create FOOBAR platform/tests/config/create/foobar | pcregrep -q '[[:alnum:]]{25}'

amp -k config ls | pcregrep -q 'FOOBAR'

amp -k config rm FOOBAR | pcregrep -q 'FOOBAR'
