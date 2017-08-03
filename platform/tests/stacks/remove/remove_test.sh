#!/bin/bash

amp -k stack rm visualizer 2>/dev/null | pcregrep -q -e [a-z0-9]
