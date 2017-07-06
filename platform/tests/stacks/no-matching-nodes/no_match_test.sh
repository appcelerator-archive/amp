#!/usr/bin/env bash

amp -k stack up -c platform/tests/stacks/no-matching-nodes/example.yml
amp -k stack ls 2>/dev/null | pcregrep -q "\s*example\s*0/1\s*1\s*ERROR\s*"
amp -k stack rm example
