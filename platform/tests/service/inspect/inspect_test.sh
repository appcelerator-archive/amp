#!/bin/bash

amp -k service inspect pinger_pinger 2>/dev/null | grep -q "pinger"
