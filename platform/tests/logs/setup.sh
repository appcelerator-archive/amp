#!/bin/bash

amp -k stack deploy -c examples/stacks/pinger/pinger.yml pinger
amp -k stack ls 2>/dev/null | grep -q "\spinger\s"
