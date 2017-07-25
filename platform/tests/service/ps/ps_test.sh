#!/bin/bash

amp -k service ps pinger_pinger 2>/dev/null | grep -o -w -E -q '[[:alnum:]]{25}'

amp -k service ps pinger_pinger 2>/dev/null | grep -Evq 'SHUTDOWN'

amp -k service ps pinger_pinger -a 2>/dev/null | grep -E -i 'RUNNING|SHUTDOWN'
