#!/bin/bash

amp -k stack rm visualizer 2>/dev/null | grep -Eq "[a-z0-9]"
