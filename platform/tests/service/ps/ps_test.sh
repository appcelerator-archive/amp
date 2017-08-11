#!/bin/bash

amp -k service ps pinger_pinger 2>/dev/null | pcregrep -q '[[:alnum:]]{25}'

amp -k service ps pinger_pinger 2>/dev/null | grep -vq 'SHUTDOWN'

amp -k service ps pinger_pinger -a 2>/dev/null | pcregrep -iq 'RUNNING|SHUTDOWN'
