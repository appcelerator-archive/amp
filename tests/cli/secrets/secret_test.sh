#!/bin/bash
# on success, amp secret create returns an alphanumeric ID of length 25
amp -k secret create TEST tests/cli/secrets/create/test | grep -o -w -E -q '[[:alnum:]]{25}'

amp -k secret ls | grep -q 'TEST'

amp -k secret rm TEST | grep -q 'TEST'
