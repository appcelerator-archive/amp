#!/bin/bash

amp -k service ps pinger_pinger 2>/dev/null | pcregrep -q '[[:alnum:]]{25}'

amp -k service ps pinger_pinger 2>/dev/null | pcregrep -vq 'SHUTDOWN'

amp -k service ps pinger_pinger -a 2>/dev/null | pcregrep -i 'RUNNING|SHUTDOWN'
