#!/usr/bin/env bash

amp -k stack up -c platform/tests/service/no-matching-nodes/sample.yml
amp -k service ls 2>/dev/null | pcregrep -q "\s*sample\s*0/0\s*0\s*NO MATCHING NODE\s*"
amp -k stack rm sample
