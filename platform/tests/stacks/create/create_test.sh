#!/bin/bash

amp="amp -s localhost"
$amp stack up -c platform/stacks/visualizer.stack.yml
#$amp stack ls
$amp stack ls 2>/dev/null | grep -q "\svisualizer\s"
