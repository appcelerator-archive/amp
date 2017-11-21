#!/usr/bin/env bash

amp -k metrics cpu 2>/dev/null | grep -q -Eo '^amp_.*\s*[0-9]+([.][0-9]{3})$'

amp -k metrics cpu --average 2>/dev/null | grep -q -Eo '^amp_.*\s*[0-9]+([.][0-9]{3})$'