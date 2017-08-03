#!/bin/bash

amp -k stack up -c examples/stacks/ui/ui.stack.yml visualizer
amp -k stack ls 2>/dev/null | pcregrep -q "\svisualizer\s"
